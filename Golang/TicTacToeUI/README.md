Dette prosjektet er ufullstendig. Det fungerer trolig kun å kjøre dersom mappen prosjektet omfatter
er i %GOPATH%/src.

Et Tic Tac Toe prosjekt relatert til faget IS-105 "Datakommunikasjon og operativsystemer". Prosjektet
er dessverre ufullstendig ettersom vi fikk en uventet Covid-19 relatert oppgave som krevde mer tid enn
forventet og gav oss mindre til tid til å fokusere på dette prosjektet.

Prosjektet mangler "polering", det er for eksempel mulig å "sette seg fast" i en meny, altså å gå til
en meny, men ikke komme seg tilbake igjen og å sjekke om det oppstår feil et par plasser for å hindre
at menyen går til en uønsket plass med manglende data.

Men det som derimot fungerer er 1v1 lokalt og 1v1 på nett, men 2v2 på nett fungerer for en eller annen grunn
ikke, å debugge dette problemet fikk vi ikke tid til å gjøre.

Prosjektets nettkode baserer seg på et nettkode bibliotek jeg skrev, som du kan finne under Golang/nettkode.
GUIen baserer seg på et bibliotek som heter "GOSCITER" -> https://github.com/sciter-sdk/go-sciter

Bruksinstruksjoner
Dersom du vil teste online modusen ->
1: Naviger i en konsoll til "TicTacToe UI/Tjener"
2: Utfør kommandoen "go run TicTacToeServer.go"
3: Naviger i en ny konsoll til "TicTacToe UI/Klient"
4: Utfør kommandoen "go run main.go" og åpne så mange klienter du ønsker
5: På en av klientene, naviger i brukergrensesnittet til
	Play Online -> *skriv et brukernavn*, trykk connect ->
	Host Tournament -> (kun 2 spillere fungerer dessverre) -> Create Lobby
6: Klient nummer 2
	Play Online -> *skriv et brukernavn*, trykk connect ->
	Join Tournament -> skriv inn Lobby ID som er på Klient nummer 1 sin skjerm ->
	trykk join. NB: Vi planla å legge til funksjonalitet for å kopiere Lobby ID på et trykk, men
	vi kom ikke så langt og dermed må Lobby ID'en skrives manuelt. Det kan av den årsak være vanskelig
	å skille med karakterer slik som liten l og stor i.
	
	