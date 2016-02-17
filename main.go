package main

import (
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
	// fmt.Printf("{country: %s, cities:[", v.NewDataSet.Table[0].Country)
	/*for _, element := range v.NewDataSet.Table {
		// element is the element from someSlice for where we are
		fmt.Print(element.City + ", ")
	}
	fmt.Print("]}\n")*/
	fmt.Println(v.NewDataSet.Table)
}
