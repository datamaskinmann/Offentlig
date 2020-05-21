package gui

import (
	"fmt"
	"ioWrapper"
	"log"
	"nethandler"
	"os"
	"reflect"
	"shell"
	"sort"
	"strconv"
	"strings"
	"tictactoenet"

	"tictactoegame"

	"github.com/sciter-sdk/go-sciter"

	"github.com/sciter-sdk/go-sciter/window"
)

var MainUIHTML string
var PlayOnlineUIHTML string
var PlayOfflineUIHTML string
var Local1v1UIHTML string
var LobbyUIHTML string
var LobbyCreatorUIHTML string
var UICSS string
var WaitingScreenUIHTML string
var NameSetterUIHTML string
var LobbyJoinerUIHTML string
var LeaderboardUIHTML string
var localGame bool = true

var game tictactoegame.TicTacToeGame

var username string

var wd *window.Window
var dom *sciter.Element

func init() {
	var rect = sciter.NewRect(0, 0, getWindowWidth(), getWindowHeight())
	var wdCreationErr error
	wd, wdCreationErr = window.New(sciter.SW_MAIN|sciter.SW_CONTROLS|sciter.SW_ENABLE_DEBUG|sciter.SW_RESIZEABLE, rect)

	wd.DefineFunction("Online", online)
	wd.DefineFunction("Offline", offline)
	wd.DefineFunction("Exit", exit)
	wd.DefineFunction("MMBack", mmBack)
	wd.DefineFunction("Local1v1", local1v1)
	wd.DefineFunction("DoMove", domove)
	wd.DefineFunction("HostTournament", hostTournament)
	wd.DefineFunction("CreateLobby", createLobby)
	wd.DefineFunction("HTBack", htBack)
	wd.DefineFunction("JoinTournament", joinTournament)
	wd.DefineFunction("clickJoin", clickJoin)
	wd.DefineFunction("SetUsername", setUsername)

	if wdCreationErr != nil {
		log.Fatal("Failed to create UI window", wdCreationErr)
	}

	var errors [11]error

	// window.LoadFile funker ikke... Så vi gjør dette og bruker heller window.LoadHtml...
	MainUIHTML, errors[0] = ioWrapper.ReadFile("./menu.html")
	PlayOnlineUIHTML, errors[1] = ioWrapper.ReadFile("./play_online.html")
	PlayOfflineUIHTML, errors[2] = ioWrapper.ReadFile("./play_offline.html")
	Local1v1UIHTML, errors[3] = ioWrapper.ReadFile("./game.html")
	LobbyUIHTML, errors[4] = ioWrapper.ReadFile("./lobby.html")
	LobbyCreatorUIHTML, errors[5] = ioWrapper.ReadFile("./lobbyCreator.html")
	WaitingScreenUIHTML, errors[6] = ioWrapper.ReadFile("./waitingScreen.html")
	NameSetterUIHTML, errors[7] = ioWrapper.ReadFile("./nameSetter.html")
	LobbyJoinerUIHTML, errors[8] = ioWrapper.ReadFile("./lobbyJoiner.html")
	LeaderboardUIHTML, errors[10] = ioWrapper.ReadFile("./leaderboard.html")
	UICSS, errors[10] = ioWrapper.ReadFile("./tictactoeui.css")

	for _, err := range errors {
		if err != nil {
			log.Fatal("Failed to load UI component", err)
		}
	}
}

func getWindowWidth() int {
	width, getWidthErr := shell.Execute("wmic",
		"path",
		"Win32_VideoController",
		"get",
		"CurrentHorizontalResolution")

	if getWidthErr != nil {
		fmt.Println(getWidthErr)
		return 1280
	}

	conv, convErr := strconv.Atoi(strings.SplitAfter(strings.TrimSpace(width), "\n")[1])

	if convErr != nil {
		fmt.Println(convErr)
		return 1280
	}

	return conv
}

func getWindowHeight() int {
	height, getHeightErr := shell.Execute("wmic",
		"path",
		"Win32_VideoController",
		"get",
		"CurrentVerticalResolution")

	if getHeightErr != nil {
		return 720
	}

	conv, convErr := strconv.Atoi(strings.SplitAfter(strings.TrimSpace(height), "\n")[1])

	if convErr != nil {
		return 720
	}

	return conv
}

