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

type Select struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type TrainReq struct {
	TrainNo     string `json:"train_no"`
	FromStation string `json:"from-station"`
	ToStation   string `json:"to-station"`
	Date        string `json:"__date"`
}

type TrainRes struct {
	TrainNo   string `json:"train_no"`
	TrainCode string `json:"train_code"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`

	StartStation     string `json:"start_station"`
	StartStationName string `json:"start_station_name"`
	EndStation       string `json:"end_station"`
	EndStationName   string `json:"end_station_name"`

	FromStation     string `json:"from_station"`
	FromStationName string `json:"from_station_name"`
	ToStation       string `json:"to_station"`
	ToStationName   string `json:"to_station_name"`

	TwoSeat     string `json:"two_seat"`
	OneSeat     string `json:"one_seat"`
	SpecialSeat string `json:"special_seat"`
}
