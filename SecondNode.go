package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/drone/routes"
)

var cacheMap map[string]string

type MapValueResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MapAllValuesResponse struct {
	AllValues []MapValueResponse `json:"allValues"`
}

func main() {

	mux := routes.New()
	mux.Put("/keys/:keyID/:value", PutValue)

	mux.Get("/keys/:keyID", RetrieveSingleValue)
	mux.Get("/keys", RetrieveAllValues)

	http.Handle("/", mux)
	http.ListenAndServe(":3001", nil)
}

func PutValue(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get(":keyID")
	value := r.URL.Query().Get(":value")

	if cacheMap == nil {
		cacheMap = make(map[string]string)
	}
	cacheMap[key] = value
	fmt.Println(cacheMap)
	w.WriteHeader(http.StatusOK)

}

func RetrieveSingleValue(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get(":keyID")
	var output MapValueResponse
	isKeyFound := false
	fmt.Println("Key to search for is : ", key)

	if cacheMap == nil {
		fmt.Println("Making new map")
		cacheMap = make(map[string]string)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Cache is empty"}`))
	} else {
		fmt.Println("Entered else part")
		for mapKey, mapVal := range cacheMap {
			fmt.Println(key, "  ", mapKey)
			if key == mapKey {
				fmt.Println("Entered if part")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				output.Key = mapKey
				output.Value = mapVal
				outputJSON, err := json.Marshal(output)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}
				w.Write(outputJSON)
				isKeyFound = true
				break
			}
		}
		if !isKeyFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "value not found"}`))
		}
		fmt.Println(isKeyFound)

	}

}

func RetrieveAllValues(w http.ResponseWriter, r *http.Request) {
	var output MapAllValuesResponse
	var currValue MapValueResponse
	for mapKey, mapVal := range cacheMap {
		currValue.Key = mapKey
		currValue.Value = mapVal
		output.AllValues = append(output.AllValues, currValue)
	}
	w.WriteHeader(http.StatusOK)
	outputJSON, _ := json.Marshal(output)
	w.Header().Set("Content-Type", "application/json")
	w.Write(outputJSON)
}
