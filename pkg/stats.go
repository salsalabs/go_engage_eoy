package eoy

import (
	"fmt"
	"math"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//update accepts a key and retrieves or creates a new stats record.  The record
//is updated with the current activity, then written back to the database.
func update(rt *Runtime, r goengage.Fundraise, key string) {
	g := GivingStat{}
	rt.DB.Where("id = ?", key).First(&g)
	if g.CreatedDate == nil {
		g.ID = key
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

//Stats reads a channel of Stats to retrieve ActivityIDs.  Those
//are used to populate the Activity table in the database.
func Stats(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Stats: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		rt.Log.Printf("%v Stats\n", r.ActivityID)
		update(rt, r, r.ActivityID)
		update(rt, r, r.SupporterID)
		update(rt, r, r.ActivityFormID)
		t := fmt.Sprintf("%d", r.Year)
		update(rt, r, t)
		t = fmt.Sprintf("%d-%02d", r.Year, r.Month)
		update(rt, r, t)
	}
	rt.Log.Println("Stats: end")

	return nil
}
