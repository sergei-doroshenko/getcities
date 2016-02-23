package test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/sergei-doroshenko/getcities/converter"
)

func TestConvertToJson(t *testing.T) {

	var table []converter.Entry
	table = append(table, converter.Entry{xml.Name{}, "Belarus", "Riga Airport"})
	table = append(table, converter.Entry{xml.Name{}, "Belarus", "Vitebsk"})
	table = append(table, converter.Entry{xml.Name{}, "Belarus", "Minsk"})

	var dataSet = converter.DataSet{}
	dataSet.XMLName = xml.Name{}
	dataSet.Table = table

	var xmlRecord = converter.XmlRecord{}
	xmlRecord.XMLName = xml.Name{}
	xmlRecord.NewDataSet = dataSet

	fmt.Println(xmlRecord)

	var json1 = `{"country":"Belarus","cities":["Riga Airport","Vitebsk","Minsk"],"error":""}`

	json := converter.ConvertToJson(xmlRecord)
	fmt.Println(json)
	if json != json1 {
		t.Error("Expected "+json1+", got ", json)
	}
}
