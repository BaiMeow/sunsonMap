package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	key, path string
)

type School struct {
	Address  string     `json:"-"`
	City     string     `json:"-"`
	Pos      [2]float64 `json:"pos"`
	Students []string   `json:"students"`
}

func main() {
	flag.StringVar(&key, "key", "", "高德地图api的key")
	flag.StringVar(&path, "path", "", "数据路径")
	flag.Parse()
	if key == "" || path == "" {
		log.Fatal("missing key or path")
	}
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(bufio.NewReader(file))
	var data = make(map[string]*School)
	//忽略标题行
	reader.Read()
	for i := 1; i <= 57; i++ {
		line, err := reader.Read()
		//忽略未知高校的学生
		if err != nil {
			panic(err)
		}
		if line[2] == "" {
			continue
		}
		if _, ok := data[line[2]]; !ok {
			if line[5] == "" {
				line[5] = line[2]
			}
			data[line[2]] = &School{
				City:    line[4],
				Address: line[5],
			}
		}
		data[line[2]].Students = append(data[line[2]].Students, line[1])
	}
	for _, v := range data {
		v.FillGeocode()
	}
	SchoolsJsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal("json fail:", err)
	}
	os.WriteFile("public/data/data.json", SchoolsJsonBytes, 0666)
}

type amapResp struct {
	Status string `json:"status"`
	Pois   []struct {
		Location string `json:"location"`
	} `json:"pois"`
}

func (s *School) FillGeocode() {
	resp, err := http.Get(fmt.Sprintf("https://restapi.amap.com/v3/place/text?key=%s&keywords=%s&city=%s", key, s.Address, s.City))
	if err != nil {
		log.Fatal("http fail", s.Address, err)
	}
	var respData amapResp
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		log.Fatal("json fail", err)
	}
	if respData.Status != "1" {
		return
	}
	var longitude, latitude float64
	if _, err := fmt.Sscanf(respData.Pois[0].Location, "%e,%e", &longitude, &latitude); err != nil {
		log.Fatal(err)
	}
	s.Pos = [2]float64{longitude, latitude}
}
