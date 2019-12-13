package eoy

import (
	"time"
)

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
		return "Count"
	case AllAmount:
		return "Amount"
	case OneTimeCount:
		return "Count"
	case OneTimeAmount:
		return "Amount"
	case RecurringCount:
		return "Count"
	case RecurringAmount:
		return "Amount"
	case OfflineCount:
		return "Count"
	case OfflineAmount:
		return "Amount"
	case RefundsCount:
		return "Count"
	case RefundsAmount:
		return "Amount"
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

//FieldHeader returns the field header from a stat record.  Returns
//header text and the number of columns.  A null name indicates
//that the header does not have a column of its own.
func (s Stat) FieldHeader(i int) (name *string, c int) {
	t := ""
	switch FieldOrder(i) {
	case AllCount:
		t = "All"
		c = 2
	case OneTimeCount:
		t = "One Time"
		c = 2
	case RecurringCount:
		t = "Recurring"
		c = 2
	case OfflineCount:
		t = "Offline"
		c = 2
	case RefundsCount:
		t = "Refunds"
		c = 2
	case Largest:
		t = "Overall"
		c = 2
	}
	if len(t) > 0 {
		name = &t
	}
	return name, c
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
		rt.Cell(sheetName, row, col+f-int(AllCount), v, y)
	}
	row++
	return row
}

//Headers adds Excel cells has headers for each stat type (all, one time, recurring, etc.)
func (s Stat) Headers(rt *Runtime, sheetName string, row, col int) int {
	for i := AllCount; i < StatFieldCount; i++ {
		f := int(i)
		n, c := s.FieldHeader(f)
		colOffset := col + f - int(AllCount)
		if n != nil {
			rt.Cell(sheetName, row, col+f-1, *n, rt.HeaderStyle)
			if c > 1 {
				left := Axis(row, colOffset)
				right := Axis(row, colOffset+c-1)
				err := rt.Spreadsheet.MergeCell(sheetName, left, right)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	row++
	return row
}
