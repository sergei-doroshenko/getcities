package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	// "sync"
	"time"
)

// const url string = "http://www.webservicex.net/globalweather.asmx/GetCitiesByCountry?CountryName="

var config Config
var msgs = make(chan string)
var done = make(chan bool)

// var wg sync.WaitGroup // 1

type Config struct {
	Url     string `json:"url"`
	Queries int    `json:"queriew"`
}

type XmlRecord struct {
	XMLName    xml.Name `xml: "string"`
	NewDataSet struct {
		XMLName xml.Name `xml: "NewDataSet"`
		Table   []struct {
			XMLName xml.Name `xml: "Table"`
			Country string   `xml:"Country"`
			City    string   `xml: City"`
		} `xml: "Table"`
	} `xml: "NewDataSet"`
}

type LogRecord struct {
	Country string   `json:"country"`
	Cities  []string `json:"cities"`
	Error   string   `json:"error"`
}

func getXml(country string) string {
	resp, err := http.Get(config.Url + country)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return ""
	}
	return html.UnescapeString(string(contents))
}

func convertToJson(xmlRecord XmlRecord) string {
	records := xmlRecord.NewDataSet.Table
	if len(records) > 0 {
		var c string = xmlRecord.NewDataSet.Table[0].Country
		var ct []string = make([]string, 0)

		for _, element := range xmlRecord.NewDataSet.Table {
			ct = append(ct, element.City)
		}

		rec := LogRecord{
			Country: c,
			Cities:  ct}

		jsonData, _ := json.Marshal(rec)
		return string(jsonData)
	}
	return ""
}

func createLogFile(filename string) {
	_, err3 := os.Create(filename)
	if err3 != nil {
		fmt.Println(err3.Error())
	} else {
		fmt.Println("Created file: ", filename)
	}
}

func appendToFile(filename, text string) {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// path/to/whatever does not exist
		fmt.Println("File not exists")
		createLogFile(filename)
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600) //0644
	if err != nil {
		// panic(err)
		fmt.Println("File doesn't exist")
		fmt.Println("Open file error: ", err.Error())

	} else {
		fmt.Println(f.Name())

		if _, err = f.WriteString(text + "\n"); err != nil {
			//panic(err)
			fmt.Println("Can't write to file")
			fmt.Println(err.Error())
		}

	}

	defer f.Close()
}

func logToFile(msgCount int) {

	current_time := time.Now().Local()
	filename := "cities-" + current_time.Format("2006-01-02") + ".log"

	for msgCount > 0 {
		text := <-msgs
		fmt.Println("appender get messege: ", text)
		appendToFile(filename, text)
		msgCount--
	}

	done <- true

}

func getCities(country string) {
	// defer wg.Done() // 3

	fmt.Println("Start getting cities from ", country)
	xmlData := getXml(country)

	// fmt.Println(xmlData)
	v := XmlRecord{}
	err := xml.Unmarshal([]byte(xmlData), &v)

	if err != nil {
		fmt.Printf("error: %v", err)
	}

	jsonStr := convertToJson(v)

	if jsonStr == "" {
		errJsonData, _ := json.Marshal(LogRecord{Country: country, Error: "Error while getting cities"})
		errJson := string(errJsonData)
		msgs <- errJson

	} else {
		fmt.Println("Converted json: ", jsonStr)
		msgs <- jsonStr
	}
}

func readConfig() Config {
	dat, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Can't read configuration.")
		os.Exit(1)
	}

	config := Config{}

	if err := json.Unmarshal(dat, &config); err != nil {
		panic(err)
	}
	fmt.Println("url: " + config.Url)
	return config
}

func main() {
	config = readConfig()

	countries := os.Args[1:]
	// countries := []string{"Ukrain", "Poland", "Russia3838"}
	// wg.Add(1) // 2

	for _, country := range countries {
		go getCities(country)
	}

	// go getCities()

	// wg.Add(1) // 2
	// go getCities()
	// json := <-msgs

	// wg.Wait() // 4
	go logToFile(len(countries))

	<-done
	fmt.Println("main finished")
}
