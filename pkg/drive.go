package eoy

import (
	"log"

	goengage "github.com/salsalabs/goengage/pkg"
)

//Drive reads fundraising activities from Engage and hands them off to
//channels for downstream processing.
func Drive(rt *Runtime, done chan bool) (err error) {
	log.Println("Drive: start")

	payload := goengage.ActivityRequestPayload{
		Type:         goengage.FundraiseType,
		ModifiedFrom: "2010-09-01T00:00:00.000Z",
		ModifiedTo:   "2020-09-01T00:00:00.000Z",
		Offset:       0,
		Count:        rt.Env.Metrics.MaxBatchSize,
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
	for count > 0 {
		err := n.Do()
		if err != nil {
			return err
		}
		count = int32(len(resp.Payload.Activities))
		rqt.Payload.Offset = rqt.Payload.Offset + count
		// fmt.Printf("%-36s %-10s %-10s %-10s %7s %7s %7s\n",
		// 	"ActivityID",
		// 	"ActivityDate",
		// 	"ActivityType",
		// 	"DonationType",
		// 	"Total",
		// 	"Recurring",
		// 	"OneTime")

		for _, r := range resp.Payload.Activities {
			for i, c := range rt.Channels {
				c <- r
				log.Printf("%v Drive: chan %d\n", r.ActivityID, i)

				// fmt.Printf("Drive: %-36s %04d-%02d-%02d %-10s %-10s %7.2f %7.2f %7.2f\n",
				// 	r.ActivityID,
				// 	r.Year,
				// 	r.Month,
				// 	r.Day,
				// 	r.ActivityType,
				// 	r.DonationType,
				// 	r.TotalReceivedAmount,
				// 	r.RecurringAmount,
				// 	r.OneTimeAmount)
			}
		}
	}
	for _, c := range rt.Channels {
		close(c)
	}
	done <- true
	log.Println("Drive: end")
	return nil
}
