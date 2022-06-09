package main

import (
	"log"
	"main/components/globals"
	"main/components/links"
	"main/components/users"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	check(globals.Initialize(globals.Context{
		DBName: "eventschedule.db",
	}))
	defer check(globals.Destroy())
	check(users.Initialize())
	check(links.Initialize())
}
