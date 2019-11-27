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
		rt.Log.Printf("%v Activity\n", r.ActivityID)
		rt.DB.Create(&r)
	}
	rt.Log.Println("Activity: end")
	return nil
}
