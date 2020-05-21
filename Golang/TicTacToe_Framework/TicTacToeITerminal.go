package main

import (
	"fmt"

	"./input"
	"./tictactoegame"
)

func main() {
	Game := tictactoegame.Create() // Opprette en tictactoe instans
	Game.Draw()                    // Tegne brettet

	var turn int = 1

	fmt.Println("For å plassere en brikke, skriv (Y-Akse X-Akse) eksempel '0 0'")
	for {
		fmt.Println("Det er", tictactoegame.GetValueStr(turn), "sin tur")
		args, convError := input.ReadIntSlice() // Lese string fra terminalen og konvertere det til int slice
		if convError != nil {
			fmt.Println("Det har oppstått en feil. Vennligst oppgi tall som argumenter")
			continue // Det har oppstått en feil, gå til toppen av løkken
		}
		moveError := Game.DoMove(args[0], args[1], turn)
		if moveError != nil { // Er trekket gyldig?
			fmt.Println(moveError)
			continue
		}
		Game.Draw()
		if Game.CheckDraw() {
			fmt.Println("Spillet ble uavgjort!")
			break // bryte løkken
		}
		if Game.CheckWin(turn) {
			fmt.Println(tictactoegame.GetValueStr(turn), "vant spillet!")
			break
		}
		switch turn {
		case 1:
			turn = 2
			break
		case 2:
			turn = 1
			break
		}
	}
	fmt.Println("Trykk enter for å ende spillet")
	input.ReadString()
}