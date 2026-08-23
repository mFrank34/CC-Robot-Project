package retrieve

import (
	"Robot-Project/external/router"
	"Robot-Project/internal/model"
	"encoding/json"
	"fmt"
)

func Status(host, targetId string) (model.Status, error) {
	body, err := router.Get(host, "/id/"+targetId+"/status")
	if err != nil {
		fmt.Println("Error capturing Bot Status: ", err)
		return model.Status{}, err
	}

	var result model.Status

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing Json request", err)
		return model.Status{}, err
	}

	return result, nil
}

func Ids(host string) ([]string, error) {
	body, err := router.Get(host, "/ids")
	if err != nil {
		fmt.Println("Error making requires for ids:", err)
		return []string{}, err
	}

	var bot model.Ids
	if err := json.Unmarshal(body, &bot); err != nil {
		fmt.Println("Error parsing Json request", err)
		return []string{}, err
	}

	if len(bot.Ids) == 0 {
		return []string{}, fmt.Errorf("no bot ids found")
	}

	return bot.Ids, nil
}
