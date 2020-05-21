package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"

	"event"
	"netobjects"
	"runtime"
)

type Client struct {
	ipaddress               string
	port                    int
	protocol                string
	connected               bool
	server                  net.Conn
	id                      int
	ReadyEvent              *event.Event
	LostConnectionEvent     *event.Event
	LogEvent                *event.Event
	OtherClientConnected    *event.Event
	OtherClientDisconnected *event.Event
	handlerFuncs            map[reflect.Type]interface{}
	stopChan                chan (interface{})
}

// ReadyEventArgs brukes for å kalle ReadyEvent.Raise()
type ReadyEventArgs struct {
	ID int
}

// LogEventArgs brukes for å kalle LogEvent
type LogEventArgs struct {
	Msg string
}

// OtherClientConnectedEventArgs brukes for å kalle OtherClientConnectedEvent
type OtherClientConnectedEventArgs struct {
	ID int
}

// OtherClientDisconnectedEventArgs brukes for å kalle OtherClientDisconnectedEvent
type OtherClientDisconnectedEventArgs struct {
	ID int
}

/* Create oppretter en klient.
IPAddress parameteret er IP-Adressen du ønsker å binde klienten til
Port parameteret er porten du ønsker å binde klient til
Protocol parameteret er tcp/udp */
func Create(IPAddress string, Port int, Protocol string) *Client {
	client := Client{
		ipaddress:               IPAddress,
		port:                    Port,
		protocol:                Protocol,
		handlerFuncs:            make(map[reflect.Type]interface{}),
		id:                      -1,
		ReadyEvent:              &event.Event{},
		LostConnectionEvent:     &event.Event{},
		LogEvent:                &event.Event{},
		OtherClientConnected:    &event.Event{},
		OtherClientDisconnected: &event.Event{},
		stopChan:                make(chan (interface{})),
	}
	// Legge til en håndteringsfunksjon for når klienten mottar et objekt av typen 'SetID'.
	// Dette objektet sender tjeneren når klienten kobler på og brukes for å sette tjenerens ID.
	client.AddHandlerFunction(reflect.TypeOf(netobjects.SetID{}), client.handleSetID)
	// Legge til en håndteringsfunksjon for når klient mottar et objekt av typen 'ClientConnected'
	// Dette objektet sender tjeneren når en annen klient kobler til
	client.AddHandlerFunction(reflect.TypeOf(netobjects.ClientConnected{}), client.handleOtherClientConnected)
	// Legge til en håndteringsfunksjon for når klient mottar et objekt av typen 'ClientDisconnected'
	// Dette objektet sender tjeneren når en annen klient kobler fra
	client.AddHandlerFunction(reflect.TypeOf(netobjects.ClientDisconnected{}), client.handleOtherClientDisconnected)

	return &client
}

// Connect kobler klienten på tjeneren.
func (client *Client) Connect() error {
	// Forsøke å koble på tjeneren
	connection, connectionError := net.Dial(strings.ToLower(client.protocol),
		fmt.Sprint(client.ipaddress, ":", client.port))

	client.server = connection
	// Oppstod det en feil under tilkoblingen?
	if connectionError != nil {
		client.LogEvent.RaiseAsync(LogEventArgs{Msg: fmt.Sprintf("Error while connecting to %s:%d", client.ipaddress, client.port)}, client)
		return connectionError
	}

	client.LogEvent.RaiseAsync(LogEventArgs{Msg: "Connected to server!"}, client)
	client.connected = true
	// Opprette en GoRoutine som leser data fra tjeneren
	go client.readServer()
	return nil
}

// Disconnect kobler klienten fra tjeneren
func (client *Client) Disconnect() {
	client.server.Close()
	close(client.stopChan)
}

// IsConnected returnerer true hvis klienten er tilkoblet en tjener
func (client *Client) IsConnected() bool {
	return client.connected
}

// GetID returnerer ID'en til denne klienten
// Returnerer en error dersom ID'en ikke har blitt satt
// av tjeneren enda
func (client *Client) GetID() (int, error) {
	// ID 0 = alltid tjeneren, dermed hvis klient skulle på ett eller
	// annet magisk vis hatt ID 0, er dette en feil
	if client.id <= 0 {
		return -1, errors.New("The client's ID is either invalid or not yet set by the server")
	}
	return client.id, nil
}

/*AddHandlerFunction tar inn som parametere en reflect.Type go en interface{}
Type parameteret er typen struct klienten må motta for å kjøre funksjonen passert inn som andreparameter */
func (client *Client) AddHandlerFunction(Type reflect.Type, function interface{}) {
	client.handlerFuncs[Type] = function
}

