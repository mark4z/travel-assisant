package main

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
	FromStation string `json:"from" form:"from" binding:"required"`
	ToStation   string `json:"to" form:"to" binding:"required"`
	Date        string `json:"date" form:"date" binding:"required"`
}

type SearchReq struct {
	TrainNo string `json:"no" form:"no"`
	TrainReq
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
