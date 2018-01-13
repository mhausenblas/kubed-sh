package main

import (
	"fmt"
	"strings"
)

func husage(line string) {
	fmt.Println("Available commands:")
	fmt.Printf("%s", completer.Tree("    "))
}

func hecho(line string) {
	l := strings.Split(line, " ")
	fmt.Println(l[1])
}
