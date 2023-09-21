package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
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
)

func init() {
	jar, _ := cookiejar.New(nil)
	client = http.Client{
		Jar: jar,
	}
}

func main() {
	m = mapper()

	no := "G3133"
	date := "2023-10-01"
	fromZh := "郑州"
	toZh := "长治"
	index(m[fromZh], fromZh, m[toZh], toZh, date)
	id := train(no, m[fromZh], m[toZh], date)
	stations := pass(id, m[fromZh], m[toZh], date)
	for i := len(stations) - 1; i > 0; i-- {
		train(no, m[stations[0]], m[stations[i]], date)
	}
	for i := 0; i < len(stations)-1; i++ {
		train(no, m[stations[i]], m[stations[len(stations)-1]], date)
	}
}

func index(from, fromZh, to, toZh, date string) {
	res, err := client.Get(fmt.Sprintf(indexUrl, fromZh, from, toZh, to, date))
	if err != nil {
		panic(err)
	}
	_ = res.Body.Close()
}

func train(wantNo, from, to, date string) string {
	resp, err := client.Get(fmt.Sprintf(trainUrl, from, to, date))
	if err != nil {
		panic(err)
	}
	jsonStr, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	info := &Info{}
	err = json.Unmarshal(jsonStr, info)
	if err != nil {
		panic(err)
	}
	extraInfo := func(info string) (id, no, s, o, t string) {
		infos := strings.Split(info, "|")
		return infos[2], infos[3], infos[32], infos[31], infos[30]
	}

	for _, v := range info.Data.Result {
		id, no, s, o, t := extraInfo(v)
		if no == wantNo {
			fmt.Println(fmt.Sprintf("train :%s [%s:%s] \ttwo：%s\tone：%s\tspecial：%s\tid:%s\t\n  ", no, m[from], m[to], s, o, t, id))
			return id
		}
	}
	return ""
}

func pass(id, from, to, date string) []string {
	res, err := client.Get(fmt.Sprintf(passUrl, id, from, to, date))
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	interval := &IntervalInfo{}
	_ = json.Unmarshal(body, interval)
	var p []string
	for i := range interval.Data.Data {
		p = append(p, interval.Data.Data[i].StationName)
	}
	return p
}

func mapper() map[string]string {
	response, err := client.Get(mapperUrl)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(response.Body)
	_ = response.Body.Close()

	resp := string(body)
	city := strings.Split(resp, "@")
	city = city[1:]
	m := make(map[string]string, len(city)*2)
	for c := range city {
		fields := strings.Split(city[c], "|")
		m[fields[1]] = fields[2]
		m[fields[2]] = fields[1]
	}
	return m
}
