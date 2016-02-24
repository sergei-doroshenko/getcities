package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/sergei-doroshenko/getcities/config"
	"github.com/sergei-doroshenko/getcities/converter"
	"github.com/sergei-doroshenko/getcities/datasource"
	"github.com/sergei-doroshenko/getcities/logger"
)

func main() {
	conf := config.ReadConfig()
	dtsource := datasource.GetDataSource{Url: conf.Url}
	conv := converter.Converter{DataSource: dtsource}
	fileLogger := logger.FileLogger{LogFileName: conf.LogFileNane}

	countries := os.Args[1:]

	var msgs = make(chan string, len(countries))
	var wg sync.WaitGroup

	for _, country := range countries {
		wg.Add(1)
		go conv.GetCities(country, msgs, &wg)
	}

	wg.Add(1)
	go fileLogger.Log(len(countries), msgs, &wg)
	wg.Wait()

	fmt.Println("main finished")
}
