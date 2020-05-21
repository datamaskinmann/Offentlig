/*
* @Author Sindre Fredriksen & Sven Sørensen
* @Version 28.04.2020
 */

/*
* textManagement pakken har som ansvar og håndtere all
* behandling og analyse av tekst.
 */
package textManagement

import (
	"regexp"
)

/*
* RegExp tar inn regular expression string og en tekst
* string som parameter som den utfører en regular
* expression analyse på.
* Returnerer en slice av type string som inneholder
* alle tilfeller som passer mot regExp uttrykket.
 */
func RegExp(regExp string, tekst string) []string {

	r, _ := regexp.Compile(regExp)

	res := r.FindAllString(tekst, -1)

	return res

}

/*
* RegExpSubMatch tar inn to strings og returnerer en array array
* av typen string, funksjonen sjekker og deler opp substrings som
* matcher regular expression mønsteret mot teksten som blir inputtet.
 */
func RegExpSubMatch(regExp string, tekst string) [][]string {

	r, _ := regexp.Compile(regExp)

	res := r.FindAllStringSubmatch(tekst, -1)

	return res

}
