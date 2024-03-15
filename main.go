package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error
	LOhash = make(map[string]logOwners)
	MThash = make(map[string]messageTypes)
	DB, err = sql.Open("sqlite3", "logaro.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", showMain)
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	mux.HandleFunc("GET /j.js", jsHandler)
	mux.HandleFunc("GET /styles.css", cssHandler)
	mux.HandleFunc("POST /addlog", storeMessage)
	mux.HandleFunc("POST /registerprog", registerLO)
	mux.HandleFunc("POST /removeprog", removeLO)
	mux.HandleFunc("POST /registertype", registerMT)
	mux.HandleFunc("POST /removetype", removeMT)
	mux.HandleFunc("POST /search", showMainSearches)
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	log.Println(tnotice00 + server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Println(terr+":", err)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpls/favicon.ico")
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpls/j.js")
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpls/styles.css")
}

func showMain(w http.ResponseWriter, r *http.Request) {
	filterLO, _ := strconv.Atoi(r.FormValue("process-filter"))
	filterMT, _ := strconv.Atoi(r.FormValue("type-filter"))
	limit := 500
	limitPos := 0
	if r.FormValue("type-filter") != "on" {
		limit = 10000
		limitPos = 0
	}
	log.Println(filterLO, filterMT, limit, limitPos)
	data := ViewData{
		Title:       "LOGARO",
		Process:     syncLogOwners(),
		MessTypes:   syncMessageTypes(),
		DataStrings: syncMessages(filterLO, filterMT, limit, limitPos),
	}
	tmpl, _ := template.ParseFiles("tmpls/index.html")
	tmpl.Execute(w, data)
}

func showMainSearches(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("sval")
	data := ViewData{
		Title:       "LOGARO",
		SearchValue: ttext01 + search + "'",
		Process:     syncLogOwners(),
		MessTypes:   syncMessageTypes(),
		DataStrings: syncMessagesSearch(search),
	}
	tmpl, _ := template.ParseFiles("tmpls/index.html")
	tmpl.Execute(w, data)
}

func storeMessage(w http.ResponseWriter, r *http.Request) {
	var result messIn
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(terr+":", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(terror02)
	}
	//Простая проверка
	if result.LogOwnerToken == "" || result.Message == "" || result.MessageType == 0 || result.SourceUtime == 0 {
		w.Write([]byte(terror01))
		return
	}
	if storeMS(result.LogOwnerToken, result.MessageType, result.SourceUtime, result.Message) {
		log.Println(tnotice01, result.Message)
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte(terr))
	}
}

func storeMS(LogOwnerToken string, MessageType int, SourceUtime int64, Message string) bool {
	result := false
	loId := LOhash[LogOwnerToken].Id
	if loId > 0 && MessageType != 0 && Message != "" {
		_, err := DB.Exec("insert into log_messages(log_owner, message_type, utime, source_utime, message) values(?, ?, ?, ?, ?)", loId, MessageType, time.Now().Unix(), SourceUtime, Message)
		if err != nil {
			log.Println(err)
		} else {
			result = true
		}
	}
	return result
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
	rows, err := DB.Query(queryBase+limits, filterLO, filterMT, limitPos, limit)
	if err != nil {
		log.Fatal(terr+" [sm01]", err)
	}
	defer rows.Close()
	for rows.Next() {
		var ms messReadable
		var LogOwner, MessageType, MessageColor sql.NullString
		var Utime, Source_utime int64
		err = rows.Scan(&ms.Id, &LogOwner, &MessageType, &MessageColor, &Utime, &Source_utime, &ms.Message)
		if err != nil {
			log.Fatal(terr+" [sm02]", err)
		}
		ms.Utime = time.Unix(Utime, 0).Format("02.01.2006 15:05")
		ms.Source_utime = time.Unix(Source_utime, 0).Format("02.01.2006 15:05")
		if LogOwner.Valid {
			ms.LogOwner = LogOwner.String
		} else {
			ms.LogOwner = "--*"
		}
		ms.MessageColor = MessageColor.String
		if MessageType.Valid {
			ms.MessageType = MessageType.String
		} else {
			ms.MessageType = "--*"
		}
		rv = append(rv, ms)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(terr+" [sm03]", err)
	}
	return rv
}

func syncMessagesSearch(sval string) []messReadable {
	var rv []messReadable
	var LogOwner, MessageType, MessageColor sql.NullString
	queryBase := "select log_messages.id, log_owners.tname, message_types.tname, message_types.hexcolor, utime, source_utime, message from log_messages left join log_owners on log_owners.id=log_messages.log_owner LEFT join message_types on message_types.id=log_messages.message_type where log_messages.message LIKE '%" + sval + "%'"
	log.Println(queryBase)
	rows, err := DB.Query(queryBase)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var ms messReadable
		var Utime, Source_utime int64
		err = rows.Scan(&ms.Id, &LogOwner, &MessageType, &MessageColor, &Utime, &Source_utime, &ms.Message)
		if err != nil {
			log.Fatal(terr+" [sm02]", err)
		}
		ms.Utime = time.Unix(Utime, 0).Format("02.01.2006 15:05")
		ms.Source_utime = time.Unix(Source_utime, 0).Format("02.01.2006 15:05")
		if LogOwner.Valid {
			ms.LogOwner = LogOwner.String
		} else {
			ms.LogOwner = "--*"
		}
		ms.MessageColor = MessageColor.String
		if MessageType.Valid {
			ms.MessageType = MessageType.String
		} else {
			ms.MessageType = "--*"
		}
		rv = append(rv, ms)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return rv
}
