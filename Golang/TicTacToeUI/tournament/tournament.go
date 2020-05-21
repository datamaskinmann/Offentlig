package tournament

import (
	"event"
	"fmt"
	"math/rand"
	"player"
	"strings"
	"time"
)

// Tournament structet er 'grunnmuren' til hele tictactoe prosjektet vårts tournament mode
type Tournament struct {
	Players             []*player.Player
	ID                  string
	size                int
	TournamentFullEvent *event.Event
}

var tournMap map[string]*Tournament = make(map[string]*Tournament)

// New oppretter et nytt tournament og genererer en tilfeldig tournament ID.
func New(size int) *Tournament {

	var id string
	for { // Bare i tilfelle det skulle mirakuløst skje at vi genererer to like tournament IDer
		id = generateTournamentID()
		if _, ok := tournMap[id]; !ok {
			break
		}
	}

	var retTournament = &Tournament{Players: []*player.Player{}}

	retTournament.size = size
	retTournament.TournamentFullEvent = &event.Event{}
	retTournament.ID = id

	tournMap[id] = retTournament

	return retTournament
}

func (tournament Tournament) Size() int {
	return tournament.size
}

func generateTournamentID() string {
	rand.Seed(time.Now().UnixNano())

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÆØÅ" +
		"abcdefghijklmnopqrstuvwxyzæøå" +
		"0123456789")

	lengde := 15

	var b strings.Builder

	for i := 0; i < lengde; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	str := b.String()
	return str
}

//AddPlayer legger en spiller til i tournamentet
func (tournament *Tournament) AddPlayer(player *player.Player) {
	if len(tournament.Players)+1 > tournament.size {
		fmt.Println("The tournament is already full")
		return
	}
	tournament.Players = append(tournament.Players, player)
	if len(tournament.Players) == tournament.size {
		tournament.TournamentFullEvent.Raise(new(interface{}), tournament)
	}
}

// RemovePlayer fjerner en spiller fra tournamentet
func (tournament *Tournament) RemovePlayer(id int) {
	for i := 0; i < len(tournament.Players); i++ {
		if id == tournament.Players[i].ID {
			tournament.Players[i] = tournament.Players[len(tournament.Players)-1]
			tournament.Players = tournament.Players[:len(tournament.Players)-1]
		}
	}
}

func (tournament *Tournament) getPlayer(id int) *player.Player {
	for i := 0; i < len(tournament.Players); i++ {
		if id == tournament.Players[i].ID {
			return tournament.Players[i]
		}
	}
	return nil
}

// GetPlayers returns a list of the tournaments players
func (tournament Tournament) GetPlayers() []*player.Player {
	return tournament.Players
}

func (tournament *Tournament) getIDs() []int {
	var playerIDs []int

	for i := 0; i < len(tournament.Players); i++ {
		playerIDs = append(playerIDs, tournament.Players[i].ID)
	}
	return playerIDs
}

func except(slice1 []int, slice2 []int) []int {
	var retArr []int

	for i := 0; i < len(slice1); i++ {
		var contains bool = false

		for x := 0; x < len(slice2); x++ {
			if slice1[i] == slice2[x] {
				contains = true
				break
			}
		}

		if !contains {
			retArr = append(retArr, slice1[i])
		}
	}
	return retArr
}

// MatchPlayer forsøker å finne en match til spilleren med ID passert inn som parameter
func (tournament *Tournament) MatchPlayer(id int) {
	paramPlayer := tournament.getPlayer(id)

	if paramPlayer == nil {
		fmt.Println("No such player exists, try again.")
		return
	}

	if paramPlayer.GetMatch() != -1 {
		fmt.Println("Player is already matched...")
		return
	}
	IDColl := tournament.getIDs()

	matchAgainst := except(IDColl, paramPlayer.PlayedAgainst)
	fmt.Println(matchAgainst)

	for i := 0; i < len(matchAgainst); i++ {
		var buffer *player.Player = tournament.getPlayer(matchAgainst[i])
		if buffer.MatchedAgainst == -1 {
			paramPlayer.OnMatch(buffer.ID)
			buffer.OnMatch(id)
			return
		}
	}
}

func (tournament *Tournament) Finished() bool {
	for i := 0; i < len(tournament.Players); i++ {
		if len(tournament.Players) != len(tournament.Players[i].PlayedAgainst) {
			return false
		}
	}
	return true
}

// GetTournament returnerer tournamentet med ID'en spesifisert i parameteret eller nil dersom et slikt tournament ikke ekisterer
func GetTournament(id string) *Tournament {
	fmt.Println(id)
	for _, a := range tournMap {
		fmt.Println(a)
	}
	return tournMap[id]
}
