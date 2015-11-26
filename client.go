package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cznic/sortutil"
	"github.com/drone/routes"

	"github.com/spaolacci/murmur3"
)

type MapValueResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var nodeMap map[uint64]string
var keys sortutil.Uint64Slice

func main() {
	nodeMap = make(map[uint64]string)
	FirstNode := "http://localhost:3000/"
	SecondNode := "http://localhost:3001/"
	ThirdNode := "http://localhost:3002/"

	//Sort the map

	keys = append(keys, murmur3.Sum64([]byte(FirstNode)))
	keys = append(keys, murmur3.Sum64([]byte(SecondNode)))
	keys = append(keys, murmur3.Sum64([]byte(ThirdNode)))

	keys.Sort()
	fmt.Println("Keys array is : ", keys)

	for _, element := range keys {
		switch element {

		case murmur3.Sum64([]byte(FirstNode)):
			nodeMap[element] = FirstNode
		case murmur3.Sum64([]byte(SecondNode)):
			nodeMap[element] = SecondNode
		case murmur3.Sum64([]byte(ThirdNode)):
			nodeMap[element] = ThirdNode
		}

	}

	mux := routes.New()
	mux.Put("/keys/:keyID/:value", PutValue)
	mux.Get("/keys/:keyID", RetrieveValue)

	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)

}

func getNode(key string) string {
	keyHash := murmur3.Sum64([]byte(key))
	var returnIndex = len(keys) - 1
	for index, element := range keys {
		if keyHash < element {
			if index > 0 {
				returnIndex = index - 1
			}
			break

		}
	}
	return nodeMap[keys[returnIndex]]

}

func PutValue(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Entered insertvalue")
	key := r.URL.Query().Get(":keyID")
	value := r.URL.Query().Get(":value")
	address := getNode(key)
	var Mybuffer bytes.Buffer
	Mybuffer.WriteString(address)
	Mybuffer.WriteString("keys")
	Mybuffer.WriteString("/")
	Mybuffer.WriteString(key)
	Mybuffer.WriteString("/")
	Mybuffer.WriteString(value)

	//Create a request
	req, err := http.NewRequest("PUT", Mybuffer.String(), nil)
	if err != nil {
		fmt.Println("error: body, _ := ioutil.ReadAll(resp.Body) -- line 592")
		panic(err)
	}
	client1 := &http.Client{}
	resp, err := client1.Do(req)
	if err != nil {
		fmt.Println("error: Unable to submit request")
		panic(err)
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Data was not persisted")
		panic(err.Error())
	} else {
		fmt.Println("Key : ", key, " || Value : ", value, " --- added to node : ", address)
	}

}

func RetrieveValue(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get(":keyID")
	address := getNode(key)
	var Mybuffer bytes.Buffer
	Mybuffer.WriteString(address)
	Mybuffer.WriteString("keys")
	Mybuffer.WriteString("/")
	Mybuffer.WriteString(key)

	//Create a request
	req, err := http.NewRequest("GET", Mybuffer.String(), nil)
	if err != nil {
		panic(err)
	}
	client1 := &http.Client{}
	resp, err := client1.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var MyResponse MapValueResponse

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		err = json.Unmarshal(body, &MyResponse)
		if err != nil {
			panic(err.Error())
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		MyoutputJSON, err := json.Marshal(MyResponse)
		if err != nil {
			w.Write([]byte(`{    "error": "Unable to marshal response.`))
			panic(err.Error())
		}
		fmt.Println("Retrieved from node : ", address, " for key : ", key)
		w.Write(MyoutputJSON)
	} else {
		w.WriteHeader(resp.StatusCode)
	}

}
