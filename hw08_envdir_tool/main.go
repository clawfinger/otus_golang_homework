package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Println("Usage: go-envdir [DIR] [Programm]")
		os.Exit(1)
	}

	env, err := ReadDir(args[1])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	returnCode := RunCmd(args[2:], env)
	os.Exit(returnCode)
}
