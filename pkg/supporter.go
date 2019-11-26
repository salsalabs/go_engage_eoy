package eoy

import (
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//Supporter reads a channel of activities to retrieve supporterIDs.  Those
//are used to populate the supporter table in the database.
func Supporter(rt *Runtime, c chan goengage.Fundraise) (err error) {
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		s := goengage.Supporter{
			SupporterID: r.SupporterID,
		}
		rt.DB.Where("supporter_id = ?", r.SupporterID).First(&s)
		if s.CreatedDate == nil {
			t, err := goengage.FetchSupporter(rt.Env, r.SupporterID)
			if err != nil {
				return err
			}
			if t == nil {
				x := time.Now()
				s.CreatedDate = &x
			} else {
				s = *t
			}
			rt.DB.Create(&s)
		} else {
			rt.DB.First(&s)
		}
	}
	return nil
}
