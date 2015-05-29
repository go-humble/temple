// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

// NOTE: this file is meant to be run together with a generated file created
// by Build. By itself it won't compile because Templates is not defined.

package main

import (
	"log"
	"os"
)

type Todo struct {
	Title string
}

var todos = []Todo{
	{Title: "One"},
	{Title: "Two"},
	{Title: "Three"},
}

func main() {
	todosTmpl, err := GetTemplate("todos/index")
	if err != nil {
		log.Fatal(err)
	}
	if err := todosTmpl.Execute(os.Stdout, todos); err != nil {
		log.Fatal(err)
	}
}
