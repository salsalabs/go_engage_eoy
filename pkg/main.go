package eoy

//GivingStat is used to store the usual statistics about a set of donations.
type GivingStat struct {
	// supporterID, activityPageID, Year, Year-Month
	Key             string
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
}
