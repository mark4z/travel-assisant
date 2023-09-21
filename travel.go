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
)

func init() {
	jar, _ := cookiejar.New(nil)
	client = http.Client{
		Jar: jar,
	}
}

func main() {
	m = mapper()
	no := "G1887"
	date := "2023-09-28"
	fromZh := "南京南"
	toZh := "郑州东"

	index(m[fromZh], fromZh, m[toZh], toZh, date)

	id := train(no, m[fromZh], m[toZh], date)
	stations := pass(id, m[fromZh], m[toZh], date)
	walk(stations, id, date)

	//ids := trainAll(m[fromZh], m[toZh], date)
	//for _, id := range ids {
	//	stations := pass(id, m[fromZh], m[toZh], date)
	//	fullWalk(stations, id, date)
	//}
}

func fullWalk(stations []Station, id string, date string) {
	trace(id, date, stations[0], stations[len(stations)-1])
}

func walk(stations []Station, id string, date string) {
	for i := 0; i < len(stations); i++ {
		trace(id, date, stations[0], stations[i])
		time.Sleep(500 * time.Microsecond)
	}
	fmt.Println()
	for i := 0; i < len(stations); i++ {
		trace(id, date, stations[i], stations[len(stations)-1])
		time.Sleep(500 * time.Microsecond)
	}
}

func index(from, fromZh, to, toZh, date string) {
	res, err := client.Get(fmt.Sprintf(indexUrl, fromZh, from, toZh, to, date))
	if err != nil {
		panic(err)
	}
	_ = res.Body.Close()
}

func trainAll(from, to, date string) (ans []string) {
	info := trainInfo(date, from, to)

	extraInfo := func(info string) (id, no, s, o, t string) {
		infos := strings.Split(info, "|")
		return infos[2], infos[3], infos[32], infos[31], infos[30]
	}

	for _, v := range info.Data.Result {
		id, no, _, _, _ := extraInfo(v)
		if strings.HasPrefix(no, "G") || strings.HasPrefix(no, "D") {
			ans = append(ans, id)
		}
	}
	return
}

func train(wantNo, from, to, date string) string {
	info := trainInfo(date, from, to)

	extraInfo := func(info string) (id, no, s, o, t string) {
		infos := strings.Split(info, "|")
		return infos[2], infos[3], infos[32], infos[31], infos[30]
	}

	for _, v := range info.Data.Result {
		id, no, s, o, t := extraInfo(v)
		if no == wantNo {
			fmt.Println(fmt.Sprintf("train :%s [%s:%s] \ttwo：%s\tone：%s\tspecial：%s\tid:%s\t\n  ", no, m[from], m[to], t, o, s, id))
			return id
		}
	}
	return ""
}

func trainInfo(date string, from string, to string) *Info {
	resp, err := client.Get(fmt.Sprintf(trainUrl, date, from, to))
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
	return info
}

func trace(wantNo, date string, fromS, toS Station) {
	from := m[fromS.StationName]
	to := m[toS.StationName]
	info := trainInfo(date, from, to)
	extraInfo := func(info string) (id, no, special, one, two, fromS, toS, from, to, arrive, start string) {
		infos := strings.Split(info, "|")
		return infos[2], infos[3], infos[32], infos[31], infos[30], infos[4], infos[5], infos[6], infos[7], infos[8], infos[9]
	}
	var found bool
	for _, v := range info.Data.Result {
		id, no, s, o, t, _, _, _, _, _, _ := extraInfo(v)
		if id == wantNo {
			found = true
			fmt.Println(fmt.Sprintf("train :%s [%s:%s] \ttwo：%s\tone：%s\tspecial：%s\tid:%s\t\n  ", no, fromS.StationName, toS.StationName, t, o, s, id))
		}
	}
	if found {
		return
	}
}

func pass(id, from, to, date string) []Station {
	res, err := client.Get(fmt.Sprintf(passUrl, id, from, to, date))
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	interval := &IntervalInfo{}
	_ = json.Unmarshal(body, interval)
	return interval.Data.Data
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
