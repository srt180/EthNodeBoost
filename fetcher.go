package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

var search string

func init() {
	flag.StringVar(&search, "search", "", "search string for ethernodes.org query")
}

var urlFmt = `https://ethernodes.org/data?draw=1&columns[0][data]=id&columns[0][name]=&columns[0][searchable]=true&columns[0][orderable]=true&columns[0][search][value]=&columns[0][search][regex]=false&columns[1][data]=host&columns[1][name]=&columns[1][searchable]=true&columns[1][orderable]=true&columns[1][search][value]=&columns[1][search][regex]=false&columns[2][data]=isp&columns[2][name]=&columns[2][searchable]=true&columns[2][orderable]=true&columns[2][search][value]=&columns[2][search][regex]=false&columns[3][data]=country&columns[3][name]=&columns[3][searchable]=true&columns[3][orderable]=true&columns[3][search][value]=&columns[3][search][regex]=false&columns[4][data]=client&columns[4][name]=&columns[4][searchable]=true&columns[4][orderable]=true&columns[4][search][value]=&columns[4][search][regex]=false&columns[5][data]=clientVersion&columns[5][name]=&columns[5][searchable]=true&columns[5][orderable]=true&columns[5][search][value]=&columns[5][search][regex]=false&columns[6][data]=os&columns[6][name]=&columns[6][searchable]=true&columns[6][orderable]=true&columns[6][search][value]=&columns[6][search][regex]=false&columns[7][data]=lastUpdate&columns[7][name]=&columns[7][searchable]=true&columns[7][orderable]=true&columns[7][search][value]=&columns[7][search][regex]=false&columns[8][data]=inSync&columns[8][name]=&columns[8][searchable]=true&columns[8][orderable]=true&columns[8][search][value]=&columns[8][search][regex]=false&order[0][column]=0&order[0][dir]=asc&start=%d&length=%d&search[value]=%s&search[regex]=false&_=%d`

type Node struct {
	Id            string    `json:"id"`
	Host          string    `json:"host"`
	Port          int       `json:"port"`
	Client        string    `json:"client"`
	ClientVersion string    `json:"clientVersion"`
	Os            string    `json:"os"`
	LastUpdate    time.Time `json:"lastUpdate"`
	Country       string    `json:"country"`
	InSync        int       `json:"inSync"`
	Isp           string    `json:"isp"`
}

type Response struct {
	Draw            int    `json:"draw"`
	RecordsTotal    int    `json:"recordsTotal"`
	RecordsFiltered int    `json:"recordsFiltered"`
	Data            []Node `json:"data"`
}

func queryNodes() chan string {
	nodeChan := make(chan string, 1)

	go func() {
		total := 100
		length := 10

		now := time.Now()

		start := 0
		for ; start < total; start += length {
			url := fmt.Sprintf(urlFmt, start, length, search, now.UnixMilli())

			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}

			r := Response{}
			if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
				panic(err)
			}

			total = r.RecordsTotal

			for _, n := range r.Data {
				node := n
				enodeURL := fmt.Sprintf("enode://%s@%s:%d", node.Id, node.Host, node.Port)

				nodeChan <- enodeURL
			}

			resp.Body.Close()
		}
		close(nodeChan)
	}()

	return nodeChan

}
