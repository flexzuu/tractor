package time

import "time"

type UnitType string

func (t UnitType) Enum() []string {
	return []string{
		string(UnitMinute),
		string(UnitHour),
		string(UnitDay),
		string(UnitWeek),
		string(UnitMonth),
		string(UnitYear),
	}
}

func (t UnitType) String() string {
	return string(t)
}

const (
	UnitMinute UnitType = "minute"
	UnitHour   UnitType = "hour"
	UnitDay    UnitType = "day"
	UnitWeek   UnitType = "week"
	UnitMonth  UnitType = "month"
	UnitYear   UnitType = "year"
)

type CronManager struct {
	Hello      string
	PeriodUnit UnitType
	Time       time.Time `field:"time"`
	Date       time.Time `field:"date"`
	Names      []string
	Numbers    []int
	Times      []time.Time `fields:"time"`
}

// func (c *CronManager) InspectorUI() view.Element {
// 	return view.El("div", view.Attrs{"class": "mx-4 flex"}, []view.Element{
// 		view.El("atom.Knob", nil, nil),
// 		view.El("atom.Slider", nil, nil),
// 	})
// }