// SetScene switches the GUI scene
func SetScene(scene string) {
	wd.LoadHtml(scene, "")
	wd.SetOption(sciter.SCITER_SET_DEBUG_MODE, 1)
	wd.SetCSS(UICSS, "", "")
	dom, _ = wd.GetRootElement()
}

func SetTitle(title string) {
	wd.SetTitle(title)
}

func RunAndShow() {
	wd.Show()
	wd.Run()
}

func exit(vals ...*sciter.Value) *sciter.Value {
	os.Exit(0)
	return nil
}

func online(vals ...*sciter.Value) *sciter.Value {
	if !nethandler.IsConnected() {
		setWaitState("Attempting to connect to server...")
		err := nethandler.Connect()

		if err != nil {
			fmt.Println("Failed to connect")
			return nil
		}
		nethandler.AddHandlerFunction(reflect.TypeOf(tictactoenet.CreateLobby{}), createLobbyHandler)
		nethandler.AddHandlerFunction(reflect.TypeOf(tictactoenet.JoinLobby{}), joinLobbyHandler)
		nethandler.AddHandlerFunction(reflect.TypeOf(tictactoenet.GetLobby{}), getLobbyMembersHandler)
		nethandler.AddHandlerFunction(reflect.TypeOf(tictactoenet.Match{}), matchHandler)
		nethandler.AddHandlerFunction(reflect.TypeOf(tictactoenet.PlayerJoinLobby{}), playerJoinLobbyHandler)
		nethandler.AddHandlerFunction(reflect.TypeOf(tictactoenet.Move{}), moveHandler)
		nethandler.AddHandlerFunction(reflect.TypeOf(tictactoenet.TournamentFinished{}), tournamentFinishedHandler)
		online() // Rekursivt kalle funksjonen så vi kommer ut av if statementen :(
	} else {
		switch username {
		case "":
			SetScene(NameSetterUIHTML)
			break
		default:
			SetScene(PlayOnlineUIHTML)
			break
		}
	}
	return nil
}

func tournamentFinishedHandler(e tictactoenet.TournamentFinished, senderID int) {
	SetScene(LeaderboardUIHTML)
	sorted := make(pairList, len(e.Results))

	i := 0
	for k, v := range e.Results {
		sorted[i] = pair{k, v}
		i++
	}

	sort.Sort(sort.Reverse(sorted))
	ul, getUlError := dom.Select("#leaderboardUL")
	if getUlError != nil || len(ul) == 0 {
		fmt.Println("Failed to get UL element")
		return
	}
	for _, a := range sorted {
		li, _ := sciter.CreateElement("li", fmt.Sprintf("%s: %v points", a.Key, a.Value))
		ul[0].Append(li)
	}

}

