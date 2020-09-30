package main

import (
	"fmt"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	authUsingSurf()
}

func authUsingSurf() {
	creds := readCredentials()
	br := surf.NewBrowser()
	br.SetUserAgent(agent.Chrome())
	fatal(br, br.Open("https://stackoverflow.com"))
	fatal(br, br.Click("a.login-link"))
	form, err := br.Form("#login-form")
	fatal(br, err)
	fatal(br, form.Input("email", creds.email))
	fatal(br, form.Input("password", creds.password))
	fatal(br, form.Submit())

	fatal(br, br.Click("a.my-profile"))

	pageContent := br.Body()
	if strings.Contains(pageContent, "Edit profile and settings") {
		writeLog("Successfully logged in")
		os.Exit(0)
	}
	writePageContent(br)
	writeErrorLog("Could not find a marker string in page content")
}

func readCredentials() credentials {
	result := credentials{}

	const emailArgument = "--email"
	const passwordArgument = "--password"
	for i, arg := range os.Args {
		if i == 0 {
			//skip the initial argument
			continue
		}
		const separator = "="
		if !strings.Contains(arg, separator) {
			writeErrorLog("Unknown argument ", arg)
		}
		parts := strings.Split(arg, separator)
		name := parts[0]
		value := parts[1]

		if name == emailArgument {
			result.email = value
		} else if name == passwordArgument {
			result.password = value
		} else {
			writeErrorLog("Unknown argument ", arg)
		}
	}
	if result.email == "" {
		writeErrorLog(emailArgument, " argument is required")
	}
	if result.password == "" {
		writeErrorLog(passwordArgument, " argument is required")
	}
	return result
}

func writePageContent(browser *browser.Browser) {
	file, err := os.Create("output.html")
	if err != nil {
		writeErrorLog(err)
	}
	defer closeFile(file)
	_, err = io.WriteString(file, browser.Body())
	if err != nil {
		writeErrorLog(err)
	}
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func fatal(browser *browser.Browser, err error) {
	if err != nil {
		writePageContent(browser)
		writeErrorLog(err)
	}
}

func writeLog(values ...interface{}) {
	log.Println(values)
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFile(file)
	stringToAppend := fmt.Sprintln(time.Now().Format(time.RFC3339), values)
	_, err = file.WriteString(stringToAppend)
	if err != nil {
		log.Fatal(err)
	}
}

func writeErrorLog(values ...interface{}) {
	writeLog(values)
	os.Exit(1)
}

type credentials struct {
	email    string
	password string
}
