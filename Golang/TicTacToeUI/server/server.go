package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"event"
	"netobjects"
)

type Server struct {
	ipaddress      string
	port           int
	protocol       string
	listener       net.Listener
	maxConnections int
	// Ready eventet trenger en funksjon med to tomme interface parametere eks: func ready(args interface{}, sender interface{})
	ReadyEvent              *event.Event
	ClientConnectedEvent    *event.Event
	ClientDisconnectedEvent *event.Event
	ServerFullEvent         *event.Event
	LogEvent                *event.Event
	clients                 map[int]net.Conn
	handlerFuncs            map[reflect.Type]interface{}
	isActive                bool
	stopChan                chan (bool)
	terminationChan         chan (int)
}

type ClientConnectedEventArgs struct {
	ID int
}

type ClientDisconnectedEventArgs struct {
	ID int
}

type LogEventArgs struct {
	Msg string
}

// Create oppretter en tjener. IPAddress er IP-addressen til tjeneren. Port er porten til tjeneren. Protocol er "TCP" eller "UDP". Passer inn -1 som MaxConnections for uendlig mulige tilkoblinger
func Create(IPAddress string, Port int, Protocol string, MaxConnections int) *Server {
	return &Server{
		ipaddress:               IPAddress,
		port:                    Port,
		protocol:                Protocol,
		maxConnections:          MaxConnections,
		ReadyEvent:              &event.Event{},
		ClientConnectedEvent:    &event.Event{},
		ClientDisconnectedEvent: &event.Event{},
		ServerFullEvent:         &event.Event{},
		LogEvent:                &event.Event{},
		clients:                 make(map[int]net.Conn),
		handlerFuncs:            make(map[reflect.Type]interface{}),
		stopChan:                make(chan (bool)),
		terminationChan:         make(chan (int)),
	}
}

// Start servereren passer inn -1 som maxConnections for uendlig mulige tilkoblinger
func (server *Server) Start() error {
	server.LogEvent.RaiseAsync(LogEventArgs{Msg: "Starting server..."}, server)
	var listenerError error
	// Opprette en lytter på ip, port og protokoll fra structets verdier
	server.listener, listenerError = net.Listen(strings.ToLower(server.protocol), fmt.Sprint(server.ipaddress, ":", server.port))
	// Oppstod det en feil ved opprettingen av lytteren?
	if listenerError != nil {
		// Ja -> returner feil
		return listenerError
	}

	// Opprette en goroutine som kjører server.listen
	go server.listen()

	return nil
}

// Stop stopper tjeneren
func (server *Server) Stop() error {
	if server.listener == nil {
		server.LogEvent.RaiseAsync(LogEventArgs{Msg: "Cannot stop a server that has not been started"}, server)
		return errors.New("Cannot stop a server that has not been started")
	}
	close(server.stopChan)
	err := server.listener.Close()

	for _, c := range server.clients {
		c.Close()
	}

	if err != nil {
		server.LogEvent.RaiseAsync(LogEventArgs{Msg: fmt.Sprintf("Error while stopping server: %s", err)}, server)
		return err
	}
	return nil
}

//TerminateClient tar som parameter ID'en til en klient.
// Tilkoblingen med denne klienten vil så bli avsluttet.
// Returnerer en error hvis en klient med denne ID'en ikke er tilkoblet
func (server *Server) TerminateClient(id int) error {
	server.LogEvent.RaiseAsync(LogEventArgs{Msg: fmt.Sprint("Attempting to terminate client with id", id)},
		server)
	if server.clients[id] == nil {
		return errors.New(fmt.Sprint("No client with ID:", id, "is connected"))
	}
	server.clients[id].Close()
	server.terminationChan <- id
	return nil
}

// Global waitgroup
var handleWaitGroup sync.WaitGroup

