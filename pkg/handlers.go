package eoy

import (
	"fmt"
	"log"
	"math"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//Activity reads a channel of activities to retrieve ActivityIDs.  Those
//are used to populate the Activity table in the database.
func Activity(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Activity: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		if rt.GoodYear(r.ActivityDate) {
			rt.DB.Create(&r)
		}
	}
	rt.Log.Println("Activity: end")
	return nil
}

//Date reads a channel of activities. Those are used to populate
//the Dates table in the database.  Note that the selected year
//is not used to filter out records.
func Date(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Date: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		y := Year{
			ID: r.Year,
		}
		rt.DB.FirstOrInit(&y, y)
		if y.CreatedDate == nil {
			t := time.Now()
			y.CreatedDate = &t
			rt.DB.Create(&y)
		}
		m := Month{
			ID:    fmt.Sprintf("%4d-%02d", r.Year, r.Month),
			Year:  r.Year,
			Month: r.Month,
		}
		rt.DB.FirstOrInit(&m, m)
		if m.CreatedDate == nil {
			t := time.Now()
			m.CreatedDate = &t

			m.Year = r.Year
			m.Month = r.Month
			rt.DB.Create(&m)
		}
	}
	rt.Log.Println("Date: end")
	return nil
}

//Drive reads fundraising activities from Engage and hands them off to
//channels for downstream processing.
func Drive(rt *Runtime, done chan bool) (err error) {
	rt.Log.Println("Drive: start")

	payload := goengage.ActivityRequestPayload{
		Type:         goengage.FundraiseType,
		ModifiedFrom: "2001-09-01T00:00:00.000Z",
		//ModifiedTo:   "2020-09-01T00:00:00.000Z",
		Offset: 0,
		Count:  rt.Env.Metrics.MaxBatchSize,
	}
	rqt := goengage.ActivityRequest{
		Header:  goengage.RequestHeader{},
		Payload: payload,
	}
	var resp goengage.FundraiseResponse
	n := goengage.NetOp{
		Host:     rt.Env.Host,
		Endpoint: goengage.SearchActivity,
		Method:   goengage.SearchMethod,
		Token:    rt.Env.Token,
		Request:  &rqt,
		Response: &resp,
	}
	count := int32(rqt.Payload.Count)
	for count == int32(rqt.Payload.Count) {
		err := n.Do()
		if err != nil {
			return err
		}
		m := fmt.Sprintf("Drive: read from offset %6d / %6d\n", resp.Payload.Offset, resp.Payload.Total)
		rt.Log.Printf(m)
		if rqt.Payload.Offset%500 == 0 {
			log.Printf(m)
		}
		rqt.Payload.Offset = rqt.Payload.Offset + count
		count = int32(len(resp.Payload.Activities))
		for _, r := range resp.Payload.Activities {
			// Need this here so that the stats will work.
			r.Year = r.ActivityDate.Year()
			r.Month = int(r.ActivityDate.Month())
			r.Day = r.ActivityDate.Day()

			for _, c := range rt.Channels {
				c <- r
			}
		}
	}
	for _, c := range rt.Channels {
		close(c)
	}
	done <- true
	rt.Log.Println("Drive: end")
	return nil
}

//Form reads a channel of activities to retrieve ActivityIDs.  Those
//are used to populate the Activity table in the database.
func Form(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Form: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		if rt.GoodYear(r.ActivityDate) {
			s := ActivityForm{}
			rt.DB.Where("id = ?", r.ActivityFormID).First(&s)
			if s.CreatedDate == nil {
				s.ID = r.ActivityFormID
				s.Name = r.ActivityFormName
				t := time.Now()
				s.CreatedDate = &t
				rt.DB.Create(&s)
			}
		}
	}
	rt.Log.Println("Form: end")
	return nil
}

