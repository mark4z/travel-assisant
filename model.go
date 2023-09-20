package main

type Info struct {
	Data struct {
		Result []string `json:"result"`
	} `json:"data"`
}

type IntervalInfo struct {
	Data struct {
		Data []struct {
			StationName string `json:"station_name"`
		} `json:"data"`
	} `json:"data"`
}
