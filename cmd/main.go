package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	eoy "github.com/salsalabs/goengage-eoy/pkg"
	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const sleepDuration = "10s"

type actor func(rt *eoy.Runtime, c chan goengage.Fundraise) (err error)

func main() {
	var (
		app   = kingpin.New("Engage EOY Report", "A command-line app to create an Engage EOY")
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
	db.AutoMigrate(&eoy.ActivityForm{})
	db.AutoMigrate(&eoy.GivingStat{})

	//Channel used by the downstream processors...
	var channels []chan goengage.Fundraise
	for i := 0; i < 5; i++ {
		c := make(chan goengage.Fundraise)
		channels = append(channels, c)
	}
	functions := []actor{
		eoy.Activity,
		eoy.Form,
		eoy.Stats,
		eoy.Supporter,
		eoy.Transaction,
	}

	done := make(chan bool)
	rt := eoy.NewRuntime(e, db, channels)
	var wg sync.WaitGroup
	for i := range functions {
		go (func(i int, rt *eoy.Runtime, wg *sync.WaitGroup) {
			wg.Add(1)
			c := rt.Channels[i]
			r := functions[i]
			err := r(rt, c)
			if err != nil {
				log.Panic(err)
			}
			wg.Done()
		})(i, rt, &wg)
	}
	go (func(rt *eoy.Runtime, wg *sync.WaitGroup, done chan bool) {
		wg.Add(1)
		err := eoy.Drive(rt, done)
		if err != nil {
			log.Panic(err)
		}
		wg.Done()
	})(rt, &wg, done)
	d, _ := time.ParseDuration(sleepDuration)
	time.Sleep(d)
	wg.Wait()
	log.Printf("All tasks are done.  Time to build the output.")
}
