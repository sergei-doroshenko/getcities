package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Url         string `json:"url"`
	Queries     int    `json:"queriew"`
	LogPath     string `json:"logpath"`
	LogFileNane string
}

func (s Config) getFileName() string {
	current_time := time.Now().Local()
	return filepath.Join(s.LogPath, "cities-"+current_time.Format("2006-01-02")+".log")
}

func (s Config) createLogsFolder() {
	if err := os.MkdirAll(s.LogPath, 0777); err != nil {
		panic(err)
	}
}

func ReadConfig() Config {
	dat, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	config := Config{}

	if err := json.Unmarshal(dat, &config); err != nil {
		panic(err)
	}
	config.createLogsFolder()
	config.LogFileNane = config.getFileName()

	fmt.Println("Config url: ", config.Url)
	fmt.Println("Config log path: ", config.LogPath)
	fmt.Println("Config log file name: ", config.LogFileNane)
	return config
}
