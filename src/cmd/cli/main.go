package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Ids struct {
	Ids []string `json:"ids"`
}

type Message struct {
	From    string    `json:"from"`
	Payload string    `json:"payload"`
	SentAt  time.Time `json:"timestamp"`
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

	// 1. Fetch the IDs
	body, err := request(host, "/ids")
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	var result Ids
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	if len(result.Ids) == 0 {
		fmt.Println("No IDs found.")
		return
	}

	fmt.Println("Current found ids:")
	for i, id := range result.Ids {
		fmt.Printf("ID [%d]: %s\n", i, id)
	}

	// 2. Prepare the POST message payload
	command := Message{
		From:    "Michael",
		Payload: "forward",
		SentAt:  time.Now(),
	}

	pack, err := json.Marshal(command)
	if err != nil {
		fmt.Println("Error parsing struct to JSON:", err)
		return
	}

	// 3. Construct the clean URL using fmt.Sprintf
	url := fmt.Sprintf("%s/%s/message", host, result.Ids[0])

	// 4. Send the POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pack))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 5. Read the POST response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Printf("\nStatus: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(responseBody))
}
