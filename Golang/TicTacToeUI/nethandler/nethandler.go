package nethandler

import (
	"client"
	"errors"
	"netobjects"
	"reflect"
	"tictactoenet"
)

var tcpClient *client.Client

func init() {
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.SetName{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.Move{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.CreateLobby{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.JoinLobby{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.GetLobby{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.PlayerJoinLobby{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.Match{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.TournamentFinished{}))
}

// AddHandlerFunction må eksistere for å unngå import cycle
/*AddHandlerFunction tar inn som parametere en reflect.Type go en interface{}
Type parameteret er typen struct klienten må motta for å kjøre funksjonen passert inn som andreparameter */
func AddHandlerFunction(t reflect.Type, function interface{}) {
	tcpClient.AddHandlerFunction(t, function)
}

// Connect forsøker å koble klienten på tjeneren og returnerer error dersom denne tilkoblingen mislykkes
func Connect() error {
	tcpClient = client.Create("88.91.152.113", 4444, "tcp")
	conErr := tcpClient.Connect()

	if conErr != nil {
		return conErr
	}
	return nil
}

// CreateLobbyRequest sender en forespørsel til tjeneren om å opprette en lobby
func CreateLobbyRequest(lobbySize int) error {
	if !tcpClient.IsConnected() {
		return errors.New("The client is not connected to the server")
	}

	tcpClient.Write(tictactoenet.CreateLobby{LobbySize: lobbySize}, 0) // Tjenerens ID er alltid 0
	return nil
}

func SendMove(yPos int, xPos int, turnval int, recipiantID int) {
	tcpClient.Write(tictactoenet.Move{XPos: xPos, YPos: yPos, Turnval: turnval}, recipiantID)
}

// IsConnected returnerer true hvis klienten er tilkoblet
func IsConnected() bool {
	if tcpClient == nil {
		return false
	}
	return tcpClient.IsConnected()
}

// SetUsername signaliserer til tjeneren at denne klienten vil bli assosiert med navnet spesifisert i parameteret
func SetUsername(name string) {
	tcpClient.Write(tictactoenet.SetName{Name: name}, 0) // 0 er alltid tjenerens id
}

// JoinLobby forsøker å koble klienten til lobbyen med id'en spesifisert i parameteret
func JoinLobby(lobbyID string) {
	tcpClient.Write(tictactoenet.JoinLobby{LobbyID: lobbyID}, 0) // 0 er alltid tjenerens id
}

// GetLobbyMembers sender en forespørsel til tjener om å sende liste over medlemmer av lobby
func GetLobbyMembers(lobbyID string) {
	tcpClient.Write(tictactoenet.GetLobby{LobbyID: lobbyID}, 0) // 0 er alltid tjenerens id
}

func SendMatchResult(e tictactoenet.Match) {
	tcpClient.Write(e, 0)
}
