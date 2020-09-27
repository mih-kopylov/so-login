package main

import (
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"io"
	"log"
	"os"
	"strings"
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
		log.Println("Successfully logged in")
		os.Exit(0)
	}
	writePageContent(br)
	log.Fatal("Could not find a marker string in page content")
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
			log.Fatal("Unknown argument ", arg)
		}
		parts := strings.Split(arg, separator)
		name := parts[0]
		value := parts[1]

		if name == emailArgument {
			result.email = value
		} else if name == passwordArgument {
			result.password = value
		} else {
			log.Fatal("Unknown argument ", arg)
		}
	}
	if result.email == "" {
		log.Fatal(emailArgument, " argument is required")
	}
	if result.password == "" {
		log.Fatal(passwordArgument, " argument is required")
	}
	return result
}

func writePageContent(browser *browser.Browser) {
	file, err := os.Create("output.html")
	if err != nil {
		log.Fatal(err)
	}
	defer closeFile(file)
	_, err = io.WriteString(file, browser.Body())
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
}

type credentials struct {
	email    string
	password string
}
