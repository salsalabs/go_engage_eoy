package main

import (
	"fmt"
	"log"
	"os"

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
	db.AutoMigrate(&eoy.Stat{})
	db.AutoMigrate(&eoy.Year{})
	db.AutoMigrate(&eoy.Month{})

	var channels []chan goengage.Fundraise
	rt := eoy.NewRuntime(e, db, channels)
	fmt.Println("Harvest start")
	err = rt.Harvest("eoy_test.xlsx")
	if err != nil {
		panic(err)
	}
	fmt.Println("Harvest end")
}
