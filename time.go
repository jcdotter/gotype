// Copyright 2023 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.
// Author: jcdotter

package vals

import (
	"math"
	"time"
	"unsafe"
)

// ------------------------------------------------------------ /
// GOTYPE CUSTOM TYPE IMPLEMENTATION
// implementation of custom type of time.Time
// enabling seemless type conversion
// consolidated standard golang funtionality in single pkg
// and expanded transformation and computation functionality
// ------------------------------------------------------------ /

const (
	ISO8601N   = `2006-01-02 15:04:05.000000000`
	ISO8601    = `2006-01-02 15:04:05.000`
	SqlDate    = `2006-01-02T15:04:05Z`
	TimeFormat = `2006-01-02 15:04:05`
	DateFormat = `2006-01-02`
)

type TIME time.Time

// TIME returns gotype VALUE as gotype TIME
func (v VALUE) TIME() TIME {
	switch v.KIND() {
	case Int:
		return INT(*(*int)(v.ptr)).TIME()
	case Int8:
		return INT(*(*int8)(v.ptr)).TIME()
	case Int16:
		return INT(*(*int16)(v.ptr)).TIME()
	case Int32:
		return INT(*(*int32)(v.ptr)).TIME()
	case Int64:
		return INT(*(*int64)(v.ptr)).TIME()
	case Float32:
		return FLOAT(*(*float32)(v.ptr)).TIME()
	case Uint:
		return UINT(*(*uint)(v.ptr)).TIME()
	case Uint8:
		return UINT(*(*uint8)(v.ptr)).TIME()
	case Uint16:
		return UINT(*(*uint16)(v.ptr)).TIME()
	case Uint32:
		return UINT(*(*uint32)(v.ptr)).TIME()
	case Uint64:
		return UINT(*(*uint64)(v.ptr)).TIME()
	case Float64:
		return (*FLOAT)(v.ptr).TIME()
	case String:
		return (*STRING)(v.ptr).TIME()
	case Bytes:
		return STRING(*(*[]byte)(v.ptr)).TIME()
	case Time:
		return *(*TIME)(v.ptr)
	case Pointer:
		return v.Elem().TIME()
	}
	panic("cannot convert value to TIME")
}

// ------------------------------------------------------------ /
// TYPE CONVERSION FUNCTIONS
// implementation of functions to convert values to new types
// ------------------------------------------------------------ /

// Natural returns gotype TIME as golang time.Time
func (t TIME) Native() time.Time {
	return time.Time(t).UTC()
}

// Interface returns gotype TIME as a golang interface{}
func (t TIME) Interface() any {
	return t.Native()
}

// VALUE returns gotype TIME as gotype VALUE
func (t TIME) VALUE() VALUE {
	a := (any)(t)
	return *(*VALUE)(unsafe.Pointer(&a))
}

// Encode returns a gotype encoding of TIME
func (t TIME) Encode() ENCODING {
	return append([]byte{byte(Time)}, t.Bytes()...)
}

// Bytes returns gotype TIME as []byte by first
// converting to int64 of nanoseconds and then to []byte
func (t TIME) Bytes() []byte {
	return INT(t.Int()).Bytes()
}

// String returns gotype TIME as string
func (t TIME) String() string {
	tt := TIME{}
	if tt == t {
		return ""
	}
	return t.Native().Format(ISO8601)
}

// STRING returns gotype TIME as a gotype STRING
func (t TIME) STRING() STRING {
	return STRING(t.String())
}

// Serialize returns gotype TIME as serialized string
func (t TIME) Serialize() string {
	return `"` + t.String() + `"`
}

// Bool returns gotype TIME as bool
// false if empty, true if a Time
func (t TIME) Bool() bool {
	return t.Float64() != 0
}

// Int returns gotype TIME as int in Unix Time
func (t TIME) Int() int {
	return int(t.Native().UnixNano())
}

// Uint returns gotype TIME as uint in Unix Time
func (t TIME) Uint() uint {
	return uint(t.Native().UnixNano())
}

// Float returns gotype TIME as float64 in Unix Time
func (t TIME) Float64() float64 {
	return float64(t.Native().UnixNano())
}

// Time returns gotype TIME as time.Time
func (t TIME) Time() time.Time {
	return time.Time(t)
}

// ------------------------------------------------------------ /
// GOLANG STANDARD IMPLEMENTATIONS
// implementations of functions natively available for
// time in golang
// referenced packages: time
// ------------------------------------------------------------ /

