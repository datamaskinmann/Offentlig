package tictactoenet

// CreateLobby brukes både av tjener og klient.
// av klient til å 'forespørre' å opprette en lobby
// av tjener til å signalisere at forespørselen er godkjent
// tjeneren setter LobbyID
type CreateLobby struct {
	LobbySize int
	LobbyID   string
}

// Move sendes fra en klient til en annen
// for å melde om et trekk gjort
type Move struct {
	XPos    int
	YPos    int
	Turnval int
}

// SetName brukes av klient til å navnsette seg selv hos tjener
type SetName struct {
	Name string
}

// PlayerJoinLobby blir sendt fra en klient til en/flere andre klienter for å signalisere at han/hun har
// koblet til lobbyen
type PlayerJoinLobby struct {
	Name string
}

// GetLobby kan bli sendt fra klient til tjener for å motta en oversikt over alle klientene som er koblet til
// tjeneren
type GetLobby struct {
	LobbyID string
	Names   []string
	IDs     []int
	Error   error
}

// JoinLobby sendes fra klient til tjener for å joine en lobby. Eksisterer ikke lobbyen, sendes en error tilbake fra tjeneren.
type JoinLobby struct {
	LobbyID string
	Error   error
}

// Match sendes fra tjener til klient for å symbolisere at spilleren skal spille mot spilleren med ID i OpponentID feltet.
// Match sendes tilbake fra klient til tjener for å symbolisere at matchen er ferdig.
// Result er hvor mange poeng spilleren har fått. 0 poeng = tap. 1 poeng = uavgjort. 2 poeng = vunnet.
// Turnval er turverdien til spilleren som mottar objektet. 1 = X. 2 = O.
type Match struct {
	OpponentID   int
	OpponentName string
	Result       int
	Turnval      int
}

// TournamentFinished structet blir sendt fra tjener til alle klienter når en turnering er ferdig.
// results mappet vil inneholde spillerens navn som key, og spillerens poeng som value.
type TournamentFinished struct {
	Results map[string]int
}
