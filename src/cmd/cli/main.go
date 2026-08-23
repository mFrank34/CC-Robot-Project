package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type getIDs struct {
	Ids []string `json:"ids"`
}

func request(host string, subroute string) ([]byte, error) {
	resp, err := http.Get(host + subroute)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	host := "https://cc.frankslab.uk"

	body, err := request(host, "/ids")
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return
	}

	var result getIDs
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	fmt.Println("Current found ids")
	for i, id := range result.Ids {
		fmt.Printf("ID [%d]: %s\n", i, id)
	}

}
