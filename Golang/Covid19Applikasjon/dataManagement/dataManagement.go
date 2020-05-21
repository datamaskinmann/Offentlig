/*
* @Author Sindre Fredriksen & Sven Sørensen
* @Version 28.04.2020
 */

/*
* dataManagement Pakken har ansvar for å sjekke etter
* og håndtere ny data før den blir lastet opp til
* webserveren.
*
 */
package dataManagement

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"

	"../lofiWebScraper"
	"../textManagement"
	"../threadManagement"
	"../timeManagement"
)

/*
* covidDataMap variabelen ernen global variabel og
* oppretter et map av type string.
 */
var covidDataMap map[string]covidData

/*
* covidData struct er utformingen på hvordan dataen
* skal lagres og fremstilles.
 */
type covidData struct {
	Country     string
	TotalCases  int
	TotalDeaths int
	FromUrl     string
}

/*
* pdfReader funksjonen tar inn en string filePath,
* konverterer filen til tekst og returnerer en string
* som er tekst innholdet i filen.
 */
func pdfReader(filePath string) string {

	cmd := exec.Command("./pdftotext.exe", filePath, "-")

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	return out.String()
}

/*
* findCountryTotalCases tar inn en string data og en
* string countryName, og returnerer int som er
* antallet smittede i landet hvis ikke error.
 */
func findCountryData(data string, countryName string) ([]int, error) {

	var expression string = fmt.Sprintf("(?s)(%s)[^a-zA-z0-9]*([0-9]+)[^0-9]+([0-9]+)[^0-9]+([0-9]+)", countryName)

	var regexReturn [][]string = textManagement.RegExpSubMatch(expression, data)

	if len(regexReturn) == 0 {

		return []int{}, errors.New(fmt.Sprintf("Could not find any data for country: (%s) in your dataset", countryName))
	}

	totalCases, _ := strconv.Atoi(regexReturn[0][2])
	totalDeaths, _ := strconv.Atoi(regexReturn[0][4])

	return []int{totalCases, totalDeaths}, nil
}

/*
* checkUpdate funksjonen sjekker etter nye data.
* Hvis den finner ny data, blir dataen lastet
* og behandlet for presentasjon av webserveren.
 */
func checkUpdate() {

	for {
		fmt.Println("Sjekker om det har kommet nye PDFer fra WHO")
		pdfLinker := lofiWebScraper.FindFile()

		var untrackedFiles []string

		for _, a := range pdfLinker {

			if _, ok := covidDataMap[a]; !ok {

				untrackedFiles = append(untrackedFiles, a)
			}
		}
		var alteredData bool
		if len(untrackedFiles) > 0 {
			fmt.Println("Oppdaget ny PDF fra WHO!")

			threadManagement.MainSync.Add(1)

			for _, a := range untrackedFiles {

				fmt.Println("Forsøker å laste ned: ", a)

				dlErr := lofiWebScraper.DownloadFile(a, "./temp.pdf")

				if dlErr != nil {
					fmt.Println("Greide ikke å laste ned filen... :(")
					continue
				}
				data := pdfReader("./temp.pdf")

				countryInfo, err := findCountryData(data, "Norway")

				if err != nil {
					fmt.Println("Error while reading data...", err.Error())
					continue
				}

				covidDataMap[a] = covidData{
					Country:     "Norway",
					TotalCases:  countryInfo[0],
					TotalDeaths: countryInfo[1],
					FromUrl:     a,
				}
				alteredData = true
			}

			if alteredData {
				b, _ := json.Marshal(covidDataMap)

				writeFile("./covidNorge.json", b)
			}

			threadManagement.MainSync.Done()
		}

		time.Sleep(time.Until(timeManagement.TruncateToNearestHour()))
	}

}

/*
* readFile funksjonen tar inn string path og returnerer
* en byte array som består av filens innhold hvis
* ikke error.
 */
func readFile(path string) ([]byte, error) {
	fmt.Println("Forsøker å lese fil -", path)

	data, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("Greide ikke å lese filen:", err)

		return nil, err
	}
	return data, nil
}

/*
* getCovidData funksjonen åpner ./covidNorge.json og konverterer
* dataen om til og returnerer et map av type string
* formatert til struct covidData. Hvis den ikke leses
* returneres tomt map.
 */
func getCovidData() map[string]covidData {

	dataFraFil, err := readFile("./covidNorge.json")

	if err != nil {

		return map[string]covidData{}

	}
	var data map[string]covidData
	jsonErr := json.Unmarshal(dataFraFil, &data)

	if jsonErr != nil {
		return map[string]covidData{}
	}
	return data
}

/*
* Init funksjonen starter ny go routine for checkUpdate (ny tråd).
 */
func Init() {

	covidDataMap = getCovidData()
	go checkUpdate()

}

/*
* writeFile funksjonen tar inn en string path og en array byte
* og oppretter en .json fil som innholdet i arrayen blir
* skrevet til hvis ikke error, ellers returnerer den nil.
 */
func writeFile(path string, data []byte) error {

	fmt.Println("Oppretter JSON fil...")
	f, _ := os.Create("./covidNorge.json")
	f.Close()

	err := ioutil.WriteFile("./covidNorge.json", data, 0644)

	if err != nil {

		return err

	}
	return nil
}

/*
* fileExists funksjonen tar inn en string path og sjekker
* om filen eksisterer, hvis ja returneres boolean true,
*  hvis ikke returneres boolean false.
 */
func fileExists(path string) bool {

	if _, err := os.Stat(path); os.IsNotExist(err) {

		return false

	}

	return true
}

/*
* GetCovidMapValues funksjonen tar mappet som er buffret
* i minnet og returnerer en array av verdiene (ignorerer keys).*
 */
func GetCovidMapValues() []covidData {

	var retArr []covidData
	var keys []string

	for k, _ := range covidDataMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {

		retArr = append(retArr, covidDataMap[k])

	}

	return retArr

}
