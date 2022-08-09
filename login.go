package main

import (
	"fmt"
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

	arguments += "\n"

	if _, err := f.Write([]byte(arguments)); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}

	fmt.Fprint(os.Stdout, 0)
}
