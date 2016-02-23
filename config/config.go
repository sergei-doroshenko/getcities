package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Url     string `json:"url"`
	Queries int    `json:"queriew"`
}

func ReadConfig() Config {
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
