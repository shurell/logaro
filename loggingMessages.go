package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

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
