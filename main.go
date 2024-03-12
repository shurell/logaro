package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

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
	Id       int
	Tname    string
	Hexcolor string
}

type logOwners struct {
	Id       int
	Tname    string
	Apitoken string
}

type ViewData struct {
	Title       string
	Process     []logOwners
	MessTypes   []messageTypes
	DataStrings []messReadable
}

func main() {
	var err error
	DB, err = sql.Open("sqlite3", "logaro.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", showMain)
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	mux.HandleFunc("GET /styles.css", cssHandler)
	mux.HandleFunc("POST /addlog", storeMessage)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	fmt.Println("LOGARO Web interface available at http://" + server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(terr+":", err)
	}

}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpls/favicon.ico")
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpls/styles.css")
}

func showMain(w http.ResponseWriter, r *http.Request) {
	prc := syncLogOwners()
	mtp := syncMessageTypes()
	filterLO, _ := strconv.Atoi(r.FormValue("process-filter"))
	filterMT, _ := strconv.Atoi(r.FormValue("type-filter"))
	limit := 100
	limitPos := 0
	fmt.Println(filterLO, filterMT, limit, limitPos)
	str := syncMessages(filterLO, filterMT, limit, limitPos)

	data := ViewData{
		Title:       "LOGARO",
		Process:     prc,
		MessTypes:   mtp,
		DataStrings: str,
	}
	tmpl, _ := template.ParseFiles("tmpls/index.html")
	tmpl.Execute(w, data)
}

func storeMessage(w http.ResponseWriter, r *http.Request) {
	var m mess
	m.utime = time.Now().Unix()
	_, err := DB.Exec("insert into log_messages(log_owner, message_type, utime, source_utime, message) values(?, ?, ?, ?, ?)", m.logOwner, m.messageType, m.utime, m.source_utime, m.message)
	if err != nil {
		log.Println(err)
	}
}

func syncMessages(filterLO, filterMT, limit, limitPos int) []messReadable {
	var rv []messReadable
	var limits string
	const queryBase = "select log_messages.id, log_owners.tname, message_types.tname, message_types.hexcolor, utime, source_utime, message from log_messages left join log_owners on log_owners.id=log_messages.log_owner LEFT join message_types on message_types.id=log_messages.message_type"
	if filterLO > 0 && filterMT == 0 {
		limits = " where log_messages.log_owner = ? and log_messages.message_type != ? LIMIT ?, ?"
	}

	if filterLO > 0 && filterMT > 0 {
		limits = " where log_messages.log_owner = ? and log_messages.message_type = ? LIMIT ?, ?"
	}

	if filterLO == 0 && filterMT > 0 {
		limits = " where log_messages.log_owner != ? and log_messages.message_type = ? LIMIT ?, ?"
	}
	fmt.Println(queryBase + limits)
	rows, err := DB.Query(queryBase+limits, filterLO, filterMT, limitPos, limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var ms messReadable
		var Utime, Source_utime int64
		err = rows.Scan(&ms.Id, &ms.LogOwner, &ms.MessageType, &ms.MessageColor, &Utime, &Source_utime, &ms.Message)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ms)
		ms.Utime = time.Unix(Utime, 0).Format("02.01.2006 15:05")
		ms.Source_utime = time.Unix(Source_utime, 0).Format("02.01.2006 15:05")
		rv = append(rv, ms)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return rv
}

func syncMessageTypes() []messageTypes {
	var rv []messageTypes
	rows, err := DB.Query("select id, tname, hexcolor from message_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var mt messageTypes
		err = rows.Scan(&mt.Id, &mt.Tname, &mt.Hexcolor)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(mt.Id, mt.Tname, mt.Hexcolor)
		rv = append(rv, mt)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return rv
}

func syncLogOwners() []logOwners {
	var rv []logOwners
	rows, err := DB.Query("select id, tname, apitoken from log_owners")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var lo logOwners
		err = rows.Scan(&lo.Id, &lo.Tname, &lo.Apitoken)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(lo.Id, lo.Tname, lo.Apitoken)
		rv = append(rv, lo)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return rv
}
