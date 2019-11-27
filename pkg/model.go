package eoy

import (
	"fmt"
	"log"
	"math"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//ActivityForm contains a basic set of values for an activity page.
type ActivityForm struct {
	ID          string
	Name        string
	CreatedDate *time.Time
}

//Activity reads a channel of activities to retrieve ActivityIDs.  Those
//are used to populate the Activity table in the database.
func Activity(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Activity: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		rt.Log.Printf("%v Activity\n", r.ActivityID)
		rt.DB.Create(&r)
	}
	rt.Log.Println("Activity: end")
	return nil
}

//Dates reads a channel of activities to retrieve DatesIDs.  Those
//are used to populate the Dates table in the database.
func Dates(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Dates: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		rt.Log.Printf("%v Dates\n", r.ActivityID)
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
	rt.Log.Println("Dates: start")
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
		log.Printf("Drive: read %d from offset %6d\n", count, rqt.Payload.Offset)
		rqt.Payload.Offset = rqt.Payload.Offset + count
		count = int32(len(resp.Payload.Activities))
		for _, r := range resp.Payload.Activities {
			rt.Log.Printf("%v Drive\n", r.ActivityID)
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
		s := ActivityForm{}
		rt.Log.Printf("%v Form\n", r.ActivityID)
		rt.DB.Where("id = ?", r.ActivityFormID).First(&s)
		if s.CreatedDate == nil {
			s.ID = r.ActivityFormID
			s.Name = r.ActivityFormName
			t := time.Now()
			s.CreatedDate = &t
			rt.DB.Create(&s)
		}
	}
	rt.Log.Println("Form: end")
	return nil
}

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

//Supporter reads a channel of activities to retrieve supporterIDs.  Those
//are used to populate the supporter table in the database.
func Supporter(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Supporter: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		rt.Log.Printf("%v Supporter\n", r.ActivityID)

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
		rt.Log.Printf("%v Transaction\n", r.ActivityID)

		if len(r.Transactions) != 0 {
			for _, c := range r.Transactions {
				c.ActivityID = r.ActivityID
				rt.DB.Create(&c)
			}
		}
	}
	rt.Log.Println("Transaction: start")
	return nil
}
