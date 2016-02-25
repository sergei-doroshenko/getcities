package test

import (
	"encoding/xml"
	"fmt"
	"sync"
	"testing"

	"github.com/sergei-doroshenko/getcities/converter"
)

const countryName string = "belarus"

type MockDataSource struct {
}

var belarusCitiesXml string = `
<?xml version="1.0" encoding="utf-8"?>
<string xmlns="http://www.webserviceX.NET">
	<NewDataSet>
	  <Table>
	    <Country>Belarus</Country>
	    <City>Riga Airport</City>
	  </Table>
	  <Table>
	    <Country>Belarus</Country>
	    <City>Vitebsk</City>
	  </Table>
	  <Table>
	    <Country>Belarus</Country>
	    <City>Minsk</City>
	  </Table>
	</NewDataSet>
</string>
`
var belarusCitiesJson = `{"country":"Belarus","cities":["Riga Airport","Vitebsk","Minsk"],"error":""}`

func (s MockDataSource) GetData(country string) ([]byte, error) {
	var m map[string][]byte
	m = make(map[string][]byte)
	m[countryName] = []byte(belarusCitiesXml)
	return m[countryName], nil
}

func TestGetData(t *testing.T) {
	fmt.Println("TestGetData...")
	mock := MockDataSource{}
	data, _ := mock.GetData(countryName)

	xmlStr := string(data)
	if xmlStr != belarusCitiesXml {
		t.Error("Expected "+belarusCitiesXml+", got ", xmlStr)
	}

}

func TestGetXml(t *testing.T) {
	fmt.Println("TestGetXml...")
	dtsource := MockDataSource{}
	conv := converter.Converter{DataSource: dtsource}

	xmlStr, _ := conv.GetXml(countryName)

	if xmlStr != belarusCitiesXml {
		t.Error("Expected "+belarusCitiesXml+", got ", xmlStr)
	}

}

func TestGetXmlRecord(t *testing.T) {
	fmt.Println("TestGetXmlRecord...")
	xmlRecord0 := converter.XmlRecord{}
	xml.Unmarshal([]byte(belarusCitiesXml), &xmlRecord0)

	dtsource := MockDataSource{}
	conv := converter.Converter{DataSource: dtsource}

	xmlRecord1, _ := conv.GetXmlRecord(countryName)
	if !xmlRecord0.Equals(xmlRecord1) {
		t.Error("XmlRecords mismatch.")
	}
}

func TestConvertToJson(t *testing.T) {
	fmt.Println("TestConvertToJson...")
	dtsource := MockDataSource{}
	conv := converter.Converter{DataSource: dtsource}
	xmlRecord, _ := conv.GetXmlRecord(countryName)
	jsonStr, _ := conv.ConvertToJson(xmlRecord)
	if string(jsonStr) != belarusCitiesJson {
		t.Error("Expected "+belarusCitiesJson+", got ", jsonStr)
	}
}

func TestGetCities(t *testing.T) {
	fmt.Println("TestGetCities...")
	dtsource := MockDataSource{}
	conv := converter.Converter{DataSource: dtsource}

	var msgs = make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go conv.GetCities(countryName, msgs, &wg)

	message := <-msgs
	fmt.Println("Message: ", message)
	if message != belarusCitiesJson {
		t.Error("Expected "+belarusCitiesJson+", got ", message)
	}
	wg.Wait()
}
