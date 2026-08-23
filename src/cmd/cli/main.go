package main

import (
	"encoding/json"
	"fmt"

	"Robot-Project/external/router"
	"Robot-Project/external/send"
	"Robot-Project/internal/model"
)

func main() {
	host := "https://cc.frankslab.uk"

	// 1. Fetch the IDs
	body, err := router.Get(host, "/ids")
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	var result model.Ids
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

	targetId := result.Ids[2]

	// 2. Send the message
	msgResp, statusCode, err := send.Message(host, targetId, model.Forward, "", "Michael")
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Printf("\n[message] Status: %d\n", statusCode)
	fmt.Printf("[message] Response: %s\n", string(msgResp))

	// 3. Fetch the bot's status
	statusResp, err := router.Get(host, "/id/"+targetId+"/status")
	if err != nil {
		fmt.Println("Error fetching status:", err)
		return
	}

	fmt.Printf("\n[status] Response: %s\n", string(statusResp))
}
