// Denne packagen har som hensikt å implementere en form for
// events slik vi kjenner de i java

package event

import (
	"errors"
	"fmt"
	"reflect"
)

type Event struct {
	subscribers []interface{} // Metodene vi vil kalle når event.Raise() blir kalt
}

// Length returns the amount of subscribers the event has
func (e Event) Length() int {
	return len(e.subscribers)
}

// Add a subscriber func to the event, the function must have eventArgs interface{} and sender interface{} as params
func (e *Event) Add(function interface{}) error {
	// Er function parameteret faktisk en funksjon?
	if reflect.TypeOf(function).Kind() != reflect.Func {
		// Nei, det er det ikke, returnerer error
		return errors.New(fmt.Sprint(reflect.TypeOf(function), " is not a func"))
	}

	// Inneholder e.subscribers slicen allerede denne funksjonen?
	if contains(e.subscribers, function) {
		// Ja, det gjør det, returnerer error
		return errors.New(fmt.Sprint(reflect.TypeOf(function), " is already subscribed to ", e))
	}

	// Tilføye funksjonen i e.subscribers slicen
	e.subscribers = append(e.subscribers, function)
	// Ingen error, returner nill
	return nil
}

// Contains sjekker om slice s inneholder element e
func contains(s []interface{}, e interface{}) bool {
	// Iterere gjennom slice s
	for _, a := range s {
		// Er 'a' lik element e?
		if a == e {
			// Ja -> returner true
			return true
		}
	}
	// Vi har slicen uten å returnere true -> returner false
	return false
}

// remove fjerner element e fra slice s og returnerer slicen uten element e og en error
func remove(s []interface{}, e interface{}) ([]interface{}, error) {
	// Finne indexen av element e i slice s
	i, iErr := indexOf(s, e)
	// Oppstod det en feil?
	if iErr != nil {
		// Ja, det gjorde det returner nil og en error
		return nil, iErr
	}

	// Flytte det ønsket slettede elementet bakerst i slicen
	s[len(s)-1] = s[i]
	// Sette dette elementet til nil (nullpeker i minnet)
	s[len(s)-1] = nil

	// Returnere slicen fra element 0, til slicens lengde -1
	// (Alt slicen inneholdt fra før av minus elementet vi nettopp fjerna)
	// Samt ingen error
	return s[:len(s)-1], nil
}

// indexOf finner indexen av element e i slice s og returnerer
// en int som representerer indexen av elementet og en error
func indexOf(s []interface{}, e interface{}) (int, error) {
	// Iterere gjennom slice s
	for i := 0; i < len(s); i++ {
		// Er elementet vi for øyeblikket itererer på
		// Element 'e' spesifisert i parameteret?
		if s[i] == e {
			// Ja -> returner indexen og nilerror
			return i, nil
		}
	}
	// Vi har iterert gjennom hele slicen og element 'e' spesifisert i parameteret
	// finnes ikke i slicen, returnerer -1 og en feilmelding
	return -1, errors.New("The slice does not contain the specified element")
}

// Remove removes a subscriber func from the event, the function passed will take the eventArgs as its first parameter, and the sender interface{} as the second parameter
func (e *Event) Remove(function interface{}) error {
	// Fjerne funksjonen spesifisert i parameteret fra e.subscribers
	slice, removeError := remove(e.subscribers, function)
	// Oppstod det en feil under fjerningen av elementet?
	if removeError != nil {
		// Ja -> returner feilen
		return removeError
	}
	// Sette e.subscribers til slicen som remove funksjonen returnerte
	e.subscribers = slice
	return nil
}

// RaiseAsync rasies the event on a new goroutine. The parameters eventArgs and sender can be a custom struct, or new(interface{}) for the equivalent of null in Java
func (e *Event) RaiseAsync(eventArgs interface{}, sender interface{}) {
	// Iterere gjennom alle funksjonene som e.subscribers inneholder
	for _, a := range e.subscribers {
		// Opprette en goroutine som ved hjelp av refleksjon kaller på funksjonene som e.subscribers inneholder
		// og passerer inn eventArgs og sender som parametere
		go reflect.ValueOf(a).
			Call([]reflect.Value{reflect.ValueOf(eventArgs), reflect.ValueOf(sender)})
	}
}

// Raise rasies the event on the goroutine of the caller. The parameters eventArgs and sender can be a custom struct, or new(interface{}) for the equivalent of null in Java
func (e *Event) Raise(eventArgs interface{}, sender interface{}) {
	// Iterere gjennom alle funksjonene som e.subscribers inneholder
	for _, a := range e.subscribers {
		// Opprette en goroutine som ved hjelp av refleksjon kaller på funksjonene som e.subscribers inneholder
		// og passerer inn eventArgs og sender som parametere
		reflect.ValueOf(a).
			Call([]reflect.Value{reflect.ValueOf(eventArgs), reflect.ValueOf(sender)})
	}
}
