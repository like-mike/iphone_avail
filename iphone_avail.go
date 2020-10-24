package main

import (
	"log"

	"github.com/like-mike/iphone_avail/check_avail"
	"github.com/like-mike/iphone_avail/config"
	"github.com/like-mike/iphone_avail/email"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	errCh := make(chan error, 3)

	go check_avail.CheckAvail("MGLP3LL/A", "UNLOCKED/US", 98031, errCh) // Silver 128GB - Unlocked
	go check_avail.CheckAvail("MGKH3LL/A", "TMOBILE/US", 98031, errCh)  // Silver 128gb - TMO

	go check_avail.CheckAvail("MGKU3LL/A", "TMOBILE/US", 98031, errCh) // Blue 512gb - TMO

	for i := 0; i < 3; i++ {
		select {
		case err = <-errCh:
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
				_ = email.SendRecepients(config.Env.Recepients, "Error", err.Error())
			}
		}
	}

}
