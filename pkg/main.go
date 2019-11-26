package go_engage_eoy

import "time"

//GivingStat is used to store the usual statistics about a set of donations.
type GivingStat struct {
	Key string
	//Type is ONE_TIME, RECURRING, REFUND, OFFLINE, ALL
	Type string
	Max  struct {
		At     *time.Time `gorm:"max_at"`
		Amount float64    `gorm:"max_amount"`
	}
	Min struct {
		At     *time.Time `gorm:"min_at"`
		Amount float64    `gorm:"min_amount"`
	}
	Count   int32
	Average float64
	Sum     float64
}