//update accepts a key and retrieves or creates a new stats record.  The record
//is updated with the current activity, then written back to the database.
func update(rt *Runtime, r goengage.Fundraise, key string) {
	g := Stat{}
	rt.DB.Where("id = ?", key).First(&g)
	if g.CreatedDate == nil {
		g.ID = key
		t := time.Now()
		g.CreatedDate = &t
		rt.DB.Create(&g)
	}
	amount := 0.0
	if r.WasImported {
		g.OfflineCount++
		g.OfflineAmount += r.OneTimeAmount
		amount = g.OfflineAmount
	} else {
		switch r.DonationType {
		case goengage.OneTime:
			g.OneTimeCount++
			g.OneTimeAmount += r.OneTimeAmount
			amount = g.OneTimeAmount
		case goengage.Recurring:
			g.RecurringCount++
			g.RecurringAmount += r.RecurringAmount
			amount = g.RecurringAmount
		}
		for _, t := range r.Transactions {
			switch t.Type {
			case goengage.Refund:
				g.RefundsCount++
				g.RefundsAmount += t.Amount
			}
		}
	}
	g.AllCount++
	g.AllAmount += amount
	g.Largest = math.Max(g.Largest, amount)
	if amount > 0.0 {
		if g.Smallest < 1.0 {
			g.Smallest = amount
		} else {
			g.Smallest = math.Min(g.Smallest, amount)
		}
	}
	rt.DB.Model(&g).Updates(&g)
}

//Stats reads a channel of Stats.  Those records are used to populate
//Stats records for each of the accumulation topics.
func Stats(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Stats: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		if rt.GoodYear(r.ActivityDate) {
			//Debug with Nicaragua
			if r.ActivityFormID == "dd4fc813-c019-4408-95dd-7bf708fcbac7" {
				fmt.Printf("%-10s %-36s %-36s %-20s %6.2f + %6.2f = %6.2f\n",
					r.ActivityFormName,
					r.ActivityID,
					r.ActivityDate,
					r.DonationType,
					r.OneTimeAmount,
					r.RecurringAmount,
					r.TotalReceivedAmount,
				)
				for _, t := range r.Transactions {
					fmt.Printf("%-10s %-36s %-36s Transaction: %-12s %-12s  %7.2f\n",
						r.ActivityFormName,
						r.ActivityID,
						r.ActivityDate,
						t.Reason,
						t.Type,
						t.Amount,
					)

				}
			}
			update(rt, r, r.ActivityID)
			update(rt, r, r.SupporterID)
			update(rt, r, r.ActivityFormID)
			t := fmt.Sprintf("%d", r.Year)
			update(rt, r, t)
			t = fmt.Sprintf("%d-%02d", r.Year, r.Month)
			update(rt, r, t)
		}
	}
	rt.Log.Println("Stats: end")

	return nil
}

//Supporter reads a channel of activities to retrieve supporterIDs.  Those
//are used to populate the supporter table in the database.
func Supporter(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Supporter: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		if rt.GoodYear(r.ActivityDate) {
			s := goengage.Supporter{
				SupporterID: r.SupporterID,
			}
			rt.DB.FirstOrInit(&s, s)

			// rt.DB.Where("supporter_id = ?", r.SupporterID).First(&s)
			if s.CreatedDate == nil {
				t, err := goengage.FetchSupporter(rt.Env, r.SupporterID)
				if err != nil {
					return err
				}
				if t == nil {
					x := time.Now()
					s.CreatedDate = &x
				} else {
					s = *t
				}
				rt.DB.Create(&s)
			}
		}
	}
	rt.Log.Println("Supporter: end")
	return nil
}

//Transaction reads a channel of activities to retrieve TransactionIDs.  Those
//are used to populate the Transaction table in the database.
func Transaction(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Transaction: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		if rt.GoodYear(r.ActivityDate) {
			if len(r.Transactions) != 0 {
				for _, c := range r.Transactions {
					c.ActivityID = r.ActivityID
					rt.DB.Create(&c)
				}
			}
		}
	}
	rt.Log.Println("Transaction: end")
	return nil
}
