package eoy

import (
	"log"

	goengage "github.com/salsalabs/goengage/pkg"
)

//Transaction reads a channel of activities to retrieve TransactionIDs.  Those
//are used to populate the Transaction table in the database.
func Transaction(rt *Runtime, c chan goengage.Fundraise) (err error) {
	log.Println("Transaction: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		log.Printf("%v Transaction\n", r.ActivityID)

		if len(r.Transactions) != 0 {
			for _, c := range r.Transactions {
				rt.DB.Create(&c)
			}
		}
	}
	log.Println("Transaction: start")
	return nil
}
