/*
* @Author Sindre Fredriksen & Sven Sørensen
* @Version 28.04.2020
 */

/*
* webServer pakken starter en webserver og serverer
* data fra maskinen som serveren er hostet på til
* klienten.
 */
package webServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"../dataManagement"
	"../threadManagement"
	"../timeManagement"
)

/*
* templ variabelet er en tom referanse til et Template
* struct i golangs html/template pakke.
 */
var templ *template.Template

/*
* Init funksjonen setter opp en file server og denne file
* serveren kan aksessere filer i ./tmpfiles directoryen
* på maskinen som hoster webserveren. På denne måten kan
* bruker agenter gjøre get requests for å få servert disse
* filene. Init definerer også at indexhandler funksjonen
* skal kalles dersom en bruker agent ønsker å aksessere
* rooten til webserveren.
* Laster også templates fra ./templated directoryen som
* ligger lokalt på host maskinen.
 */
func Init() {

	templ = template.Must(template.ParseGlob("./templates/*"))

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/getCovidNorgeJson", getCovidNorgeJson)

	http.Handle("/tmpfiles/",

		http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("./tmpfiles"))))

	fmt.Println("Forsøker å starte en webserver på port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {

		panic(err)

	}

	fmt.Println("Startet webserver på port 8080 suksessfullt!")

}

/*
* indexHandler funksjonen kalles når en klient ønsker
* å aksessere rooten til webserveren.
 */
func indexHandler(w http.ResponseWriter, r *http.Request) {

	threadManagement.MainSync.Wait()

	templ.ExecuteTemplate(w, "index.html", nil)
}

/*
* getCovidNorgeJson funksjonen tar imot en http request og
* skriver verdiene fra dataManagement.GetCovidMapValues til
* en .json fil og returnerer den som en http response.
 */
func getCovidNorgeJson(w http.ResponseWriter, r *http.Request) {

	threadManagement.MainSync.Wait()

	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Expires", timeManagement.TruncateToNearestHour().Format(http.TimeFormat))

	json, _ := json.Marshal(dataManagement.GetCovidMapValues())

	w.Header().Set("Status", "200")

	w.Write(json)

}
