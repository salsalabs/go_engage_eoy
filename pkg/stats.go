package eoy

import (
	"log"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//Stats reads a channel of Stats to retrieve ActivityIDs.  Those
//are used to populate the Activity table in the database.
func Stats(rt *Runtime, c chan goengage.Fundraise) (err error) {
	log.Println("Stats: start")
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		log.Printf("%v Stats\n", r.ActivityID)

		g := GivingStat{}
		rt.DB.Where("id = ?", r.ActivityID).First(&g)
		if g.CreatedDate == nil {
			g.ID = r.ActivityID
			t := time.Now()
			g.CreatedDate = &t
			rt.DB.Create(&g)
		}
	}
	log.Println("Stats: end")

	return nil
}
