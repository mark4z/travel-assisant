package main

type Info struct {
	Data struct {
		Result []string `json:"result"`
	} `json:"data"`
}

type IntervalInfo struct {
	Data struct {
		Data []Station `json:"data"`
	} `json:"data"`
}

type Station struct {
	StationName string `json:"station_name"`
	ArriveTime  string `json:"arrive_time"`
	StartTime   string `json:"start_time"`
}
