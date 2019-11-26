package eoy

import (
	goengage "github.com/salsalabs/goengage/pkg"
)

//Activity reads a channel of activities to retrieve ActivityIDs.  Those
//are used to populate the Activity table in the database.
func Activity(rt *Runtime, c chan goengage.Fundraise) (err error) {
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		r.Year = r.ActivityDate.Year()
		r.Month = int(r.ActivityDate.Month())
		r.Day = r.ActivityDate.Day()
		rt.DB.Create(r)
	}
	return nil
}
