package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const url string = "http://www.webservicex.net/globalweather.asmx/GetCitiesByCountry?CountryName="

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
}

func getXml(country string) string {
	resp, err := http.Get(url + country)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	return html.UnescapeString(string(contents))
}

func convertToJson(xmlRecord XmlRecord) string {
	records := xmlRecord.NewDataSet.Table
	if len(records) > 0 {
		var c string = xmlRecord.NewDataSet.Table[0].Country

		var ct []string = make([]string, 0)

		for _, element := range xmlRecord.NewDataSet.Table {
			// element is the element from someSlice for where we are
			ct = append(ct, element.City)
		}

		res2D := LogRecord{
			Country: c,
			Cities:  ct}

		jsonData, _ := json.Marshal(res2D)
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

func logToFile(text string) {
	current_time := time.Now().Local()
	filename := "cities-" + current_time.Format("2006-01-02") + ".log"
	appendToFile(filename, text)
}

func main() {
	fmt.Println("Hello, start process...")
	xmlData := getXml("Ukrain")
	// fmt.Println(xmlData)

	v := XmlRecord{}
	err1 := xml.Unmarshal([]byte(xmlData), &v)

	if err1 != nil {
		fmt.Printf("error: %v", err1)
		return
	}

	json := convertToJson(v)
	// fmt.Println("Converted json: ", json)

	logToFile(json)
}
