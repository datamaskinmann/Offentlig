package main

import (
	"fmt"
	"log"
	"math/rand"
	"netobjects"
	"os"
	"os/signal"
	"player"
	"reflect"
	"server"
	"sync"
	"syscall"
	"tictactoenet"
	"time"
	"tournament"
)

var tcpServer *server.Server
var players map[int]*player.Player = make(map[int]*player.Player)

var blockMain sync.WaitGroup

func init() {
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.SetName{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.Move{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.CreateLobby{}))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.GetLobby{}))
	netobjects.AddNetworkableObject(reflect.TypeOf((tictactoenet.JoinLobby{})))
	netobjects.AddNetworkableObject(reflect.TypeOf(tictactoenet.Match{}))

}

func registerHandlers() {
	tcpServer.AddHandlerFunction(reflect.TypeOf(tictactoenet.SetName{}), onSetName)
	tcpServer.AddHandlerFunction(reflect.TypeOf(tictactoenet.CreateLobby{}), onCreateLobby)
	tcpServer.AddHandlerFunction(reflect.TypeOf(tictactoenet.GetLobby{}), onGetLobby)
	tcpServer.AddHandlerFunction(reflect.TypeOf(tictactoenet.JoinLobby{}), onJoinLobby)
	tcpServer.AddHandlerFunction(reflect.TypeOf(tictactoenet.Match{}), onMatchFinished)
}

func main() {
	fmt.Println("Trykk CTRL + C for å avslutte")
	//var ip string = os.Args[1]

	var ip string = "0.0.0.0"
	var port int = 4444
	var portReadError error = nil
	//port, portReadError := strconv.Atoi(os.Args[2])

	if portReadError != nil {
		log.Fatalf("%d is not a valid port!", port)
	}

	tcpServer = server.Create(ip, port, "tcp", -1)
	registerHandlers()
	servErr := tcpServer.Start()

	if servErr != nil {
		log.Fatalf("Unable to start server: %s", servErr.Error())
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	go func() {
		<-sigs
		tcpServer.Stop()
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	tcpServer.LogEvent.Add(onServerLog)

	blockMain.Add(1)

	blockMain.Wait()
}

func onServerLog(e server.LogEventArgs, sender interface{}) {
	fmt.Println(e.Msg)
}

func onSetName(e tictactoenet.SetName, senderID int) {
	players[senderID] = player.New(e.Name, senderID)
}

func onCreateLobby(e tictactoenet.CreateLobby, senderID int) {
	var lobby = tournament.New(e.LobbySize)
	lobby.AddPlayer(players[senderID])
	players[senderID].InLobby = lobby.ID
	lobby.TournamentFullEvent.Add(onLobbyFull)

	tcpServer.Write(tictactoenet.CreateLobby{LobbyID: lobby.ID}, senderID)
}

func onGetLobby(e tictactoenet.GetLobby, senderID int) {
	if tournament.GetTournament(e.LobbyID) == nil {
		e.Error = fmt.Errorf("No lobby with ID: (%s) exists", e.LobbyID)
		tcpServer.Write(e, senderID)
		return
	}
	var retPlName []string
	var retPlID []int

	for _, a := range tournament.GetTournament(e.LobbyID).GetPlayers() {
		retPlName = append(retPlName, a.Name)
		retPlID = append(retPlID, a.ID)
	}
	e.IDs = retPlID
	e.Names = retPlName

	fmt.Println("Sending lobby data to", senderID)
	tcpServer.Write(e, senderID)
}

func onJoinLobby(e tictactoenet.JoinLobby, senderID int) {
	if tournament.GetTournament(e.LobbyID) == nil {
		fmt.Println("No such lobby exists join lobby")
		e.Error = fmt.Errorf("No lobby with ID: (%s) exists", e.LobbyID)
		tcpServer.Write(e, senderID)
		return
	}
	var tourn = tournament.GetTournament(e.LobbyID)
	if len(tourn.Players) < tourn.Size() {
		tourn.AddPlayer(players[senderID])
		players[senderID].InLobby = e.LobbyID
		tcpServer.Write(e, senderID)
		for _, a := range tourn.GetPlayers() {
			if a.ID != senderID {
				tcpServer.Write(tictactoenet.PlayerJoinLobby{Name: players[senderID].Name}, a.ID)
			}
		}
	}
}

func onLobbyFull(e interface{}, s interface{}) {
	var tourn = s.(*tournament.Tournament)
	genAllMatches(tourn)
}

func onMatchFinished(e tictactoenet.Match, senderID int) {
	players[senderID].OnMatchFinished(e.Result)
	var tourn = tournament.GetTournament(players[senderID].InLobby)
	if !tourn.Finished() {
		tourn.MatchPlayer(senderID)
		if players[senderID].GetMatch() != -1 {
			var turnVals = genTurnvals()
			tcpServer.Write(tictactoenet.Match{OpponentID: players[senderID].MatchedAgainst,
				OpponentName: players[players[senderID].MatchedAgainst].Name,
				Turnval:      turnVals[0]}, senderID)
			tcpServer.Write(tictactoenet.Match{OpponentID: players[senderID].ID,
				OpponentName: players[senderID].Name,
				Turnval:      turnVals[1]},
				players[senderID].MatchedAgainst)
		}
	} else {
		var result map[string]int = make(map[string]int)
		for _, v := range tourn.GetPlayers() {
			result[v.Name] = v.Points()
		}

		for i := 0; i < len(tourn.GetPlayers()); i++ {
			tcpServer.Write(tictactoenet.TournamentFinished{Results: result}, tourn.GetPlayers()[i].ID)
		}
	}
}

func genAllMatches(tourn *tournament.Tournament) {
	var genList []int
	for _, a := range tourn.GetPlayers() {
		var contains bool
		for _, b := range genList {
			if b == a.ID {
				contains = true
				break
			}
		}
		if contains {
			continue
		}
		tourn.MatchPlayer(a.ID)
		if a.MatchedAgainst == -1 {
			continue
		}
		genList = append(genList, a.ID)
		genList = append(genList, a.MatchedAgainst)
		var turnVals = genTurnvals()
		tcpServer.Write(tictactoenet.Match{OpponentID: a.MatchedAgainst,
			OpponentName: players[a.MatchedAgainst].Name,
			Turnval:      turnVals[0]}, a.ID)
		tcpServer.Write(tictactoenet.Match{OpponentID: a.ID,
			OpponentName: players[a.ID].Name,
			Turnval:      turnVals[1]}, a.MatchedAgainst)
	}
}

func genTurnvals() []int {
	var ret []int
	rand.Seed(time.Now().UnixNano())
	var rnd = rand.Intn(2)
	switch rnd {
	case 0:
		ret = append(ret, 1)
		ret = append(ret, 2)
		return ret
	default:
		ret = append(ret, 2)
		ret = append(ret, 1)
		return ret
	}
}
