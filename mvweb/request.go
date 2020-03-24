package mvweb

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	mvwebApi = "https://mediathekviewweb.de/api/query"
)

func NewRequest(maxResults uint, queries ...Query) Request {
	//{"queries":[{"fields":["channel"],"query":"ard"},{"fields":["topic"],"query":"sendung"}],"sortBy":"timestamp","sortOrder":"desc","future":false,"offset":0,"size":15}
	return Request{
		Queries:   queries,
		SortBy:    "timestamp",
		SortOrder: "desc",
		Future:    false,
		Offset:    0,
		Size:      maxResults,
	}
}

func NewQuery(field, query string) Query {
	return Query{Fields: []string{field}, Query: query}
}

type Query struct {
	Fields []string `json:"fields"`
	Query  string   `json:"query"`
}

type Request struct {
	Queries   []Query `json:"queries"`
	SortBy    string  `json:"sortBy"`
	SortOrder string  `json:"sortOrder"`
	Future    bool    `json:"future"`
	Offset    uint    `json:"offset"`
	Size      uint    `json:"size"`
}

func (m *Request) Run() (Response, error) {
	var res Response
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(&m)
	if err != nil {
		return res, err
	}

	resp, err := http.Post(mvwebApi, "text/plain;charset=UTF-8", &buf)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&res)
	return res, err
}
