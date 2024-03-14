package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	fmt.Println("LOGARO Web interface available at http://" + server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(terr+":", err)
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

func registerLO(w http.ResponseWriter, r *http.Request) {
	var result logOwners
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(terr+":", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	if storeLO(result.Tname, result.Apitoken) {
		fmt.Println("Добавлен источник: ", result.Tname)
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte("error"))
	}
}

func removeLO(w http.ResponseWriter, r *http.Request) {
	var result logOwners
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(terr+":", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	if result.Id != 0 && result.Apitoken != "" {
		fmt.Println("Удаляем программу-источник из справочника:", result.Id, result.Apitoken)
		res, err := DB.Exec("delete from log_owners where id=? and apitoken=?", result.Id, result.Apitoken)
		if err != nil {
			fmt.Println(res)
			w.Write([]byte("ok"))
		}
	} else {
		w.Write([]byte("error"))
	}

}

func registerMT(w http.ResponseWriter, r *http.Request) {
	var result messageTypes
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(terr+":", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	if storeMT(result.Tname, result.Hexcolor) {
		fmt.Println("Добавлен тип сообщения: ", result.Tname)
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte("error"))
	}
}

func removeMT(w http.ResponseWriter, r *http.Request) {
	var result messageTypes
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(terr+":", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	if result.Id != 0 && result.Tname != "" {
		fmt.Println("Удаляем тип сообщения из справочника:", result.Id, result.Tname)
		res, err := DB.Exec("delete from message_types where id=? and tname=?", result.Id, result.Tname)
		if err != nil {
			fmt.Println(res)
			w.Write([]byte("ok"))
		}
	} else {
		w.Write([]byte("error"))
	}

}

func showMain(w http.ResponseWriter, r *http.Request) {
	prc := syncLogOwners()
	mtp := syncMessageTypes()
	filterLO, _ := strconv.Atoi(r.FormValue("process-filter"))
	filterMT, _ := strconv.Atoi(r.FormValue("type-filter"))
	limit := 500
	limitPos := 0
	if r.FormValue("type-filter") != "on" {
		limit = 10000
		limitPos = 0
	}
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

func showMainSearches(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("sval")
	data := ViewData{
		Title:       "LOGARO",
		SearchValue: "Результаты поиска по тексту: '" + search + "'",
		Process:     syncLogOwners(),
		MessTypes:   syncMessageTypes(),
		DataStrings: syncMessagesSearch(search),
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

func storeLO(a, b string) bool {
	result := false
	var count int
	err := DB.QueryRow("select count(id) from log_owners where apitoken=", b).Scan(&count)
	if err != nil {
		log.Println(err)
	}
	if count == 0 && a != "" {
		_, err = DB.Exec("insert into log_owners(tname, apitoken) values(?, ?)", a, b)
		if err != nil {
			log.Println(err)
		} else {
			result = true
		}
	}
	return result
}

func storeMT(a, b string) bool {
	result := false
	var count int
	err := DB.QueryRow("select count(id) from message_types where tname=", a).Scan(&count)
	if err != nil {
		log.Println(err)
	}
	if count == 0 && a != "" {
		_, err = DB.Exec("insert into message_types(tname, hexcolor) values(?, ?)", a, b)
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

func syncMessagesSearch(sval string) []messReadable {
	var rv []messReadable

	queryBase := "select log_messages.id, log_owners.tname, message_types.tname, message_types.hexcolor, utime, source_utime, message from log_messages left join log_owners on log_owners.id=log_messages.log_owner LEFT join message_types on message_types.id=log_messages.message_type where log_messages.message LIKE '%" + sval + "%'"
	fmt.Println(queryBase)
	rows, err := DB.Query(queryBase)
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
