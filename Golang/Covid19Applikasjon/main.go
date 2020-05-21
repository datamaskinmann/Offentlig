/*
* @Author Sindre Fredriksen & Sven Sørensen
* @Version 05.04.2020
 */

/*
* Denne pakken står for aktivisering av programmet
* som skal søke etter og finne .pdf filer hos WHO
* og laste disse ned for videre analyse for å finne
* antall smittede i Norge dag for dag og laste dem
* opp til en webserver.
 */

package main

import (
	"./dataManagement"
	"./webServer"
)

/*
* Denne funksjonen initialiserer programmet ved å gjøre
* et kall og kjøre func.Init i pakkene dataManagement og
* webServer.
 */
func init() {

	dataManagement.Init()
	webServer.Init()

}

func main() {

}	
