package eoy

import (
	"log"
	"os"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	goengage "github.com/salsalabs/goengage/pkg"
)

//Runtime contains the variables that we need to run this application.
type Runtime struct {
	Env         *goengage.Environment
	DB          *gorm.DB
	Log         *log.Logger
	Channels    []chan goengage.Fundraise
	Spreadsheet *excelize.File
	CountStyle  int
	ValueStyle  int
	KeyStyle    int
}

//ActivityForm contains a basic set of values for an activity page.
type ActivityForm struct {
	ID          string
	Name        string
	CreatedDate *time.Time
}

//Stat is used to store the usual statistics about a set of donations.
type Stat struct {
	// supporterID, activityPageID, Year, Year-Month
	ID              string
	AllCount        int32
	AllAmount       float64
	OneTimeCount    int32
	OneTimeAmount   float64
	RecurringCount  int32
	RecurringAmount float64
	OfflineCount    int32
	OfflineAmount   float64
	RefundsCount    int32
	RefundsAmount   float64
	Largest         float64
	Smallest        float64
	CreatedDate     *time.Time
}

//Month is used to provide a primary key for storing stats by month.
type Month struct {
	//ID is YYYY-MM
	ID          string
	Year        int
	Month       int
	CreatedDate *time.Time
}

//Year is used to provide a primary key for storing stats by year.
type Year struct {
	ID          int
	CreatedDate *time.Time
}

//Decorator describes how all of the sheets should behave.  An "index" is the
//offset from the start of the results returned by the database, where descriptive
//values (from the keys) are followed by a stats object.
type Decorator interface {
	//Index is the offset from the beginning of the results
	//returned by the database.  Keys are first, stats follow.
	Value(index int) interface{}
	Style(index int) string
	Axis(row, index int) string
}

//Sheet contains the stuff that we need to create and populate a sheet
//in the EOY spreadsheet.
type Sheet struct {
	Name   string
	Keys   []string
	Titles []string
}

//NewRuntime creates a runtime object and initializes the rt.
func NewRuntime(e *goengage.Environment, db *gorm.DB, channels []chan goengage.Fundraise) *Runtime {
	w, err := os.Create("eoy.log")
	if err != nil {
		log.Panic(err)
	}
	s := excelize.NewFile()
	countStyle, _ := s.NewStyle(`{"number_format": 3}`)
	valueStyle, _ := s.NewStyle(`{"number_format": 3}`)
	keyStyle, _ := s.NewStyle(`{"number_format": 0}`)

	rt := Runtime{
		Env:         e,
		DB:          db,
		Log:         log.New(w, "EOY: ", log.LstdFlags),
		Channels:    channels,
		Spreadsheet: s,
		CountStyle:  countStyle,
		ValueStyle:  valueStyle,
		KeyStyle:    keyStyle,
	}

	return &rt
}
