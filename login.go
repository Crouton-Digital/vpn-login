// Package: main
// Language: go
// Path: login.go
// Author: Rolands Kalpins
// Date: 2022-08-09

package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var arguments string
	for _, arg := range os.Args[:] {
		arguments += arg + " "
	}

	arguments += " *** With ENV ***"
	arguments += os.Getenv("username")
	arguments += " " + os.Getenv("password")

	arguments += "\n"

	if _, err := f.Write([]byte(arguments)); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}

	os.WriteFile(os.Getenv("auth_control_file"), []byte("1"), 0644)
	os.Exit(0)
}
