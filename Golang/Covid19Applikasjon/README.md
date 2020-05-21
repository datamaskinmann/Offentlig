Dette prosjektet ble laget i henhold til en oppgave vi fikk i IS-105 (Datakommunikasjon og operativsystem).
Oppgaven gikk ut på å lage en Covid 19 oversikt for Norge. Per 1. Mai 2020 er applikasjonen dessverre utdatert
ettersom WHO.INT (Siden vi henter og leser rapportene fra) endret rapportstrukturen deres som fører til at
vi måtte ha endret store deler av koden for å få den til å fungere igjen.

PDFTOTEXT er hentet fra https://www.xpdfreader.com/download.html

Dersom du ikke bruker Windows og skulle allikevel ønske å kjøre programmet, gå til linken gitt ovenfor og
last ned XPDFREADER for ditt operativsystem, gå så i "dataManagement.go" og erstatt "./pdfreadertotext.exe" til
ditt format.

Instruksjoner for kjøring av program

1: Naviger med konsoll til samme directory som "main.go"
2: "go run main.go"

Det er alt som skal til for å kjøre programmet

3: Gå til din foretrukkende webbrowser. Foretrekkelig Google Chrome, vi har opplevd at applikasjonen
ikke fungerer i Firefox på grunn av en feil ved lasting av Google charts, noe som vi avhenger av
for å generere grafer.
4: Naviger til localhost:8080