func (server *Server) listen() {
	defer func() {
		server.LogEvent.RaiseAsync(LogEventArgs{Msg: "Succesfully Closed server..."}, server)
		runtime.Goexit()
	}()
	// Holder styr på ID'en til klientene
	var counter int = 0

	server.isActive = true

	server.ReadyEvent.RaiseAsync(new(interface{}), server)
	for {
		select {
		case <-server.stopChan:
			server.isActive = false
			runtime.Goexit()
		default:
			// Akseptere klient som ønsker å tilkoble
			client, clientError := server.listener.Accept()
			// Oppstod det en feil under aksepteringen av klienten?
			if clientError != nil {
				// Ja -> continue (gå til toppen av løkken)
				server.LogEvent.RaiseAsync(LogEventArgs{Msg: fmt.Sprint("Error while accepting client:", clientError)},
					server)
				continue
			}
			// Inkrementere waitgroupen
			handleWaitGroup.Add(1)

			server.LogEvent.RaiseAsync(LogEventArgs{Msg: "Client connected!"}, server)
			counter++

			// Klienten kobla på tjeneren, legg den til
			// i server.clients mappet og sett keyen til
			// verdien av counter
			server.clients[counter] = client
			// Opprette SetID struct og sette ID feltet til
			// verdien av counter
			setClientID := netobjects.SetID{ID: counter}

			server.ClientConnectedEvent.RaiseAsync(ClientConnectedEventArgs{ID: counter}, server)
			// Sende setClientID structet til klienten slik
			// den setter ID'en sin
			server.Write(setClientID, counter)

			// Goroutine som leser data fra klienten
			go server.readClient(setClientID.ID)

			// Send ClientConnected til alle bortsett fra den som nettopp kobla til
			server.WriteAllExcept(netobjects.ClientConnected{ID: counter}, counter)

			// Er maxConnections ikke -1 og er counter større enn maxConnections?
			if server.maxConnections != -1 && counter >= server.maxConnections {
				// Vente til waitgroupen vår er ferdig
				handleWaitGroup.Wait()
				server.LogEvent.RaiseAsync(
					LogEventArgs{Msg: fmt.Sprint("Max connections (", server.maxConnections, ") achieved, closing listener...")},
					server)
				server.ServerFullEvent.RaiseAsync(new(interface{}), server)
				// Bryte for-løkken
				return
			}
		}
	}
}

// IsActive returnerer en bool som representerer om listeneren på serrveren har blit starta riktig
func (server Server) IsActive() bool {
	return server.isActive
}

func (server *Server) readClient(id int) {
	defer func() {
		// Slette klienten fra server.clients mappet
		delete(server.clients, id)
		// Gi beskjed om at klienten har frakoblet
		server.ClientDisconnectedEvent.RaiseAsync(ClientDisconnectedEventArgs{ID: id}, server)
		server.LogEvent.RaiseAsync(LogEventArgs{Msg: fmt.Sprintf("Client %d disconnected! :(", id)}, server)
		server.WriteAllExcept(netobjects.ClientDisconnected{ID: id}, id)
		runtime.Goexit()
	}()
	// Disconnect timeout = 1 sekund
	var disconnectTimeout time.Duration = time.Second
	// Inkrementeres dersom det oppstår en feil under lesingen av en klient
	var tries int = 0

	// Finne tilkoblingen til klient i server.clients mappet vårt
	// ved hjelp av id parameteret
	var client net.Conn = server.clients[id]
	handleWaitGroup.Done()
	// Opprette buffered input/output-leser på klientens nettstrøm
	reader := bufio.NewReader(client)
	for {
		select {
		case <-server.stopChan:
			runtime.Goexit()
		case termID := <-server.terminationChan:
			if termID == id {
				return
			}
			continue
		default:
			// Lese data fra klienten til delimiter byte '\n' er nådd
			value, readError := reader.ReadBytes('\n')
			// Oppstod det en feil under lesingen av klienten?
			if readError == io.EOF {
				// Ja -> frys goroutinen i ett sekund
				time.Sleep(disconnectTimeout)
				// Er tries lik 5?
				// (Har det gått 5 sekunder uten at
				// tjeneren lyktes i å lese data fra klienten?)
				if tries == 5 {
					server.LogEvent.RaiseAsync(LogEventArgs{Msg: "Client failed to respond after 5 seconds...\nClient disconnected"},
						server)
					// Lukke tilkoblingen til klienten
					return
				}
				// Inkrementere tries og gå til toppen av for-løkken
				// (forsøke å lese data fra klienten igjen)
				tries++
				continue
			} else if readError != nil {
				server.LogEvent.RaiseAsync(LogEventArgs{Msg: "Error while reading from client"}, server)
			}
			tries = 0 // Klienten klarte å koble seg på igjen før 5 sec hadde gått
			// På all data som sendes fra tjener til klient, tilføyes
			// den tiltenkte mottakerens ID på slutten av bytestrømmen
			receiver := int(value[len(value)-2])

			// Er den tiltenkte mottakerens ikke 0
			// (Tjeneren har alltid ID 0)
			if receiver != 0 {
				// Forsøke å videresende dataen
				// Til den tiltenkte mottakeren
				writeError := server.write(value)
				if writeError != nil {
					server.LogEvent.RaiseAsync(LogEventArgs{Msg: writeError.Error()}, server)
				}
				continue
			}
			// Finne senderens ID som er den 3. siste
			// byten i bytestrømmen
			sender := int(value[len(value)-3])

			// Deserialisere bytestrømmen til tiltenkt
			// struktur
			xObj := netobjects.GetObject(value)

			// Finnes det en håndteringsfunksjon for denne tiltenkte
			// strukturen?
			if server.handlerFuncs[reflect.TypeOf(xObj).Elem()] == nil {
				// Nei -> panic
				panic(fmt.Sprint("No handler function found for objects of type",
					reflect.TypeOf(xObj)))
			}
			// Kalle ved hjelp av refleksjon håndteringsfunksjonen for denne strukturen
			go reflect.ValueOf(server.handlerFuncs[reflect.TypeOf(xObj).Elem()]).
				Call([]reflect.Value{reflect.ValueOf(xObj).Elem(), reflect.ValueOf(sender)})
		}
	}
}

