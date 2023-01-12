package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	solrClient = "http://localhost:8983/api/collections"
)

var client = http.Client{
	Transport: &http.Transport{},
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type SolrCreateCommand struct {
	Data string `json:"data"`
}

type CreateSchemaCommand struct {
	Create Create `json:"create"`
}
type Create struct {
	Name              string `json:"name"`
	NumShards         int    `json:"numShards"`
	ReplicationFactor int    `json:"replicationFactor"`
}

func main() {
	ping, _ := Ping(solrClient)
	fmt.Print(fmt.Sprintf("Solr instance status %v", ping))

	createSchema()

}

func createSchema() {
	cS := CreateSchemaCommand{
		Create: Create{
			Name:              "users",
			NumShards:         1,
			ReplicationFactor: 1,
		},
	}
	marshal, err := json.Marshal(cS)
	if err != nil {
		return
	}

	payload, err := json.Marshal(SolrCreateCommand{Data: string(marshal)})
	if err != nil {
		return
	}

	r, err := http.NewRequest("POST", fmt.Sprintf(solrClient), bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		panic(res.Status)
	}
}

func Ping(domain string) (int, error) {
	url := "http://" + domain
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return resp.StatusCode, nil
}
