package eoy

import (
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	goengage "github.com/salsalabs/goengage/pkg"
)

//Runtime contains the variables that we need to run this application.
type Runtime struct {
	Env      *goengage.Environment
	DB       *gorm.DB
	Log      *log.Logger
	Channels []chan goengage.Fundraise
}

//NewRuntime creates a runtime object and initializes the log.
func NewRuntime(e *goengage.Environment, db *gorm.DB, channels []chan goengage.Fundraise) *Runtime {
	w, err := os.Create("eoy.log")
	if err != nil {
		log.Panic(err)
	}

	rt := Runtime{
		Env:      e,
		DB:       db,
		Log:      log.New(w, "EOY: ", log.LstdFlags),
		Channels: channels,
	}

	return &rt
}

//GivingStat is used to store the usual statistics about a set of donations.
type GivingStat struct {
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
