package eoy

import "time"

//Stat is used to store the usual statistics about a set of donations.
type Stat struct {
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

//FieldOrder defines integer values for field order in the Stat record.
// See https://yourbasic.org/golang/iota/
type FieldOrder int

//Field orders for a state record.
const (
	ID FieldOrder = iota
	AllCount
	AllAmount
	OneTimeCount
	OneTimeAmount
	RecurringCount
	RecurringAmount
	OfflineCount
	OfflineAmount
	RefundsCount
	RefundsAmount
	Largest
	Smallest
	StatFieldCount
)

//Header returns the header from a stat record.
func (s Stat) Header(i int) string {
	switch FieldOrder(i) {
	case ID:
		return "Internal Key"
	case AllCount:
		return "All Count"
	case AllAmount:
		return "All Amount"
	case OneTimeCount:
		return "OneTime Count"
	case OneTimeAmount:
		return "OneTime Amount"
	case RecurringCount:
		return "Recurring Count"
	case RecurringAmount:
		return "Recurring Amount"
	case OfflineCount:
		return "Offline Count"
	case OfflineAmount:
		return "Offline Amount"
	case RefundsCount:
		return "Refunds Count"
	case RefundsAmount:
		return "Refunds Amount"
	case Largest:
		return "Largest"
	case Smallest:
		return "Smallest"
	}
	return ""
}

//Value returns the field value from a stat record.  Note that the field values
//match the order of fields in Stat.
func (s Stat) Value(i int) interface{} {
	switch FieldOrder(i) {
	case ID:
		return s.ID
	case AllCount:
		return s.AllCount
	case AllAmount:
		return s.AllAmount
	case OneTimeCount:
		return s.OneTimeCount
	case OneTimeAmount:
		return s.OneTimeAmount
	case RecurringCount:
		return s.RecurringCount
	case RecurringAmount:
		return s.RecurringAmount
	case OfflineCount:
		return s.OfflineCount
	case OfflineAmount:
		return s.OfflineAmount
	case RefundsCount:
		return s.RefundsCount
	case RefundsAmount:
		return s.RefundsAmount
	case Smallest:
		return s.Smallest
	case Largest:
		return s.Largest
	}
	return nil
}

//Style returns an Excel style code from a stat record.
func (s Stat) Style(rt *Runtime, i int) int {
	switch FieldOrder(i) {
	case ID:
		return rt.KeyStyle
	case AllCount:
		return rt.CountStyle
	case AllAmount:
		return rt.ValueStyle
	case OneTimeCount:
		return rt.CountStyle
	case OneTimeAmount:
		return rt.ValueStyle
	case RecurringCount:
		return rt.CountStyle
	case RecurringAmount:
		return rt.ValueStyle
	case OfflineCount:
		return rt.CountStyle
	case OfflineAmount:
		return rt.ValueStyle
	case RefundsCount:
		return rt.CountStyle
	case RefundsAmount:
		return rt.ValueStyle
	case Smallest:
		return rt.ValueStyle
	case Largest:
		return rt.ValueStyle
	}
	return -1
}

//Fill adds Excel cells for stat values using row and column for position.
func (s Stat) Fill(rt *Runtime, sheetName string, row, col int) int {
	for i := AllCount; i < StatFieldCount; i++ {
		f := int(i)
		v := s.Value(f)
		y := s.Style(rt, f)
		rt.Cell(sheetName, row, f, v, y)
	}
	row++
	return row
}
