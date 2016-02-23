package converter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"

	"github.com/sergei-doroshenko/getcities/config"
	"github.com/sergei-doroshenko/getcities/logger"
)

/*type XmlRecord struct {
	XMLName    xml.Name `xml: "string"`
	NewDataSet struct {
		XMLName xml.Name `xml: "NewDataSet"`
		Table   []struct {
			XMLName xml.Name `xml: "Table"`
			Country string   `xml:"Country"`
			City    string   `xml: City"`
		} `xml: "Table"`
	} `xml: "NewDataSet"`
}*/

type Entry struct {
	XMLName xml.Name `xml: "Table"`
	Country string   `xml:"Country"`
	City    string   `xml: City"`
}

type DataSet struct {
	XMLName xml.Name `xml: "NewDataSet"`
	Table   []Entry
}

type XmlRecord struct {
	XMLName    xml.Name `xml: "string"`
	NewDataSet DataSet
}

type Getable func(url string) (resp *http.Response, err error)

func GetXml(country string, conf *config.Config, fn Getable) string {
	resp, err := fn(conf.Url + country)

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

func ConvertToJson(xmlRecord XmlRecord) string {
	records := xmlRecord.NewDataSet.Table
	if len(records) > 0 {
		var c string = xmlRecord.NewDataSet.Table[0].Country
		var ct []string = make([]string, 0)

		for _, element := range xmlRecord.NewDataSet.Table {
			ct = append(ct, element.City)
		}

		rec := logger.LogRecord{
			Country: c,
			Cities:  ct}

		jsonData, _ := json.Marshal(rec)
		return string(jsonData)
	}
	return ""
}
