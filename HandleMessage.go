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

type MessageCarrier struct {
	msg []Message
}

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		/*
			Todo: The get method currently has little to no security.
			implement some sort of credential system to ensure the
			right person is the only one who can fetch the message
		 */
		origin := r.Header.Get("X-REAL-IP") // Only returns messages directed to origin
		// Susceptible to ip spoofing

		msgInt, exists := Cache.Get(origin)
		carrier := MessageCarrier{}
		mapstructure.Decode(msgInt, &carrier)

		if exists {

			w.Header().Set("Content-Type","application/json")
			resp, _ := json.Marshal(carrier)

			w.Write(resp)

			Cache.Delete(origin)

		} else {
			// Todo: return a no messages for you message
		}

		break
	case "POST":
		msg := Message{}
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&msg)

		carrier := MessageCarrier{}
		msg.Origin = r.Header.Get("X-REAL-IP")

		msgCarrier, exists := Cache.Get(msg.Destination)
		if exists {
			// Add msg to msgCarrier and replace cache
			mapstructure.Decode(msgCarrier, &carrier)
			carrier.msg = append(carrier.msg, msg)
			Cache.Replace(msg.Destination, carrier, 5*time.Minute) // Replace previous data, with new data
		} else {
			// Create new msgCarrier and add to cache
			Cache.Add(msg.Destination, carrier, 5 * time.Minute)
		}

		log.Print(msgCarrier)

		carrier = MessageCarrier{}


		log.Print(carrier.msg[0])

		break
	default:
		log.Print("No valid method used")
		break
	}
}