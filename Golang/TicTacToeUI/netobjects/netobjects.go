package netobjects

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var networkableObjects = map[string]reflect.Type{
	"netobjects.SetID":              reflect.TypeOf(SetID{}),
	"netobjects.ClientDisconnected": reflect.TypeOf(ClientDisconnected{}),
	"netobjects.ClientConnected":    reflect.TypeOf(ClientConnected{}),
}

/* Pga golang ikke har attributes og ingen måte å reflektere
alle structsa i en package må jeg lage en metode for å manuelt
finne riktig objekt type ved å ha et map hvor nøkkelen
er navnet på reflect.Type i string format og verdien er
reflect.Type */

// NetObject er en form for metadata om objekter som sendes over nett og brukes kun av serveren eller klienten under sending og lesing av objekter

type NetObject struct {
	Unmarshaltype string // denne brukes til å finne reflect.Type i networkableobjects mappet
	Data          string // den faktiske dataen til objektet som sendes over nettet som må unmarshalles
}

type SetID struct {
	ID int
}

type ClientDisconnected struct {
	ID int
}

type ClientConnected struct {
	ID int
}

// GetNetworkableObjects returnerer en reflect.Type array som representerer alle tillagte NetworkableObjects
func GetNetworkableObjects() []reflect.Type {
	// Opprette et array av typen reflect.Type og på lenge med networkableObjects mappet
	var ret []reflect.Type = make([]reflect.Type, len(networkableObjects))
	// Index int
	var i int
	// Iterere gjennom networkableObjects mappet
	for _, a := range networkableObjects {
		// Legge 'a' som er value i networkableObjects
		// entreet vi for øyeblikket itererer på i ret arrayet
		ret[i] = a
		i++
	}
	return ret
}

// Formatteres slik: AddNetworkableObject(reflect.TypeOf(DinStruct{}))
func AddNetworkableObject(Type reflect.Type) {
	networkableObjects[Type.String()] = Type
}

// (ignorer) Finner tilsvarende reflect.Type i networkableObjects mappet
func findObjectType(TypeName string) reflect.Type {
	return networkableObjects[TypeName]
}

// (ignorer) Sjekk om et objekt er networkable
func IsNetworkableObject(Object interface{}) bool {
	if networkableObjects[reflect.TypeOf(Object).String()] == nil {
		fmt.Println(fmt.Sprint("Object '", reflect.TypeOf(Object).String(), "' is not a Networkable object, please add it using netobjects.AddNetworkableObject()"))
		return false
	}
	return true
}

func GetObject(data []byte) interface{} {
	// De tre siste bytesa av networkobject []byte
	// er informasjon om sender og mottaker av objektet
	rawData := data[:len(data)-3]

	var netObj NetObject

	// Deserialisere data vi mottok i parameteret til NetObject struct
	json.Unmarshal(rawData, &netObj)

	// Finnes deserialiseringstypen til denne strukturen i networkableObjects mappet vårt?
	if findObjectType(netObj.Unmarshaltype) == nil {
		// Nei -> panic
		fmt.Println(fmt.Sprint("Type", netObj.Unmarshaltype, "is not listed as a networkable object"))
		return nil
	}

	// Instansiere ved hjelp av refleksjon et struct av typen spesifisert
	// i NetObjectet vi nettopp deserialiserte
	xObj := reflect.New(findObjectType(netObj.Unmarshaltype)).Interface()

	// Deserialisere dataen inni netObj structet vi nettopp deserialiserte
	// til dets tiltenkte struktur
	json.Unmarshal([]byte(netObj.Data), &xObj)

	return xObj
}
