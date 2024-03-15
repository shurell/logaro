package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func registerLO(w http.ResponseWriter, r *http.Request) {
	var result logOwners
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(terr+":", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(terror02)
	}
	if storeLO(result.Tname, result.Apitoken) {
		log.Println(tnotice02, result.Tname)
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte(terr))
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
			var tlo logOwners
			err = DB.QueryRow("select id, tname, apitoken from log_owners where tname=? and apitoken=?", a, b).Scan(&tlo.Id, &tlo.Tname, &tlo.Apitoken)
			if err != nil {
				log.Println(terr, err)
			} else {
				LOhash[b] = tlo
			}
		}
	}
	return result
}

func removeLO(w http.ResponseWriter, r *http.Request) {
	var result logOwners
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(terr+":", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(terror02)
	}
	if result.Id != 0 && result.Apitoken != "" {
		_, err := DB.Exec("delete from log_owners where id=? and apitoken=?", result.Id, result.Apitoken)
		if err != nil {
			delete(LOhash, result.Apitoken)
			w.Write([]byte("ok"))
		}
	} else {
		w.Write([]byte(terr))
	}
}

func syncLogOwners() []logOwners {
	var rv []logOwners
	clear(LOhash)
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
		LOhash[lo.Apitoken] = lo
		rv = append(rv, lo)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return rv
}
