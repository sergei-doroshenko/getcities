package logger

import (
	"fmt"
	"os"
	"sync"
)

type FileLogger struct {
	LogFileName string
}

func (s FileLogger) appendToFile(text string) {

	f, err := os.OpenFile(s.LogFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
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

func (s FileLogger) Log(msgCount int, msgs chan string, wg *sync.WaitGroup) {

	defer wg.Done()

	buffer := ""

	for msgCount > 0 {
		text := <-msgs
		buffer += text + "\n\n"
		fmt.Println("appender get messege: ", text)

		msgCount--
	}

	s.appendToFile(buffer)
}
