package player

import (
	"fmt"
)

// Player structet er structet som tournament bygger på
type Player struct {
	Name           string
	ID             int
	PlayedAgainst  []int
	MatchedAgainst int
	points         int
	InLobby        string
}

type MatchEventArgs struct {
	OpponentID int
}

type MatchFinishedEventArgs struct {
	OpponentID int
}

// New oppretter en ny spiller med navn og id fra parameter
func New(name string, id int) *Player {
	return &Player{
		ID:             id,
		Name:           name,
		PlayedAgainst:  []int{id},
		MatchedAgainst: -1,
	}
}

// OnMatch kalles når en spiller får en match. id parameteret en ID'en til motstanderen.
func (player *Player) OnMatch(id int) {
	fmt.Println(player.ID, "match", id)
	player.MatchedAgainst = id
}

// OnMatchFinished kalles når spilleren har spilt ferdig en kamp. Den inkrementerer spillerens 'points'
// egenskap med det som blir passert inn som parameter
func (player *Player) OnMatchFinished(points int) {
	player.points += points
	player.PlayedAgainst = append(player.PlayedAgainst, player.MatchedAgainst)
	player.MatchedAgainst = -1
}

// GetMatch returnerer spillerens match ID eller -1 dersom spilleren ikke er matched
func (player *Player) GetMatch() int {
	return player.MatchedAgainst
}

// Points returnerer spillerens poeng
func (player Player) Points() int {
	return player.points
}
