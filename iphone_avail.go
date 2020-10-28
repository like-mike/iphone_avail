package main

import (
	"log"

	"github.com/like-mike/iphone_avail/item"

	"github.com/like-mike/iphone_avail/check_avail"
	"github.com/like-mike/iphone_avail/config"
	"github.com/like-mike/iphone_avail/email"
)

func main() {
	err := config.Init()
	if err != nil {
		_ = email.SendRecepients(0, config.Env.ErrorRecepients, "Error", err.Error())
		log.Fatal(err)
	}

	items, err := item.GetItems()
	if err != nil {
		_ = email.SendRecepients(0, config.Env.ErrorRecepients, "Error", err.Error())
		log.Fatal(err)
	}

	errCh := make(chan error, len(items))

	for i, item := range items {
		go check_avail.CheckAvail(int64(i), item, errCh)
	}

	for i := 0; i < len(items); i++ {
		select {
		case err = <-errCh:
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
				_ = email.SendRecepients(int64(len(items))+1, config.Env.ErrorRecepients, "Error", err.Error())
			}
		}
	}

}