// Weekday specifies a day of the week (Sunday = 0, ...).
type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// Month specifies a month of the year (January = 1, ...).
type Month int

const (
	January Month = iota + 1
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func Now() TIME {
	return TIME(time.Now().UTC())
}

func NewTime(year int, month int, day int, hour int, min int, sec int, nsec int, loc *time.Location) TIME {
	return TIME(time.Date(int(year), time.Month(int(month)), int(day), int(hour), int(min), int(sec), int(nsec), loc).UTC())
}

func NewDate(year int, month int, day int) TIME {
	return TIME(time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.UTC))
}

func ParseTime(f string, v string) TIME {
	t, err := time.Parse(string(f), string(v))
	if err != nil {
		panic("could not parse format time string")
	}
	return TIME(t)
}

func (t TIME) IsDate() bool {
	return t == t.Date()
}

func (t TIME) Unix() int {
	return int(t.Native().Unix())
}

func (t TIME) UnixNano() int {
	return int(t.Native().UnixNano())
}

func (t TIME) Equal(u TIME) bool {
	return t.Native().Equal(u.Native())
}

func (t TIME) Before(u TIME) bool {
	return t.Native().Before(u.Native())
}

func (t TIME) After(u TIME) bool {
	return t.Native().After(u.Native())
}

func (t TIME) AddDate(years int, months int, days int) TIME {
	return TIME(t.Native().AddDate(int(years), int(months), int(days)))
}

func (t TIME) Add(d Duration) TIME {
	return TIME(t.Native().Add(d.Native()))
}

func (t TIME) Sub(u TIME) Duration {
	return Duration(t.Native().Sub(u.Native()))
}

func (t TIME) Round(d Duration) TIME {
	return TIME(t.Native().Round(d.Native()))
}

func (t TIME) Format(f string) string {
	return t.Native().Format(string(f))
}

func (t TIME) Date() TIME {
	y, m, d := t.Native().Date()
	return NewDate(y, int(m), d)
}

func (t TIME) Year() int {
	return int(t.Native().Year())
}

func (t TIME) Month() int {
	return int(t.Native().Month())
}

func (t TIME) Day() int {
	return t.Native().Day()
}

func (t TIME) Weekday() Weekday {
	return Weekday(t.Native().Weekday())
}

func (t TIME) Hour() int {
	return t.Native().Hour()
}

func (t TIME) Minute() int {
	return t.Native().Minute()
}

func (t TIME) Second() int {
	return t.Native().Day()
}

func (t TIME) Nanosecond() int {
	return t.Native().Nanosecond()
}

func (t TIME) Location() *time.Location {
	return t.Native().Location()
}

// Duration represents the elapsed time between two Times
// as an int64 nanosecond count. The representation limits
// the largest representable duration to approximately 290 years.
type Duration time.Duration

const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

func Since(t TIME) Duration {
	return Duration(time.Since(t.Native()))
}

func (d Duration) Native() time.Duration {
	return time.Duration(d)
}
func (d Duration) String() string {
	return d.Native().String()
}

func (d Duration) STRING() STRING {
	return STRING(d.Native().String())
}

func (d Duration) Hours() int {
	return FLOAT(d.Native().Hours()).Int()
}

func (d Duration) Minutes() int {
	return FLOAT(d.Native().Minutes()).Int()
}

func (d Duration) Seconds() int {
	return FLOAT(d.Native().Hours()).Int()
}

func (d Duration) Milliseconds() int {
	return FLOAT(d.Native().Milliseconds()).Int()
}

func (d Duration) Nanoseconds() int {
	return FLOAT(d.Native().Nanoseconds()).Int()
}

func (d Duration) Round(m Duration) Duration {
	return Duration(d.Native().Round(m.Native()))
}

func (d Duration) Truncate(m Duration) Duration {
	return Duration(d.Native().Truncate(m.Native()))
}

// Location maps time instants to the zone in use at that time.
// Typically, the Location represents the collection of time
// offsets in use in a geographical area. For many Locations the
// time offset varies depending on whether daylight savings time
// is in use at the time instant.
type Location time.Location

func (l Location) Native() time.Location {
	return time.Location(l)
}

func (l Location) String() string {
	loc := l.Native()
	p := &loc
	return p.String()
}

func (l Location) STRING() STRING {
	return STRING(l.String())
}

// ------------------------------------------------------------ /
// GOTYPE EXPANDED FUNCTIONS
// implementations of new functions for
// time in gotype
// referenced packages:
// ------------------------------------------------------------ /

