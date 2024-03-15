package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func registerMT(w http.ResponseWriter, r *http.Request) {
	var result messageTypes
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(terr+":", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(terror02)
	}
	if storeMT(result.Tname, result.Hexcolor) {
		log.Println(tnotice04, result.Tname)
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte(terr))
	}
}

func storeMT(a, b string) bool {
	result := false
	var count int
	err := DB.QueryRow("select count(id) from message_types where tname=?", a).Scan(&count)
	if err != nil {
		log.Println(terr+" [mt01]", err)
	}
	if count == 0 && a != "" {
		_, err = DB.Exec("insert into message_types(tname, hexcolor) values(?, ?)", a, b)
		if err != nil {
			log.Println(terr+" [mt02]", err)
		} else {
			result = true
			var tmt messageTypes
			err = DB.QueryRow("select id, tname, hexcolor from message_types where tname=? and hexcolor=?", a, b).Scan(&tmt.Id, &tmt.Tname, &tmt.Hexcolor)
			if err != nil {
				log.Println(terr+" [mt03]", err)
			} else {
				MThash[a] = tmt
			}
		}
	}
	return result
}

func removeMT(w http.ResponseWriter, r *http.Request) {
	var result messageTypes
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(terr+":", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(terror02)
	}
	if result.Id != 0 && result.Tname != "" {
		log.Println(tnotice05, result.Id, result.Tname)
		_, err := DB.Exec("delete from message_types where id=? and tname=?", result.Id, result.Tname)
		if err != nil {
			delete(MThash, result.Tname)
			w.Write([]byte("ok"))
		}
	} else {
		w.Write([]byte(terr))
	}
}

func syncMessageTypes() []messageTypes {
	var rv []messageTypes
	clear(MThash)
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
		MThash[mt.Tname] = mt
		rv = append(rv, mt)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return rv
}
