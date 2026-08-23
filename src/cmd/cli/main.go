package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"Robot-Project/external/retrieve"
	"Robot-Project/external/send"
	"Robot-Project/internal/model"
)

const host = "https://cc.frankslab.uk"
const from = "Michael"

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// 1. List bots
		ids, err := retrieve.Ids(host)
		if err != nil {
			fmt.Println("Error fetching bot ids:", err)
			continue
		}
		if len(ids) == 0 {
			fmt.Println("No bots found. Retrying...")
			continue
		}

		fmt.Println("\nAvailable bots:")
		for i, id := range ids {
			fmt.Printf("  [%d] %s\n", i, id)
		}

		// 2. Select a bot
		fmt.Print("Select bot (index): ")
		selection, _ := reader.ReadString('\n')
		selection = strings.TrimSpace(selection)

		var idx int
		if _, err := fmt.Sscanf(selection, "%d", &idx); err != nil || idx < 0 || idx >= len(ids) {
			fmt.Println("Invalid selection, try again.")
			continue
		}
		targetId := ids[idx]

		// 3. Issue a command
		fmt.Print("Enter command (e.g. forward, dig, drop): ")
		cmdInput, _ := reader.ReadString('\n')
		cmdInput = strings.TrimSpace(cmdInput)

		fmt.Print("Enter arg (leave blank if none): ")
		argInput, _ := reader.ReadString('\n')
		argInput = strings.TrimSpace(argInput)

		cmd := model.Command(cmdInput)

		respBody, statusCode, err := send.Message(host, targetId, cmd, argInput, from)
		if err != nil {
			fmt.Println("Error sending command:", err)
			continue
		}
		fmt.Printf("[message] Status: %d, Response: %s\n", statusCode, respBody)

		// 4. Check status
		status, err := retrieve.Status(host, targetId)
		if err != nil {
			fmt.Println("Error fetching status:", err)
			continue
		}
		fmt.Printf("[status] Fuel: %d, Inventory: %v\n", status.Fuel, status.Inventory)
	}
}
