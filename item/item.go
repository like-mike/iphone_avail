package item

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/like-mike/iphone_avail/config"
)

type Item struct {
	Serial  string `json:"serial"`
	Carrier string `json:"carrier"`
	Zip     int    `json:"zip"`
}

func GetItems() ([]Item, error) {
	// read file
	data, err := ioutil.ReadFile(config.Env.StaticPath + "/items.json")
	if err != nil {
		return nil, err
	}

	var items []Item
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("No items")
	}

	return items, nil
}
