package eoy

import (
	"log"

	goengage "github.com/salsalabs/goengage/pkg"
)

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
