package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

var client http.Client
var indexUrl = "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc&fs=#{from-zh},#{from}&ts=#{to-zh},#{to}&date=#{date}&flag=N,N,Y"
var trainUrl = "https://kyfw.12306.cn/otn/leftTicket/queryZ?leftTicketDTO.train_date=#{date}&leftTicketDTO.from_station=#{from}&leftTicketDTO.to_station=#{to}&purpose_codes=ADULT"
var mapperUrl = "https://kyfw.12306.cn/otn/resources/js/framework/station_name.js?station_version=1.9274"
var passUrl = "https://kyfw.12306.cn/otn/czxx/queryByTrainNo?train_no=#{id}&from_station_telecode=#{from}&to_station_telecode=#{to}&depart_date=#{date}"

func init() {
	jar, _ := cookiejar.New(nil)
	client = http.Client{
		Jar: jar,
	}
}

func main() {
	m := mapper()
	fromZh := "郑州"
	from := m[fromZh]
	toZh := "长治"
	to := m[toZh]
	date := "2023-10-01"
	no := "G3133"
	index(from, fromZh, to, toZh, date)
	id := train(no, from, to, date, m)
	stations := pass(id, from, to, date)
	for i := len(stations) - 1; i > 0; i-- {
		train(no, m[stations[0]], m[stations[i]], date, m)
	}
	for i := 0; i < len(stations)-1; i++ {
		train(no, m[stations[i]], m[stations[len(stations)-1]], date, m)
	}
}

func index(from, fromZh, to, toZh, date string) {
	target := indexUrl
	target = strings.ReplaceAll(target, "#{from-zh}", fromZh)
	target = strings.ReplaceAll(target, "#{from}", from)
	target = strings.ReplaceAll(target, "#{to-zh}", toZh)
	target = strings.ReplaceAll(target, "#{to}", to)
	target = strings.ReplaceAll(target, "#{date}", date)
	res, err := client.Get(target)
	if err != nil {
		panic(err)
	}
	res.Body.Close()
}

func train(wantNo, from, to, date string, m map[string]string) string {
	target := trainUrl
	target = strings.ReplaceAll(target, "#{from}", from)
	target = strings.ReplaceAll(target, "#{to}", to)
	target = strings.ReplaceAll(target, "#{date}", date)
	resp, err := client.Get(target)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	jsonStr, err := io.ReadAll(resp.Body)

	i := &Info{}
	err = json.Unmarshal(jsonStr, i)
	if err != nil {
		panic(err)
	}
	extraInfo := func(info string) (id, no, s, o, t string) {
		infos := strings.Split(info, "|")
		return infos[2], infos[3], infos[32], infos[31], infos[30]
	}

	for _, v := range i.Data.Result {
		id, no, s, o, t := extraInfo(v)
		if no == wantNo {
			fmt.Println(fmt.Sprintf("train :%s [%s:%s] \ttwo：%s\tone：%s\tspecial：%s\tid:%s\t\n  ", no, m[from], m[to], s, o, t, id))
			return id
		}
	}
	return ""
}

func pass(id, from, to, date string) []string {
	target := passUrl
	target = strings.ReplaceAll(target, "#{id}", id)
	target = strings.ReplaceAll(target, "#{from}", from)
	target = strings.ReplaceAll(target, "#{to}", to)
	target = strings.ReplaceAll(target, "#{date}", date)
	res, err := client.Get(target)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	interval := &IntervalInfo{}
	_ = json.Unmarshal(body, interval)
	p := []string{}
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
	resp := string(body)
	citys := strings.Split(resp, "@")
	citys = citys[1:]
	m := make(map[string]string, len(citys)*2)
	for c := range citys {
		fields := strings.Split(citys[c], "|")
		m[fields[1]] = fields[2]
		m[fields[2]] = fields[1]
	}
	return m
}
