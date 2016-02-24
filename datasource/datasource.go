package datasource

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type GetDataSource struct {
	Url string
}

type DataGetter interface {
	GetData(country string) ([]byte, error)
}

func (s GetDataSource) GetData(country string) ([]byte, error) {
	fmt.Println("Start getting cities from ", country)
	resp, err := http.Get(s.Url + country)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return content, nil
}
