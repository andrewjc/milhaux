package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {
	performUnauthenticatedTest()
}
func performUnauthenticatedTest() {
	// send an email from external to internal
	c, err := smtp.Dial("0.0.0.0:25")
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Mail("andrew.cranston@google.com"); err != nil {
		log.Fatal(err)
	}

	if err := c.Rcpt("andrewc@cranston.com"); err != nil {
		log.Fatal(err)
	}

	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintf(wc, "Test email body")
	if err != nil {
		log.Fatal(err)
	}

	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}
