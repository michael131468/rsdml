package main

import "fmt"
import "os"
import "github.com/michael131468/rsdml"

func main() {
	program_result := 0

	if len(os.Args[1:]) < 1 {
		fmt.Printf("Usage: %s [directories...]\n", os.Args[0])
		program_result = 1
		os.Exit(program_result)
	}

	for _, dir := range os.Args[1:] {
		err := rsdml.RecurseDirectory(dir)
		if err != nil {
			fmt.Println(err)
			program_result = 1
		}
	}

	os.Exit(program_result)
}