// readServer leser data den mottar fra tjeneren
func (client *Client) readServer() {
	defer func() {
		client.LogEvent.RaiseAsync(LogEventArgs{Msg: "Disconnected or lost connection"}, client)
		client.connected = false
		runtime.Goexit()
	}()
	// Opprette en buffered input/output-leser på tjenerens nettstrøm
	reader := bufio.NewReader(client.server)
	for {
		select {
		case <-client.stopChan:
			return
		default:
			var buffer []byte
			for b, err := reader.Peek(1); err != io.EOF && len(b) != 0; {
				value, err := reader.ReadByte()
				if err != nil {
					client.LogEvent.RaiseAsync(LogEventArgs{Msg: fmt.Sprint("Error while reading from server:", err)}, client)
					break
				}
				buffer = append(buffer, value)
				if value == 0x00A {
					// Finne ID'en til senderen av objektet
					// Dette vil alltid være den 3. siste byten i strømmen
					sender := int(buffer[len(buffer)-3])
					// Hente deserialisere dataen til dens tiltenkte struktur
					xObj := netobjects.GetObject(buffer)
					// Finnes det håndteringsfunksjon for denne strukturen?
					if client.handlerFuncs[reflect.TypeOf(xObj).Elem()] == nil {
						panic(fmt.Sprint("No handler function for object of type", reflect.TypeOf(xObj).Elem()))
					}
					// Kalle håndteringsfunksjonen på en goroutine ved hjelp av refleksjon
					go reflect.ValueOf(client.handlerFuncs[reflect.TypeOf(xObj).Elem()]).
						Call([]reflect.Value{reflect.ValueOf(xObj).Elem(), reflect.ValueOf(sender)})
					buffer = []byte{}
				}
				continue
			}
			// Kanskje adde timeout?
			_, err := reader.ReadByte()
			if err == io.EOF {
				return
			}
		}
	}
}

func (client *Client) handleSetID(ID netobjects.SetID, senderID int) {
	if senderID != 0 { // hvis setID objektet ikke ble sendt av tjeneren (tjeneren har alltid ID 0)
		return
	}
	client.id = ID.ID
	client.ReadyEvent.RaiseAsync(ReadyEventArgs{ID: client.id}, client)
}

func (client *Client) handleOtherClientConnected(e netobjects.ClientConnected, senderID int) {
	if senderID != 0 {
		return // Objektet ble ikke sendt av tjeneren og kan dermed
		// ikke anses genuint
	}
	// Har noen abonnert på eventet?
	if client.OtherClientConnected.Length() > 0 {
		client.OtherClientConnected.RaiseAsync(OtherClientConnectedEventArgs{e.ID}, client)
	}
}

func (client *Client) handleOtherClientDisconnected(e netobjects.ClientDisconnected, senderID int) {
	if senderID != 0 {
		return // Objektet ble ikke sendt av tjeneren o kan dermed
		// ikke anses som genuint
	}

	if client.OtherClientDisconnected.Length() > 0 {
		client.OtherClientDisconnected.RaiseAsync(OtherClientDisconnectedEventArgs{e.ID}, client)
	}
}

/* Write sender et objekt spesifisert i Object parameteret til
en tiltenkt mottaker med ID spesifisert i receiverID parameteret*/
func (client *Client) Write(Object interface{}, receiverID int) {
	if !netobjects.IsNetworkableObject(Object) {
		fmt.Println(fmt.Sprint("Object ", reflect.TypeOf(Object), " is not networkable"))
		return
	}

	if client.id == -1 {
		client.LogEvent.RaiseAsync(LogEventArgs{Msg: "Vennligst vent til ID'en til klienten har blit satt av serveren før du sender objekter"}, client)
		return
	}

	// Serialisere objektet passer inn som parameter til JSON
	data, _ := json.Marshal(Object)
	// Opprette instans av NetObject struct og tildele de nødvendige feltene
	netObj := netobjects.NetObject{
		Unmarshaltype: reflect.TypeOf(Object).String(),
		Data:          string(data),
	}
	// Serialisere netobjectet vi nettopp oppretta til bytestrøm
	obj, _ := json.Marshal(netObj)
	// Tilføye klientens ID, mottakerens ID og delimiter byten til slutten av bytestrømmen
	obj = append(obj, byte(client.id), byte(receiverID), '\n')
	// Sende dataen til tjeneren
	client.server.Write(obj)
}
