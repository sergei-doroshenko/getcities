package converter

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"net/url"
	"sync"

	"github.com/sergei-doroshenko/getcities/datasource"
)

type Entry struct {
	XMLName xml.Name
	Country string
	City    string
}

type DataSet struct {
	XMLName xml.Name
	Table   []Entry
}

type XmlRecord struct {
	XMLName    xml.Name
	NewDataSet DataSet
}

func (s XmlRecord) Equals(r XmlRecord) bool {
	if &s == &r {
		return true
	}

	if s.XMLName != r.XMLName || s.NewDataSet.XMLName != r.NewDataSet.XMLName {
		return false
	}

	if len(s.NewDataSet.Table) != len(r.NewDataSet.Table) {
		return false
	}

	for i, v := range s.NewDataSet.Table {
		if r.NewDataSet.Table[i] != v {
			return false
		}
	}

	return true
}

type Converter struct {
	DataSource datasource.DataGetter
}

func (s Converter) GetXml(country string) (string, error) {

	data, err := s.DataSource.GetData(url.QueryEscape(country))

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return html.UnescapeString(string(data)), nil
}

func (s Converter) GetXmlRecord(country string) (XmlRecord, error) {
	v := XmlRecord{}
	xmlData, err := s.GetXml(country)
	if err != nil {
		return v, err
	}

	if err := xml.Unmarshal([]byte(xmlData), &v); err != nil {
		fmt.Println(err.Error())
		return v, err
	}

	return v, nil
}

func (s Converter) ConvertToJson(xmlRecord XmlRecord) (string, error) {
	records := xmlRecord.NewDataSet.Table

	if len(records) > 0 {
		var c string = xmlRecord.NewDataSet.Table[0].Country
		var ct []string = make([]string, 0, len(records))

		for _, element := range xmlRecord.NewDataSet.Table {
			ct = append(ct, element.City)
		}

		type LogRecord struct {
			Country string   `json:"country"`
			Cities  []string `json:"cities"`
			Error   string   `json:"error"`
		}

		rec := LogRecord{
			Country: c,
			Cities:  ct}

		jsonData, err := json.Marshal(rec)
		if err != nil {
			return "", err
		}

		return string(jsonData), nil
	}

	return "", errors.New("Empty XmlRecord.")
}

func (s Converter) GetLogRecord(country string) (string, error) {

	xmlRecord, err := s.GetXmlRecord(country)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	logRecord, err := s.ConvertToJson(xmlRecord)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return logRecord, nil
}

func (s Converter) GetCities(country string, msgs chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	recd, err := s.GetLogRecord(country)
	if err != nil {
		msgs <- err.Error()
	} else {
		msgs <- recd
	}
}
