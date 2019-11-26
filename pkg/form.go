package eoy

import (
	goengage "github.com/salsalabs/goengage/pkg"
)

//ActivityForm contains a basic set of values for an activity page.
type ActivityForm struct {
	Key  string
	Name string
}

//Form reads a channel of activities to retrieve ActivityIDs.  Those
//are used to populate the Activity table in the database.
func Form(rt *Runtime, c chan goengage.Fundraise) (err error) {
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		a := ActivityForm{
			Key:  r.ActivityFormID,
			Name: r.ActivityFormName,
		}
		rt.DB.Create(a)
	}
	return nil
}
