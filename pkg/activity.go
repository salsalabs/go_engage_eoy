package eoy

import (
	goengage "github.com/salsalabs/goengage/pkg"
)

//Activity reads a channel of activities to retrieve ActivityIDs.  Those
//are used to populate the Activity table in the database.
func Activity(rt *Runtime, c chan goengage.Fundraise) (err error) {
	rt.Log.Println("Activity: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		rt.Log.Printf("Activity: %v\n", r.ActivityID)

		r.Year = r.ActivityDate.Year()
		r.Month = int(r.ActivityDate.Month())
		r.Day = r.ActivityDate.Day()
		rt.DB.Create(&r)
	}
	rt.Log.Println("Activity: end")
	return nil
}