// DaysSince returns the number of Days since lt until time t
func (t TIME) DaysSince(lt TIME) int {
	return int(t.Native().Sub(lt.Native()).Hours() / 24)
}

// MonthsSince returns the number of full Months since lt until time t
func (t TIME) MonthsSince(lt TIME) int {
	pm := int(float64(lt.Day()) / math.Min(float64(t.Day()), float64(lt.DaysInMonth())))
	return (t.Year() - lt.Year()) + (t.Month() - lt.Month()) - 1 + pm
}

// YearsSince returns the number of full Years since lt until time t
func (t TIME) YearsSince(lt TIME) int {
	return int(float64(t.MonthsSince(lt) / 12))
}

// DaysInMonth returns the number of calendar days in month of TIME 't'
func (t TIME) DaysInMonth() int {
	return NewDate(t.Year(), t.Month(), 0).AddDate(0, 1, 0).Day()
}

// MonthStart returns the first date of the month for time 't'
func (t TIME) MonthStart() TIME {
	return NewDate(t.Year(), t.Month(), 1)
}

// MonthEnd returns the last nanosecond of the month for time 't'
func (t TIME) MonthEnd() TIME {
	return t.MonthStart().AddDate(0, 1, 0).Add(-1 * Nanosecond)
}

// QuarterStart returns the first date of the quarter
// for time 't' with year ending in month 'ye'
func (t TIME) QuarterStart(ye int) TIME {
	ye = ye % 3
	ye = (3-((t.Month()-ye)%3))%3 - 2
	return t.AddDate(0, ye, 0).MonthStart()
}

// QuarterEnd returns the last nanosecond of the quarter
// for time 't' with year ending in month 'ye'
func (t TIME) QuarterEnd(ye int) TIME {
	ye = ye % 3
	ye = (3 - ((t.Month() - ye) % 3)) % 3
	return t.AddDate(0, ye, 0).MonthEnd()
}

// YearStart returns the first date of the year
// for time 't' with year ending in month 'ye'
func (t TIME) YearStart(ye int) TIME {
	var y int
	if t.Month() < ye+1 {
		y = 1
	}
	return NewDate(t.Year()-y, ye+1, 1)
}

// YearEnd returns the last nanosecond of the year
// for time 't' with year ending in month 'ye'
func (t TIME) YearEnd(ye int) TIME {
	return t.YearStart(ye).AddDate(0, 12, 0).Add(-1 * Nanosecond)
}

func (t TIME) IsHoliday() (Holiday, bool) {
	y, m, d := t.Year(), t.Month(), t.Day()
	for _, i := range GetUsHolidays().List {
		h := i.Date(y)
		if m == h.Month() && d == h.Day() {
			return i, true
		}
	}
	return Holiday{}, false
}

// HOLIDAYS
// methods and storage for standard and custom holidays

type Holiday struct {
	time.Time
	// Name is the common name of the holiday
	Name string
	// Date returns the date of the holiday for year 'y'
	Date func(y int) TIME
}

type Holidays struct {
	List []Holiday
}

func GetUsHolidays() Holidays {
	h := Holidays{}
	h.List = append(h.List, Holiday{Name: "New Years Day", Date: NewYears})
	h.List = append(h.List, Holiday{Name: "Martin Luther King Day", Date: MlkDay})
	h.List = append(h.List, Holiday{Name: "Inauguration Day", Date: InagurationDay})
	h.List = append(h.List, Holiday{Name: "Presidents Day", Date: PresidentsDay})
	h.List = append(h.List, Holiday{Name: "Memorial Day", Date: MemorialDay})
	h.List = append(h.List, Holiday{Name: "Juneteenth", Date: NationalIndependenceDay})
	h.List = append(h.List, Holiday{Name: "Independence Day", Date: IndependenceDay})
	h.List = append(h.List, Holiday{Name: "Labor Day", Date: LaborDay})
	h.List = append(h.List, Holiday{Name: "Columbus Day", Date: ColumbusDay})
	h.List = append(h.List, Holiday{Name: "Veterans Day", Date: VeteransDay})
	h.List = append(h.List, Holiday{Name: "Thanksgiving", Date: Thanksgiving})
	h.List = append(h.List, Holiday{Name: "Christmas", Date: Christmas})
	return h
}

func (h *Holidays) IsHoliday(t TIME) bool {
	y, m, d := t.Year(), t.Month(), t.Day()
	for _, i := range h.List {
		h := i.Date(y)
		if m == h.Month() && d == h.Day() {
			return true
		}
	}
	return false
}

