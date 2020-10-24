package check_avail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/like-mike/iphone_avail/item"

	"github.com/like-mike/iphone_avail/config"
	"github.com/like-mike/iphone_avail/email"
)

type PickupMessage struct {
	Body PickupBody `json:"body"`
}

type PickupBody struct {
	Stores []PickupStore `json:"stores"`
}

type PickupStore struct {
	StoreName         string      `json:"storeName"`
	PartsAvailability interface{} `json:"partsAvailability"`
}

type PickupStorePartsAvailability struct {
	StoreSelectionEnabled   bool   `json:"storeSelectionEnabled"`
	StorePickupQuote        string `json:"storePickupQuote"`
	StorePickupProductTitle string `json:"storePickupProductTitle"`
}

func CheckAvail(jobID int64, item *item.Item, errCh chan<- error) {

	// if true {
	// 	errCh <- fmt.Errorf("This is an error")
	// 	return
	// }

	url := fmt.Sprintf(`https://www.apple.com/shop/retail/pickup-message?pl=true&cppart=%s&parts.0=%s&location=%d`,
		item.Carrier,
		item.Serial,
		item.Zip,
	)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("accept", "*/*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errCh <- err
		return
	}

	if resp.StatusCode != 200 {
		errCh <- fmt.Errorf("Non-OK HTTP Status: %d", resp.StatusCode)
		return
	}

	var pickupMessage PickupMessage
	err = json.Unmarshal(body, &pickupMessage)
	if err != nil {
		errCh <- err
		return
	}

	availableStores := []PickupStorePartsAvailability{}
	stores := pickupMessage.Body.Stores
	for _, store := range stores {

		// unpack {}interface, assign to map
		partsAvail := store.PartsAvailability.(map[string]interface{})
		inv := make(map[string]PickupStorePartsAvailability)
		for key, value := range partsAvail {
			l := value.(map[string]interface{})
			inv[key] = PickupStorePartsAvailability{
				StoreSelectionEnabled:   l["storeSelectionEnabled"].(bool),
				StorePickupQuote:        l["storePickupQuote"].(string),
				StorePickupProductTitle: l["storePickupProductTitle"].(string),
			}
		}

		if inv[item.Serial].StoreSelectionEnabled {
			availableStores = append(availableStores, inv[item.Serial])
		}
	}

	if len(availableStores) > 0 {
		body := fmt.Sprintf("%s (%d stores)\n\n", availableStores[0].StorePickupProductTitle, len(availableStores))
		for _, store := range availableStores {
			body += fmt.Sprintf("%s is Avaiable %s \n", store.StorePickupProductTitle, store.StorePickupQuote)
		}

		// send email
		err = email.SendRecepients(jobID, config.Env.Recepients, "In Stock", body)
		if err != nil {
			errCh <- err
			return
		}
	}

	errCh <- nil
}
