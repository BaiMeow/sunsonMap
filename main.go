package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	key string
)

type School struct {
	Address  string     `json:"-"`
	City     string     `json:"-"`
	Pos      [2]float64 `json:"pos"`
	Students []string   `json:"students"`
}

type Conf struct {
	Key   string `json:"key"`
	Index []struct {
		Class int    `json:"class"`
		Path  string `json:"path"`
	} `json:"index"`
}

func main() {
	file, err := ioutil.ReadFile("./conf.json")
	if err != nil {
		log.Fatal(err)
	}
	var conf Conf
	if err := json.Unmarshal(file, &conf); err != nil {
		log.Fatal(err)
	}
	key = conf.Key
	index := make(map[int]string)
	for _, v := range conf.Index {
		index[v.Class] = fmt.Sprintf("/data/data%d.json", v.Class)
		gen(v.Path, v.Class)
	}
	indexBytes, err := json.Marshal(&index)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(".\\public\\data\\index.json", indexBytes, 0666); err != nil {
		log.Fatal(err)
	}
}

func gen(path string, class int) {
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
	for {
		line, err := reader.Read()
		//忽略未知高校的学生
		if err != nil {
			if err == io.EOF {
				break
			}
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
	result := make(chan bool, 500)
	var num int
	for _, v := range data {
		num++
		go v.FillGeocode(result)
		time.Sleep(time.Millisecond * 50)
	}
	t := time.NewTimer(time.Second * 2)
	for num > 0 {
		select {
		case <-result:
			num--
			t.Reset(time.Second * 2)
		case <-t.C:
			log.Fatal("time out")
		}
	}
	t.Stop()
	SchoolsJsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal("json fail:", err)
	}
	savePath := fmt.Sprintf("public/data/data%d.json", class)
	if err := os.WriteFile(savePath, SchoolsJsonBytes, 0666); err != nil {
		log.Fatal(err)
	}
	log.Println(class, "班完成")
}

type amapResp struct {
	Status string `json:"status"`
	Pois   []struct {
		Location string `json:"location"`
	} `json:"pois"`
}

func (s *School) FillGeocode(result chan bool) {
	resp, err := http.Get(fmt.Sprintf("https://restapi.amap.com/v3/place/text?key=%s&keywords=%s", key, s.Address))
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
	result <- true
}
