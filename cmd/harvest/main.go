package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	t := time.Now()
	y := t.Year()
	yearText := fmt.Sprintf("Year to use for reporting, default is %d", y)
	var (
		app      = kingpin.New("Engage EOY Report", "A command-line app to create an Engage EOY")
		login    = app.Flag("login", "YAML file with API token").Required().String()
		org      = app.Flag("org", "Organization name (for output file)").Required().String()
		year     = app.Flag("year", yearText).Default(strconv.Itoa(y)).Int()
		topLimit = app.Flag("top", "Number in top donors sheet").Default("20").Int()
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
	rt.Year = *year
	rt.TopDonorLimit = *topLimit
	fmt.Println("Harvest start")
	fn := fmt.Sprintf("%v %d EOY.xlsx", *org, *year)
	err = rt.Harvest(fn)
	if err != nil {
		panic(err)
	}
	fmt.Println("Harvest end")
}
