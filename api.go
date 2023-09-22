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
		query := r.URL.Query()
		fromCode, toCode := query.Get("from"), query.Get("to")
		dateStr, trainNo := query.Get("date"), query.Get("no")
		log.Printf("from: %s, to: %s, date: %s, no: %s", fromCode, toCode, dateStr, trainNo)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
