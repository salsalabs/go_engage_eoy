package eoy

import (
	"fmt"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//Dates reads a channel of activities to retrieve DatesIDs.  Those
//are used to populate the Dates table in the database.
func Dates(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Dates: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		rt.Log.Printf("%v Dates\n", r.ActivityID)
		y := Year{
			ID: r.Year,
		}
		rt.DB.FirstOrInit(&y, y)
		if y.CreatedDate == nil {
			t := time.Now()
			y.CreatedDate = &t
			rt.DB.Create(&y)
		}
		m := Month{
			ID:    fmt.Sprintf("%4d-%02d", r.Year, r.Month),
			Year:  r.Year,
			Month: r.Month,
		}
		rt.DB.FirstOrInit(&m, m)
		if m.CreatedDate == nil {
			t := time.Now()
			m.CreatedDate = &t

			m.Year = r.Year
			m.Month = r.Month
			rt.DB.Create(&m)
		}
	}
	rt.Log.Println("Dates: start")
	return nil
}
