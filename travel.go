package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

var (
	indexUrl  = "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc&fs=%s,%s&ts=%s,%s&date=%s&flag=N,N,Y"
	trainUrl  = "https://kyfw.12306.cn/otn/leftTicket/queryZ?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=ADULT"
	mapperUrl = "https://kyfw.12306.cn/otn/resources/js/framework/station_name.js?station_version=1.9274"
	passUrl   = "https://kyfw.12306.cn/otn/czxx/queryByTrainNo?train_no=%s&from_station_telecode=%s&to_station_telecode=%s&depart_date=%s"
)

var (
	client http.Client
	m      map[string]string
	__init bool
)

func init() {
	jar, _ := cookiejar.New(nil)
	client = http.Client{
		Jar: jar,
	}
}

func index(from, to, date string) {
	if __init {
		return
	}
	m = mapper()
	res, err := client.Get(fmt.Sprintf(indexUrl, m[from], from, m[to], to, date))
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		panic(fmt.Sprintf("index error %d", res.StatusCode))
	}
	_ = res.Body.Close()

	__init = true
}

func findAllTrain(from, to, date string) []*TrainRes {
	index(from, to, date)

	resp, err := client.Get(fmt.Sprintf(trainUrl, date, from, to))
	if err != nil {
		panic(err)
	}
	jsonStr, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	info := &struct {
		Data struct {
			Result []string `json:"result"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(jsonStr, info)
	if err != nil {
		panic(err)
	}
	res := make([]*TrainRes, 0)
	for _, t := range info.Data.Result {
		//filter G+D
		trainRes := decode(t)
		if strings.HasPrefix(trainRes.TrainNo, "G") || strings.HasPrefix(trainRes.TrainNo, "D") {
			res = append(res, trainRes)
		}
	}
	if len(res) == 0 {
		panic(fmt.Sprintf("can not found target %s-%s %s", m[from], m[to], date))
	}
	return res
}

func findTrainByNo(wantNo, from, to, date string) *TrainRes {
	trains := findAllTrain(from, to, date)

	for _, t := range trains {
		if t.TrainNo == wantNo {
			return t
		}
	}
	panic(fmt.Sprintf("can not found target findTrainByNo %s %s-%s %s", wantNo, m[from], m[to], date))
}

func findPassStationsByCode(id, from, to, date string) []Station {
	res, err := client.Get(fmt.Sprintf(passUrl, id, from, to, date))
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	interval := &struct {
		Data struct {
			Data []Station `json:"data"`
		} `json:"data"`
	}{}
	_ = json.Unmarshal(body, interval)
	return interval.Data.Data
}

func findTrainByCode(wantCode, from, to, date string) *TrainRes {
	trains := findAllTrain(from, to, date)

	for _, t := range trains {
		if t.TrainCode == wantCode {
			return t
		}
	}
	panic(fmt.Sprintf("can not found target findTrainByCode %s %s-%s %s", wantCode, m[from], m[to], date))
}

// save mapper to file with yyyy-mm-dd name
func mapper() map[string]string {
	fileName := "./mapper-" + time.Now().Format("2006-01-02") + ".json"
	//clean all tmp mapper file not today
	if files, err := os.ReadDir("./"); err == nil {
		for _, f := range files {
			if strings.Contains(f.Name(), "mapper") && !strings.Contains(f.Name(), time.Now().Format("2006-01-02")) {
				_ = os.Remove(f.Name())
			}
		}
	}
	temp := make(map[string]string, 128)

	if mFile, err := os.Open(fileName); err == nil {
		bytes, err := io.ReadAll(mFile)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(bytes, &temp); err != nil {
			panic(err)
		}
		return temp
	}

	response, err := client.Get(mapperUrl)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(response.Body)
	_ = response.Body.Close()

	resp := string(body)
	city := strings.Split(resp, "@")
	city = city[1:]
	for c := range city {
		fields := strings.Split(city[c], "|")
		temp[fields[1]] = fields[2]
		temp[fields[2]] = fields[1]
	}
	fd, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	jsonStr, err := json.Marshal(temp)
	if err != nil {
		panic(err)
	}
	_, _ = fd.Write(jsonStr)
	if err := fd.Close(); err != nil {
		panic(err)
	}
	return temp
}

func decode(info string) *TrainRes {
	fields := strings.Split(info, "|")
	return &TrainRes{
		TrainCode:        fields[2],
		TrainNo:          fields[3],
		SpecialSeat:      fields[32],
		OneSeat:          fields[31],
		TwoSeat:          fields[30],
		StartStation:     fields[4],
		StartStationName: m[fields[4]],
		EndStation:       fields[5],
		EndStationName:   m[fields[5]],
		FromStation:      fields[6],
		FromStationName:  m[fields[6]],
		ToStation:        fields[7],
		ToStationName:    m[fields[7]],
		StartTime:        fields[8],
		EndTime:          fields[9],
	}
}
