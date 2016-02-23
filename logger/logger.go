package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogRecord struct {
	Country string   `json:"country"`
	Cities  []string `json:"cities"`
	Error   string   `json:"error"`
}

func createLogFile(filename string) {
	_, err3 := os.Create(filename)
	if err3 != nil {
		fmt.Println(err3.Error())
	} else {
		fmt.Println("Created file: ", filename)
	}
}

func createLogsFolder(path string) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		// fmt.Fprintf("MkdirAll %q : %s", path, err.Error())
		fmt.Println(err.Error())
	}
}

func appendToFile(path, filename, text string) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Folder not exists")
		createLogsFolder(path)
	}

	fullpath := filepath.Join(path, filename)

	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		fmt.Println("File not exists")
		createLogFile(fullpath)
	}

	f, err := os.OpenFile(fullpath, os.O_APPEND|os.O_WRONLY, 0600) //0644
	if err != nil {
		fmt.Println("File doesn't exist")
		fmt.Println("Open file error: ", err.Error())

	} else {
		fmt.Println(f.Name())

		if _, err = f.WriteString(text + "\n"); err != nil {
			fmt.Println("Can't write to file")
			fmt.Println(err.Error())
		}

	}

	defer f.Close()
}

func LogToFile(msgCount int, msgs chan string, wg *sync.WaitGroup) {

	defer wg.Done()

	current_time := time.Now().Local()
	path := "logs"
	filename := "cities-" + current_time.Format("2006-01-02") + ".log"

	for msgCount > 0 {
		text := <-msgs
		fmt.Println("appender get messege: ", text)
		appendToFile(path, filename, text)
		msgCount--
	}
}
