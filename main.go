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
	id       int
	tname    string
	hexcolor string
}

type logOwners struct {
	id       int
	tname    string
	apitoken string
}

type ViewData struct {
	Title       string
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
	var str []messReadable
	str = append(str, messReadable{Id: 0, LogOwner: "Test prog", MessageColor: "FF00AA", MessageType: "Success", Utime: "12.03.2024 11:54", Source_utime: "12.03.2024 11:53", Message: "Запущено успешно"})
	data := ViewData{
		Title:       "LOGARO",
		DataStrings: str,
	}
	tmpl, _ := template.ParseFiles("tmpls/index.html")
	tmpl.Execute(w, data)
}

func storeMessage(m mess) error {
	m.utime = time.Now().Unix()
	_, err := DB.Exec("insert into log_messages(log_owner, message_type, utime, source_utime, message) values(?, ?, ?, ?, ?)", m.logOwner, m.messageType, m.utime, m.source_utime, m.message)
	if err != nil {
		return err
	}
	return nil
}

func getMessageTypes() string {
	var rv string
	rows, err := DB.Query("select id, tname, hexcolor from message_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var mt messageTypes
		err = rows.Scan(&mt.id, &mt.tname, &mt.hexcolor)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(mt.id, mt.tname, mt.hexcolor)
		rv += "<tr><td>" + strconv.Itoa(mt.id) + "</td><td>" + mt.tname + "</td><td>" + mt.hexcolor + "</td><td><a>" + tdel + "</a></td></tr>"
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return rv
}