// Instance returns the date of the 'i' instance of weekday 'wd'
// in month 'm' of year 'y'; if i < 0 returns the last instance, and
// panics if 'i' is 0 or exceeds the number of instances
func Instance(i int, wd Weekday, m Month, y int) TIME {
	s := NewDate(y, int(m), 1)
	e := s.MonthEnd().Add(-24*Hour + Nanosecond)
	f := s.Weekday()
	l := e.Weekday()
	o := 0
	if i < 0 {
		if wd > l {
			o = 7
		}
		return e.AddDate(0, 0, int(wd-l)-o)
	}
	if wd >= f {
		o = 7
	}
	r := s.AddDate(0, 0, i*7+int(wd-f)-o)
	if r.After(e) || r.Before(s) {
		panic("instance must be greater than 0 and not exceed instances in month")
	}
	return r
}

// NewYears returns the observed date for new years day of for year 'y'
func NewYears(y int) TIME {
	return HolidayObserved(NewDate(y, int(January), 1))
}

// MlkDay returns the date of Martin Luther King Jr Day for year 'y'
func MlkDay(y int) TIME {
	return Instance(3, Monday, January, y)
}

// InagurationDay returns the date of the presidential inaguration for year 'y'
func InagurationDay(y int) TIME {
	y -= (y - 1) % 4
	return HolidayObserved(NewDate(y, int(January), 20))
}

// PresidentsDay returns the date of President's Day (or Washington's Birthday) for year 'y'
func PresidentsDay(y int) TIME {
	return Instance(3, Monday, February, y)
}

// GoodFriday returns the date of good friday for year 'y'
func GoodFriday(y int) TIME {
	return Easter(y).AddDate(0, 0, -2)
}

// Easter returns the date of easter for year 'y'
func Easter(y int) TIME {
	var yr, c, n, k, i, j, l, m, d float64
	yr = float64(y)
	c = math.Floor(yr / 100)
	n = yr - 19*math.Floor(yr/19)
	k = math.Floor((c - 17) / 25)
	i = c - math.Floor(c/4) - math.Floor((c-k)/3) + 19*n + 15
	i = i - 30*math.Floor(i/30)
	i = i - math.Floor(i/28)*(1-math.Floor(i/28)*math.Floor(29/(i+1))*math.Floor((21-n)/11))
	j = yr + math.Floor(yr/4) + i + 2 - c + math.Floor(c/4)
	j = j - 7*math.Floor(j/7)
	l = i - j
	m = 3 + math.Floor((l+40)/44)
	d = l + 28 - 31*math.Floor(m/4)
	return NewDate(y, int(m), int(d))
}

// MemorialDay returns the date of Memorial Day for year 'y'
func MemorialDay(y int) TIME {
	return Instance(-1, Monday, May, y)
}

// NationalIndependenceDay returns the observed date for
// Junteenth National Independence Day for year 'y'
func NationalIndependenceDay(y int) TIME {
	return HolidayObserved(NewDate(y, int(June), 19))
}

// IndependenceDay returns the observed date for US Independence Day for year 'y'
func IndependenceDay(y int) TIME {
	return HolidayObserved(NewDate(y, int(July), 4))
}

// LaborDay returns the date of Labor Day for year 'y'
func LaborDay(y int) TIME {
	return Instance(1, Monday, September, y)
}

// ColumbusDay returns the date of Columbus Day for year 'y'
func ColumbusDay(y int) TIME {
	return Instance(2, Monday, October, y)
}

// VeteransDay returns the observed date for Veterans Day for year 'y'
func VeteransDay(y int) TIME {
	return HolidayObserved(NewDate(y, int(November), 11))
}

// Thanksgiving returns the date of Thanksgiving Day for year 'y'
func Thanksgiving(y int) TIME {
	return Instance(4, Thursday, November, y)
}

// Christmas returns the observed date for Christmas Day for year 'y'
func Christmas(y int) TIME {
	return HolidayObserved(NewDate(y, int(December), 25))
}

// HolidayObserved returns the date holiday 'h' is observed,
// Friday if on Saturday and Monday if on Sunday
func HolidayObserved(h TIME) TIME {
	if h.Weekday() == Saturday {
		h = h.AddDate(0, 0, -1)
	} else if h.Weekday() == Sunday {
		h = h.AddDate(0, 0, 1)
	}
	return h
}
