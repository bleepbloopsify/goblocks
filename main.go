package main

import (
	"fmt"

	"github.com/bleepbloopsify/blockchain/chain"
	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)

	// privkey := transaction.GenerateKey()
	// transaction.SaveKey("./private_key", privkey)
	// log.Info("Saved key")
	log.Info("Starting blockchain!")
	chain := chain.Genesis("Christopher")
	chain.AddData("Leon", 0)

	for i := 0; i < 10; i++ {
		u, _ := uuid.NewV4()
		chain.AddData(u.String(), 20)
		fmt.Println("Finished block:")
	}
	fmt.Print(chain)
}
