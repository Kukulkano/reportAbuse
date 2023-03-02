package main

// Use pogreb filesystem key-value database from
// https://github.com/akrylysov/pogreb

import (
	"fmt"
	"os"

	"github.com/akrylysov/pogreb"
)

var ipDB *pogreb.DB

func initDatabase() {
	var err error
	ipDB, err = pogreb.Open("reportAbuse.db", nil)
	if err != nil {
		fmt.Println("ERROR: Unable to init pogreb database 'reportAbuse.db'.")
		os.Exit(1)
	}
}

func closeDatabase() {
	ipDB.Close()
}
