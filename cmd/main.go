package main

import (
	"log"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	eoy "github.com/salsalabs/goengage-eoy/pkg"
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
	db.AutoMigrate(&eoy.ActivityForm{})
	db.AutoMigrate(&eoy.GivingStat{})

	var channels []chan goengage.Fundraise
	for i := 0; i < 5; i++ {
		c := make(chan goengage.Fundraise)
		channels = append(channels, c)
	}
	rt := eoy.NewRuntime(e, db, channels)
	functions := []func(rt *eoy.Runtime, c chan goengage.Fundraise) (err error) {
		eoy.Activity,
		eoy.Form,
		eoy.Stats,
		eoy.Supporter,
		eoy.Transaction,
	}
	var wg sync.WaitGroup
	for i, r := range functions {
		go (func(i int, rt *eoy.Runtime, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			c := rt.Channels[i]
			err := r(rt, c)
			if err != nil {
				rt.Log.Panic(err)
			}
		})(i, rt, &wg);
	}
}
