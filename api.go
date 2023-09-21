package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func serve() {
	m = mapper()

	http.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
		var selects []Select
		for k, v := range m {
			if k[0] > 'Z' || rune(k[0]) < 'A' {
				continue
			}
			selects = append(selects, Select{k, v})
		}
		res, _ := json.Marshal(selects)
		w.Write(res)
	})

	http.HandleFunc("/walk", func(w http.ResponseWriter, r *http.Request) {
		fromCode := r.URL.Query().Get("from")
		toCode := r.URL.Query().Get("to")
		dateStr := r.URL.Query().Get("date")
		trainNo := r.URL.Query().Get("no")

		log.Printf("start walk no: %s, date: %s, from: %s, to: %s", trainNo, dateStr, m[fromCode], m[toCode])
		index(fromCode, m[fromCode], toCode, m[toCode], dateStr)

		id := train(trainNo, fromCode, toCode, dateStr)
		stations := pass(id, fromCode, toCode, dateStr)

		var trains []*TrainRes

		for i := 0; i < len(stations); i++ {
			trains = append(trains, trace(id, dateStr, stations[0], stations[i]))
		}
		for i := 0; i < len(stations); i++ {
			trains = append(trains, trace(id, dateStr, stations[i], stations[len(stations)-1]))
		}
		res, _ := json.Marshal(trains)
		w.Header().Add("Content-Type", "application/json")
		w.Write(res)
	})

	http.ListenAndServe(":8080", nil)
}
