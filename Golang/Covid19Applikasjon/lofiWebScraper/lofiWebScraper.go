/*
* @Author Sindre Fredriksen & Sven Sørensen
* @Version 28.04.2020
 */

/*
 * lofiWebScraper pakken har ansvar for å
 * aksessere en URL, finne, laste ned
 * .pdf filer og fjerne duplikater.
 */
package lofiWebScraper

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"../textManagement"
)

/*
* client er en global variabel av
* typen http.Client som timer out
* hvis en nedlastning tar over 6 sek.
 */
var client = http.Client{
	Timeout: 6 * time.Second,
}

// Oops - GetHtml funksjonen er skamløst plyndret fra https://siongui.github.io/2016/03/19/go-download-html-from-url/

/*
* GetHtml funksjonen tar inn en string webPage som
* parameter og returnerer en string(n) som er et
* HTML dokument.
 */
func GetHtml(webPage string) string {

	response, err := client.Get(webPage)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	n, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}
	return string(n)
}

/*
* FindFile funksjonen returnerer en slice av
* type string som inneholder .pdf linker til
* WHO sine sider.
 */
func FindFile() []string {

	var htmlData string = GetHtml("https://www.who.int/emergencies/diseases/novel-coronavirus-2019/situation-reports")

	var filer []string = textManagement.RegExp("\\/docs\\/default-source[^\\?]+", htmlData)

	filer = rmDupes(filer)

	for a, b := range filer {
		filer[a] = "https://www.who.int" + b
	}
	return filer
}

/*
* rmDupes funksjonen tar imot en slice av type
* string og returnerer en slice av type string
* uten duplikater.
 */
func rmDupes(strSlice []string) []string {
	encountered := map[string]bool{}

	res := []string{}

	for i := range strSlice {

		if !encountered[strSlice[i]] {
			encountered[strSlice[i]] = true

			res = append(res, strSlice[i])
		}
	}

	return res
}

// Oops -> DownloadFile funksjonen er hentet fra https://golangcode.com/download-a-file-from-a-url/

/*
* DownloadFile funksjonen tar i mot en string url
* og en string savePath som laster ned og lagrer en
* fil.
* Returnerer error hvis det skjer en feil hvis ikke
* returneres nil.
 */
func DownloadFile(url string, savePath string) error {

	resp, err := client.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(savePath)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}
