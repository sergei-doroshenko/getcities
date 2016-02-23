package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/sergei-doroshenko/getcities/config"
	"github.com/sergei-doroshenko/getcities/converter"
	"github.com/sergei-doroshenko/getcities/logger"
)

var cfg config.Config
var msgs = make(chan string)
var wg sync.WaitGroup

func getCities(country string) {
	defer wg.Done()

	fmt.Println("Start getting cities from ", country)
	xmlData := converter.GetXml(country, &cfg, http.Get)

	v := converter.XmlRecord{}
	err := xml.Unmarshal([]byte(xmlData), &v)

	if err != nil {
		fmt.Printf("error: %v", err)
	}

	jsonStr := converter.ConvertToJson(v)

	if jsonStr == "" {
		errJsonData, _ := json.Marshal(logger.LogRecord{Country: country, Error: "Error while getting cities"})
		errJson := string(errJsonData)
		msgs <- errJson

	} else {
		fmt.Println("Converted json: ", jsonStr)
		msgs <- jsonStr
	}
}

func main() {
	cfg = config.ReadConfig()

	countries := os.Args[1:]

	for _, country := range countries {
		wg.Add(1)
		go getCities(country)
	}

	wg.Add(1)
	go logger.LogToFile(len(countries), msgs, &wg)
	wg.Wait()

	fmt.Println("main finished")
}
