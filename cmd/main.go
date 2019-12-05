package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	eoy "github.com/salsalabs/goengage-eoy/pkg"
	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	sleepDuration   = "10s"
	defaultTimezone = "Eastern"
)

type actor func(rt *eoy.Runtime, c chan goengage.Fundraise) (err error)

func main() {
	t := time.Now()
	y := t.Year()
	yearText := fmt.Sprintf("Year to use for reporting, default is %d", y)
	timezoneText := fmt.Sprintf("Choose '%v' (the default), 'Central', 'Mountain', 'Pacific' or 'Alaska'", defaultTimezone)
	var (
		app   = kingpin.New("Engage EOY Report", "A command-line app to create an Engage EOY")
		login = app.Flag("login", "YAML file with API token").Required().String()
		// org      = app.Flag("org", "Organization name (for output file)").Required().String()
		year     = app.Flag("year", yearText).Default(strconv.Itoa(y)).Int()
		topLimit = app.Flag("top", "Number in top donors sheet").Default("20").Int()
		timezone = app.Flag("timezone", timezoneText).Default(defaultTimezone).String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatal(err)
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
	db.AutoMigrate(&eoy.ActivityForm{})
	db.AutoMigrate(&eoy.Stat{})
	db.AutoMigrate(&eoy.Year{})
	db.AutoMigrate(&eoy.Month{})

	functions := []actor{
		eoy.Activity,
		eoy.Form,
		eoy.Stats,
		eoy.Supporter,
		eoy.Transaction,
		eoy.Date,
	}
	//Channel used by the downstream processors...
	var channels []chan goengage.Fundraise
	for i := 0; i < len(functions); i++ {
		c := make(chan goengage.Fundraise, 100)
		channels = append(channels, c)
	}
	var orgLocation string
	switch *timezone {
	case "Eastern":
		orgLocation = "America/New York"
	case "Central":
		orgLocation = "America/Chicago"
	case "Mountain":
		orgLocation = "America/Denver"
	case "Pacific":
		orgLocation = "America/Los Angeles"
	case "Alaska":
		orgLocation = "America/Nome"
	}

	done := make(chan bool)
	rt := eoy.NewRuntime(e, db, channels, *year, *topLimit, orgLocation)
	fmt.Printf("Runtime is %+v\n", rt)
	var wg sync.WaitGroup
	for i := range functions {
		go (func(i int, rt *eoy.Runtime, wg *sync.WaitGroup) {
			wg.Add(1)
			c := rt.Channels[i]
			r := functions[i]
			err := r(rt, c)
			if err != nil {
				rt.Log.Panic(err)
			}
			wg.Done()
		})(i, rt, &wg)
	}
	go (func(rt *eoy.Runtime, wg *sync.WaitGroup, done chan bool) {
		wg.Add(1)
		err := eoy.Drive(rt, done)
		if err != nil {
			rt.Log.Panic(err)
		}
		wg.Done()
	})(rt, &wg, done)
	<-done
	//d, _ := time.ParseDuration(sleepDuration)
	//time.Sleep(d)
	log.Printf("Waiting for tasks to complete.")
	wg.Wait()
	rt.Log.Printf("All tasks are complete.  Time to build the output.")
	// fmt.Println("Harvest start")
	// fn := fmt.Sprintf("%v %d EOY.xlsx", *org, *year)
	// err = rt.Harvest(fn)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Harvest end")
}
