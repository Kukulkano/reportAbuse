package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Logfiles     []string `json:"logfiles"`
	DatabasePath string   `json:"databasePath"`
	MyPage       string   `json:"myPage"`
	RegEx        string   `json:"regExp"`
	Mode         string   `json:"mode"`
	RegGroupDate int      `json:"regGroupDate"`
	RegGroupIP   int      `json:"regGroupIP"`
	RegGroupPage int      `json:"regGroupPage"`
	MinAttacks   int      `json:"minAttacks"`
	SmtpHost     string   `json:"smtpHost"`
	SmtpUser     string   `json:"smtpUser"`
	SmtpPwd      string   `json:"smtpPwd"`
	SmtpCopy     string   `json:"smtpCopy"`
}

var cfg Configuration

func loadConfig() {
	file, err := os.Open(*configFile)
	if err != nil {
		panic("Config file not found or not given!")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	cfg = Configuration{}
	err = decoder.Decode(&cfg)
	if err != nil {
		panic("Invalid config.json")
	}
}
