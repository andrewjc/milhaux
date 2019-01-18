package main

import (
	"fmt"
	"log"
	"net/smtp"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var messageChannel = make(chan MainEventMessage)

	for {
		for i := 0; i <= 1000; i++ {
			go spawner()
		}
		time.Sleep(time.Second * 3)
	}

	beginMainLoop(messageChannel)
}

type MainEventMessage struct {
	data        interface{}
	MessageType interface{}
}

func beginMainLoop(messageChannel chan MainEventMessage) {

	for {
		eventMessage := <-messageChannel
		if eventMessage.MessageType == 1 {

		}
	}
}

func spawner() {
	performUnauthenticatedTest()
}

func performUnauthenticatedTest() {
	// send an email from external to internal
	c, err := smtp.Dial("0.0.0.0:25")
	if err != nil {
		log.Println(err)
	}

	if err := c.Mail("andrew.cranston@google.com"); err != nil {
		log.Println(err)
	}

	if err := c.Rcpt("andrewc@cranston.com"); err != nil {
		log.Println(err)
	}

	wc, err := c.Data()
	if err != nil {
		log.Println(err)
	}

	_, err = fmt.Fprintf(wc, "Test email body")
	if err != nil {
		log.Println(err)
	}

	err = wc.Close()
	if err != nil {
		log.Println(err)
	}

	err = c.Quit()
	if err != nil {
		log.Println(err)
	}
}
