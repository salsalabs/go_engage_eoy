package eoy

import (
	"log"
	"math"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//Stats reads a channel of Stats to retrieve ActivityIDs.  Those
//are used to populate the Activity table in the database.
func Stats(rt *Runtime, c chan goengage.Fundraise) (err error) {
	log.Println("Stats: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		log.Printf("%v Stats\n", r.ActivityID)

		g := GivingStat{}
		rt.DB.Where("id = ?", r.ActivityID).First(&g)
		if g.CreatedDate == nil {
			g.ID = r.ActivityID
			t := time.Now()
			g.CreatedDate = &t
			rt.DB.Create(&g)
		}
		for _, t := range r.Transactions {
			g.AllCount++
			g.AllAmount = g.AllAmount + t.Amount
			switch r.DonationType {
			case goengage.OneTime:
				g.OneTimeCount++
				g.OneTimeAmount += t.Amount
			case goengage.Recurring:
				g.RecurringCount++
				g.RecurringAmount += t.Amount
			}
			switch t.Type {
			case goengage.Refund:
				g.RefundsCount++
				g.RefundsAmount += t.Amount
			}
			//OfflineCount    int32
			//OfflineAmount   float64
			g.Largest = math.Max(g.Largest, t.Amount)
			g.Smallest = math.Min(g.Smallest, t.Amount)

			rt.DB.Model(&g).Updates(&g)
		}
	}
	log.Println("Stats: end")

	return nil
}
