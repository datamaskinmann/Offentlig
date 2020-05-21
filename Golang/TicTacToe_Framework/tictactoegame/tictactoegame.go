package tictactoegame

import (
	"errors"
	"fmt"
)

// En laaaaang string brukt for å tegne brettet og putte verdiene på riktig plass
var brettFormat string = "         |         |\n         |         |\n    %s    |    %s    |    %s    \n         |         |\n_________|_________|_________\n         |         |\n         |         |\n    %s    |    %s    |    %s     \n         |         |\n_________|_________|_________\n         |         |\n         |         |\n    %s    |    %s    |    %s    \n         |         |\n         |         |\n"

// Brukes for å konvertere int verdier til string
// Tankegangen er at verdi 0 i matrisen representerer ingenting
// verdi 1 representerer 'X'
// verdi 2 representerer 'O'

var valueStr = map[int]string{
	0: " ",
	1: "\033[1;34mX\033[0m",
	2: "\033[1;91mO\033[0m",
}

type TicTacToeGame struct {
	board        [3][3]int
	roundsPlayed int
}

func Create() TicTacToeGame {
	return TicTacToeGame{}
}

// Returnerer stringen som inten representerer som en string på brettet. Verdi 1 = "X", verdi 2 = "O"
func GetValueStr(value int) string {
	return valueStr[value]
}

// Gjør et trekk på spillet, value 1 = plasser en X, value 2 = plasser en O
func (Game *TicTacToeGame) DoMove(yPos int, xPos int, value int) error {
	if Game.roundsPlayed == 9 {
		return errors.New("The gameboard is full")
	}
	if yPos < 0 || yPos > 2 || xPos < 0 || xPos > 2 {
		return errors.New(fmt.Sprint(yPos, xPos, " er utenfor brettet!"))
	}
	if value != 1 && value != 2 {
		return errors.New(fmt.Sprint("Verdien du oppgav er ugyldig, det må være enten 1 (X), eller 2 (O)"))
	}
	if Game.board[yPos][xPos] != 0 {
		return errors.New(fmt.Sprint("Det står allerede en brikke i", yPos, " ", yPos))
	}
	Game.roundsPlayed++

	Game.board[yPos][xPos] = value

	return nil
}

// Returnerer antall trekk som har blitt gjort i spillet
func (Game TicTacToeGame) GetRoundsPlayed() int {
	return Game.roundsPlayed
}

// Returnerer true hvis brikken som representerer value har tre på rad. 1 = X, 2 = O
func (Game *TicTacToeGame) CheckWin(value int) bool {
	var matches int = 0 // matches holder styr på hvor mange like på rad spilleren har

	for i := 0; i < 3; i++ { // Sjekke om brukeren vant på den horisontale aksen
		matches = 0
		for x := 0; x < 3; x++ {
			if Game.board[i][x] != value {
				break
			}
			matches++
			if matches == 3 {
				return true // Brukeren vant på den horisontale aksen
			}
		}
	}

	for x := 0; x < 3; x++ {
		matches = 0
		for i := 0; i < 3; i++ { // Sjekke om brukeren vant på den vertikale aksen
			if Game.board[i][x] != value {
				break
			}
			matches++
			if matches == 3 {
				return true // Brukeren vant på den vertikale aksen
			}
		}
	}

	matches = 0

	for i := 0; i < 3; i++ { // Sjekke om brukeren vant på diagonal
		if Game.board[i][i] != value {
			break
		}
		matches++
		if matches == 3 {
			return true
		}
	}

	/* Sjekke om brukeren vant på diagonal fra
	øvre høyre hjørnet eller nedre venstre hjørnet
	i og med at den første for loopen ikke dekker det
	*/

	matches = 0

	for i := 0; i < 3; i++ {
		if Game.board[i][(3-i)-1] != value {
			break
		}
		matches++
		if matches == 3 {
			return true
		}
	}
	return false
}

// Returnerer true hvis spillet er uavgjort
func (Game TicTacToeGame) CheckDraw() bool {
	// Brettet har ikke blitt fylt enda
	if Game.roundsPlayed != 9 {
		return false
	}
	// Sjekke om ingen av spillerene har vunnet
	if !Game.CheckWin(1) && !Game.CheckWin(2) {
		// Det har gått 9 runder (brettet er fylt) og ingen av spillerene har vunnet
		return true
	}
	return false
}

// Printe spillbrettet i terminalen
func (Game *TicTacToeGame) Draw() {
	var data []interface{}
	// Konvertere spillbrett matrisen [3][3]int til [9]interface{}
	// Slik at vi kan passere den som varargs
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			data = append(data, valueStr[Game.board[y][x]])
		}
	}
	fmt.Println(fmt.Sprintf(brettFormat, data...))
}
