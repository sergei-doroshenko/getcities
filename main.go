package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
)

const url string = "http://www.webservicex.net/globalweather.asmx/GetCitiesByCountry?CountryName="

type Table struct {
	XMLName xml.Name `xml: "Table"`
	Country string   `xml:"Country"`
	City    string   `xml: City"`
}

type NewDataSet struct {
	XMLName xml.Name `xml: "NewDataSet"`
	Table   []Table  `xml: "Table"`
}

type Root struct {
	XMLName    xml.Name   `xml: "string"`
	NewDataSet NewDataSet `xml: "NewDataSet"`
}

/*
type LogRecord struct {
	Country string   `json:"country"`
	Cities  []string `json:"cities"`
}
*/
type LogRecord struct {
	Country string
	Cities  []string
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

func createLogFile(filename string) *os.File {
	f, err3 := os.Create("test.log")
	if err3 != nil {
		fmt.Println(err3.Error())
	}

	fmt.Println("Created file: ", f.Name)
	return f
}

func appendToFile(filename, text string) {
	/*
		var logfile *os.File
		_, err4 := os.Stat(filename)
		if err4 != nil {
			fmt.Println(err4.Error())
			fmt.Println("Log file: " + filename + " doesn't exists.")
			logfile = createLogFile(filename)
			logfile.Chmod(0644)
			defer logfile.Close()
		}
	*/
	/*
		f := createLogFile(filename)
		defer f.Close()
	*/
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		// panic(err)
		fmt.Println("File doesn't exist")
		fmt.Println(err.Error())

	}

	defer f.Close()
	fmt.Println(f.Name())

	if _, err = f.WriteString(text); err != nil {
		//panic(err)
		fmt.Println(err.Error())
	}

}

func main() {
	fmt.Println("Hello, start process...")
	xmlData := getXml("Ukrain")
	fmt.Println(xmlData)
	v := Root{}

	err1 := xml.Unmarshal([]byte(xmlData), &v)

	if err1 != nil {
		fmt.Printf("error: %v", err1)
		return
	}

	fmt.Println("XMLName: ", v.XMLName)
	fmt.Println(v.NewDataSet.Table)
	res1D := LogRecord{
		Country: "USA",
		Cities:  []string{"New York", "Boston", "Dallas"}}

	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	logRecord := LogRecord{Country: v.NewDataSet.Table[0].Country, Cities: make([]string, 0)}
	for _, element := range v.NewDataSet.Table {
		// element is the element from someSlice for where we are
		logRecord.Cities = append(logRecord.Cities, element.City)
	}
	fmt.Println(logRecord)
	/*
		logRecord := LogRecord{
			Country: "GreateBritain",
			Cities:  []string{"Glazgo", "London", "Birmengerm"}}
		fmt.Println(logRecord)
	*/
	//logStr, _ := json.Marshal(LogRecord)
	//fmt.Println(string(logStr))

	//	appendToFile("./test.log", logStr)

	//ioutil.WriteFile("log.txt", []byte(xmlData), 0644)
}