// Løsning for å sortere
// map etter verdi røvet fra https://stackoverflow.com/questions/18695346/how-to-sort-a-mapstringint-by-its-values
type pair struct {
	Key   string
	Value int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func matchHandler(e tictactoenet.Match, senderID int) {
	turnval = e.Turnval
	opponentID = e.OpponentID
	opponentName = e.OpponentName
	turn = 1
	game = tictactoegame.Create()
	SetScene(Local1v1UIHTML)
	title, _ := dom.Select("#MenuTitle")
	if turnval == turn {
		title[0].SetText("It's your turn!")
		return
	}
	title[0].SetText(fmt.Sprintf("It's %s's turn!", opponentName))
}

func createLobbyHandler(e tictactoenet.CreateLobby, senderID int) {
	wd.SetCSS(UICSS, "", "")
	SetScene(LobbyUIHTML)
	h1, geth1Err := dom.Select("#lobbyId")

	if geth1Err != nil || len(h1) == 0 {
		// Kan oppstå dersom en spiller mottar dette objektet
		// Uten å være på lobby menyen
		fmt.Println("Failed to get lobby id element")
		return
	}
	h1[0].SetText(fmt.Sprintf("Lobby ID:%s", e.LobbyID))

	ul, getUlError := dom.Select("#lobbyMembers")
	if getUlError != nil || len(ul) == 0 {
		fmt.Println("Failed to get UL element")
		return
	}
	li, _ := sciter.CreateElement("li", username)
	ul[0].Append(li)
}

func moveHandler(e tictactoenet.Move, senderID int) {
	elem, _ := dom.Select(fmt.Sprintf("#%v_%v", e.YPos, e.XPos))
	fmt.Println("Turnval:", e.Turnval)
	game.DoMove(e.YPos, e.XPos, e.Turnval)
	switch e.Turnval {
	case 1:
		h1, _ := sciter.CreateElement("h1", "X")
		h1.SetAttr("class", "X")
		elem[0].Append(h1)
		turn = 2
		break
	case 2:
		h1, _ := sciter.CreateElement("h1", "O")
		h1.SetAttr("class", "O")
		elem[0].Append(h1)
		turn = 1
		break
	}
	trySetTitle("It's your turn")

	if game.CheckDraw() {
		nethandler.SendMatchResult(tictactoenet.Match{
			OpponentID:   opponentID,
			OpponentName: opponentName,
			Result:       1,
		})
		trySetTitle("Draw!")
		return
	}

	if game.CheckWin(1) {
		// if turnval == 1 {
		// 	title[0].SetText("You won!")
		// } else {
		trySetTitle(fmt.Sprintf("%s won!", opponentName))
		nethandler.SendMatchResult(tictactoenet.Match{
			OpponentID:   opponentID,
			OpponentName: opponentName,
			Result:       0,
		})
		// }
		return
	}

	if game.CheckWin(2) {
		// 	if turnval == 2 {
		// 		title[0].SetText("You won!")
		// 	} else {
		trySetTitle(fmt.Sprintf("%s won!", opponentName))
		nethandler.SendMatchResult(tictactoenet.Match{
			OpponentID:   opponentID,
			OpponentName: opponentName,
			Result:       0,
		})
		//}
		return
	}
}

// har opplevd problemer med å sette titlen på lobby guien, selv om
// klienten er på den skjermen... Dette er den ekstremt grusomme fiksen
func trySetTitle(text string) {
	title, err := dom.Select("#MenuTitle")
	if err != nil || len(title) <= 0 {
		return
	}
	title[0].SetText(text)
}

func joinLobbyHandler(e tictactoenet.JoinLobby, senderID int) {
	if e.Error != nil {
		// ... håndtere senere
		fmt.Println("Server says:", e.Error)
		return
	}
	nethandler.GetLobbyMembers(e.LobbyID)
	//setWaitState("Attempting to join lobby...")
}

func playerJoinLobbyHandler(e tictactoenet.PlayerJoinLobby, senderID int) {
	ul, ulgetErr := dom.Select("#lobbyMembers")

	if ulgetErr != nil || len(ul) == 0 {
		fmt.Println("Failed to get UL")
		return
	}

	li, _ := sciter.CreateElement("li", e.Name)
	ul[0].Append(li)
}

func getLobbyMembersHandler(e tictactoenet.GetLobby, senderID int) {
	if e.Error != nil {
		SetScene(LobbyJoinerUIHTML)
		fmt.Println("Can't join lobby, server says: ", e.Error)
		return
	}
	SetScene(LobbyUIHTML)
	h1, h1getErr := dom.Select("#lobbyId")

	if h1getErr != nil || len(h1) == 0 {
		// Kan oppstå dersom en spiller mottar dette objektet
		// Uten å være på lobby menyen
		fmt.Println("Failed to get lobby id element")
		return
	}
	h1[0].SetText(fmt.Sprintf("Lobby ID%s", e.LobbyID))
	ul, ulgetErr := dom.Select("#lobbyMembers")
	if ulgetErr != nil || len(ul) == 0 {
		fmt.Println("Failled to get ul")
		return
	}
	for _, a := range e.Names {
		li, _ := sciter.CreateElement("li", a)
		ul[0].Append(li)
	}
}

func setUsername(vals ...*sciter.Value) *sciter.Value {
	if !nethandler.IsConnected() {
		fmt.Println("Cannot set the user's name if the client is not connected to the server...")
		return nil
	}
	nethandler.SetUsername(vals[0].String())
	username = vals[0].String()
	SetScene(PlayOnlineUIHTML)
	return nil
}

func offline(vals ...*sciter.Value) *sciter.Value {
	game = tictactoegame.Create()
	SetScene(PlayOfflineUIHTML)
	return nil
}

func mmBack(vals ...*sciter.Value) *sciter.Value {
	localGame = true
	SetScene(MainUIHTML)
	return nil
}

func htBack(vals ...*sciter.Value) *sciter.Value {
	SetScene(PlayOnlineUIHTML)
	return nil
}

func hostTournament(vals ...*sciter.Value) *sciter.Value {
	localGame = false
	SetScene(LobbyCreatorUIHTML)
	return nil
}

func joinTournament(vals ...*sciter.Value) *sciter.Value {
	SetScene(LobbyJoinerUIHTML)
	return nil
}

func clickJoin(vals ...*sciter.Value) *sciter.Value {
	fmt.Println(vals[0].String())
	nethandler.JoinLobby(vals[0].String())
	localGame = false
	return nil
}

func local1v1(vals ...*sciter.Value) *sciter.Value {
	SetScene(Local1v1UIHTML)
	localGame = true
	turn = 1
	return nil
}

func createLobby(vals ...*sciter.Value) *sciter.Value {
	setWaitState("Attempting to create lobby...")
	nethandler.CreateLobbyRequest(vals[0].Int())
	return nil
}

func setWaitState(waitMessage string) {
	SetScene(WaitingScreenUIHTML)
	h1, _ := dom.Select("#MenuTitle")
	h1[0].SetText(waitMessage)
}

var turn int = 1
var turnval int = -1
var opponentID int = -1
var opponentName string = ""

func domove(vals ...*sciter.Value) *sciter.Value {
	if !localGame && turn != turnval {
		fmt.Println("It's not your turn")
		return nil
	}

	if game.CheckDraw() || game.CheckWin(1) || game.CheckWin(2) {
		return nil
	}

	var yx = make([]int, 2)
	for i, a := range strings.Split(vals[0].String(), "_") {
		yx[i], _ = strconv.Atoi(a)
	}
	moveError := game.DoMove(yx[0], yx[1], turn)

	if !localGame {
		nethandler.SendMove(yx[0], yx[1], turnval, opponentID)
	}

	if moveError != nil {
		fmt.Println(moveError)
	}

	pos, _ := dom.Select(fmt.Sprintf("#%s", vals[0].String()))
	title, _ := dom.Select("#MenuTitle")
	var h1 *sciter.Element

	switch turn {
	case 1:
		h1, _ = sciter.CreateElement("h1", "X")
		h1.SetAttr("class", "X")
		if localGame {
			title[0].SetText("It's O's turn")
		}
		turn = 2
		break
	case 2:
		h1, _ = sciter.CreateElement("h1", "O")
		h1.SetAttr("class", "O")
		if localGame {
			title[0].SetText("It's X's turn")
		}
		turn = 1
		break
	}

	if !localGame {
		title[0].SetText(fmt.Sprintf("It's %s's turn", opponentName))
	}
	pos[0].Append(h1)

	if game.CheckWin(1) {
		if localGame {
			title[0].SetText("X wins!")
		} else { // Den spilleren hvis tur det er, vinner alltid på sin lokale klient på sin tur
			// dermed trenger vi ikke å sjekke if turn == turnval
			nethandler.SendMatchResult(tictactoenet.Match{
				OpponentID:   opponentID,
				OpponentName: opponentName,
				Result:       2,
			})
			title[0].SetText("You win!")
		}
	}
	if game.CheckWin(2) {
		if localGame {
			title[0].SetText("O wins!")
		} else {
			nethandler.SendMatchResult(tictactoenet.Match{
				OpponentID:   opponentID,
				OpponentName: opponentName,
				Result:       2,
			})
			title[0].SetText("You win!")
		}
	}

	if game.CheckDraw() {
		if localGame {
			title[0].SetText("Draw!")
		} else {
			nethandler.SendMatchResult(tictactoenet.Match{
				OpponentID:   opponentID,
				OpponentName: opponentName,
				Result:       1,
			})
		}
	}
	return nil
}
