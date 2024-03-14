package main

import "database/sql"

var DB *sql.DB

var LogOwnersGlobal []logOwners
var MessTypesGlobal []messageTypes

type mess struct {
	id           int
	logOwner     int
	messageType  int
	utime        int64
	source_utime int64
	message      string
}

type messReadable struct {
	Id           int
	LogOwner     string
	MessageColor string
	MessageType  string
	Utime        string
	Source_utime string
	Message      string
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
