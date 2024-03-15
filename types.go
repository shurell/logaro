package main

import "database/sql"

var DB *sql.DB

var LOhash map[string]logOwners
var MThash map[string]messageTypes

type messReadable struct {
	Id           int
	LogOwner     string
	MessageColor string
	MessageType  string
	Utime        string
	Source_utime string
	Message      string
}

type messIn struct {
	LogOwnerToken string `json:"log_owner_token"`
	MessageType   int    `json:"message_type"`
	SourceUtime   int64  `json:"source_utime"`
	Message       string `json:"message"`
}

type messageTypes struct {
	Id       int    `json:"id"`
	Tname    string `json:"tname"`
	Hexcolor string `json:"hexcolor"`
}

type logOwners struct {
	Id       int    `json:"id"`
	Tname    string `json:"tname"`
	Apitoken string `json:"apitoken"`
}

type ViewData struct {
	Title       string
	SearchValue string
	Process     []logOwners
	MessTypes   []messageTypes
	DataStrings []messReadable
}