// AddHandlerFunction legger til en funksjon som blir kallet når tjeneren mottar et objekt av typen spesifisert i første parameter
func (server *Server) AddHandlerFunction(ObjectType reflect.Type, function interface{}) error {
	if reflect.TypeOf(function).Kind() != reflect.Func {
		return errors.New("The function you specified is invalid")
	}
	server.handlerFuncs[ObjectType] = function
	return nil
}

// WriteAll sender et objekt spesifisert i Object parameteret til alle tilkoblede klienter
func (server Server) WriteAll(Object interface{}) {
	// Iterere gjennom server.clients mappet
	for i := 1; i <= len(server.clients); i++ {
		// Er klienten vi for øyeblikket itererer på ikke nil?
		if server.clients[i] != nil {
			// Ja -> kall server.Write,
			// passer inn iden til klienten
			// og Object parameteret som parametere
			server.Write(Object, i)
		}
	}
}

// WriteAllExcept sender et objekt spesifisert i Object til alle tilkoblede klienter
//bortsett fra de med idene spesifisert i Exceptions parameteret
func (server Server) WriteAllExcept(Object interface{}, Exceptions ...int) {
	for i := 1; i <= len(server.clients); i++ {
		if server.clients[i] == nil {
			continue
		}

		var contains bool
		for _, a := range Exceptions {
			if i == a {
				contains = true
				break
			}
		}
		if !contains {
			// Denne ID'en skal ikke unngås
			// passer inn iden til klient
			// Send objektet til tiltenkt mottaker
			server.Write(Object, i)
		}
	}
}

/* Write sender et objekt spesifisert i Object parameteret til klienten som har
ID'en spesifisert i receiverID parameteret og returnerer error hvis ikke en klient med
den ID'en er tilkobla tjeneren*/
func (server *Server) Write(Object interface{}, receiverID int) error {
	// Er en klient med ID spesifisert i receiverID parameteret tilkobla
	// tjeneren?
	if server.clients[receiverID] == nil {
		// Nei -> returner error
		return errors.New(fmt.Sprint("No client with ID", receiverID, "is connected to the server!"))
	}
	// Serialisere Object parametere til bytestrøm
	data, _ := json.Marshal(Object)

	// Opprette NetObject og tildele
	// de nødvendige feltene verdier
	netobj := netobjects.NetObject{
		Unmarshaltype: reflect.TypeOf(Object).String(),
		Data:          string(data),
	}

	// Serialisere NetObjectet vi nettopp oppretta til bytestrøm
	obj, _ := json.Marshal(netobj)
	// Tilføye 0 (ID til tjeneren), den tiltenkte mottakerens ID og '\n' som delimiter
	// byte på slutten av bytestrømmen
	obj = append(obj, byte(0), byte(receiverID), '\n')
	server.clients[receiverID].Write(obj)
	return nil
}

// Denne funksjonen blir brukt internt for å passere sende objekter fra en klient til en annen
func (server *Server) write(data []byte) error {
	// Finne den tiltenkte mottakeren av dataen
	// som alltid er den 2. siste byten i bytestrømmen
	receiver := int(data[len(data)-2])

	// Finnes en klient med den tiltenkte mottakerens ID
	// i server.clients mappet vårt?
	if server.clients[receiver] == nil {
		// Nei -> returner feil
		return errors.New(fmt.Sprint("No client with ID", receiver, "is connected to the server!"))
	}

	server.clients[receiver].Write(data)
	return nil
}
