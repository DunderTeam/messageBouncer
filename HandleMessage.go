package main

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Name        string `json:"name"`
	Data        string `json:"data"`
	Timestamp   string `json:"timestamp"`
	Destination string `json:"destination"`
	Origin 		string `json:"origin"`
}

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		/*
			Todo: The get method currently has little to no security.
			implement some sort of credential system to ensure the
			right person is the only one who can fetch the message
		 */
		origin := GetIP(r) // Only returns messages directed to origin
		// Susceptible to ip spoofing

		msgInterface, exists := Cache.Get(origin)
		log.Print("Origin is : " + origin)
		var msgArray []Message

		mapstructure.Decode(msgInterface, &msgArray)

		if exists {

			w.Header().Set("Content-Type","application/json")
			resp, _ := json.Marshal(msgArray)

			log.Print(msgArray)

			w.Write(resp)

			Cache.Delete(origin)

			log.Print("You requested your messages")

		} else {
			// Todo: return a no messages for you message
		}

		break
	case "POST":
		msg := Message{}
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&msg)

		msg.Origin = GetIP(r)

		msgInterface, exists := Cache.Get(msg.Destination)
		var msgArray []Message

		if exists {
			// Add msg to msgCarrier and replace cache
			mapstructure.Decode(msgInterface, &msgArray)
			msgArray = append(msgArray, msg)
			Cache.Replace(msg.Destination, msgArray, 5*time.Minute) // Replace previous data, with new data
		} else {
			// Create new msgCarrier and add to cache
			msgArray = append(msgArray, msg)
			Cache.Add(msg.Destination, msgArray, 5 * time.Minute)
		}
		log.Print(msgArray)

		out := "200"



		w.Write(([]byte)(out))

		break
	default:
		log.Print("No valid method used")
		break
	}
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}