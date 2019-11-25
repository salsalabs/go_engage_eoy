package main

import (
	//"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("gorm-activity-copy", "A command-line app to copy fundraising activities to SQLite via GORM")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&goengage.Fundraise{})
	db.AutoMigrate(&goengage.Transaction{})
	db.AutoMigrate(&goengage.Supporter{})
	db.AutoMigrate(&goengage.Contact{})
	db.AutoMigrate(&goengage.CustomFieldValue{})

	payload := goengage.ActivityRequestPayload{
		Type:         goengage.FundraiseType,
		ModifiedFrom: "2010-09-01T00:00:00.000Z",
		ModifiedTo:   "2020-09-01T00:00:00.000Z",
		Offset:       0,
		Count:        e.Metrics.MaxBatchSize,
	}
	rqt := goengage.ActivityRequest{
		Header:  goengage.RequestHeader{},
		Payload: payload,
	}
	var resp goengage.FundraiseResponse
	n := goengage.NetOp{
		Host:     e.Host,
		Endpoint: goengage.SearchActivity,
		Method:   goengage.SearchMethod,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	count := int32(rqt.Payload.Count)
	for count > 0 {
		err := n.Do()
		if err != nil {
			panic(err)
		}
		count = int32(len(resp.Payload.Activities))
		rqt.Payload.Offset = rqt.Payload.Offset + count
		fmt.Printf("%-20s %-36s %-10s %-10s %7s %7s %5s\n",
			"Name",
			"ActivityDate",
			"ActivityType",
			"DonationType",
			"TotalReceivedAmount",
			"RecurringAmount",
			"OneTimeAmount")

		for _, r := range resp.Payload.Activities {
			r.Year = r.ActivityDate.Year()
			r.Month = int(r.ActivityDate.Month())
			r.Day = r.ActivityDate.Day()
			db.Create(r)

			if len(r.Transactions) != 0 {
				for _, c := range r.Transactions {
					db.Create(&c)
				}

				s := goengage.Supporter{
					SupporterID: r.SupporterID,
				}
				db.Where("supporter_id = ?", r.SupporterID).First(&s)
				if s.CreatedDate == nil {
					t, err := goengage.FetchSupporter(e, r.SupporterID)
					if err != nil {
						log.Fatal(err)
					}
					if t == nil {
						x := time.Now()
						s.CreatedDate = &x
					} else {
						s = *t
					}
					db.Create(&s)
				} else {
					db.First(&s)
				}
				name := fmt.Sprintf(`"%v %v"`, s.FirstName, s.LastName)
				fmt.Printf("%-20s %04d-%02d-%02d %-10s %-10s %7.2f %7.2f %7.2f\n",
					name,
					r.Year,
					r.Month,
					r.Day,
					r.ActivityType,
					r.DonationType,
					r.TotalReceivedAmount,
					r.RecurringAmount,
					r.OneTimeAmount)
			}
		}
	}
}
