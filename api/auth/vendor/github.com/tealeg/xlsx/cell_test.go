package xlsx

import (
	"math"
	"time"
)

type CellSuite struct{}

var _ = Suite(&CellSuite{})

// Test that we can set and get a Value from a Cell
func (s *CellSuite) TestValueSet(c *C) {
	// Note, this test is fairly pointless, it serves mostly to
	// reinforce that this functionality is important, and should
	// the mechanics of this all change at some point, to remind
	// us not to lose this.
	cell := Cell{}
	Value = "A string"
	c.Assert(Value, Equals, "A string")
}

// Test that GetStyle correctly converts the xlsxStyle.Fonts.
func (s *CellSuite) TestGetStyleWithFonts(c *C) {
	font := NewFont(10, "Calibra")
	style := NewStyle()
	Font = *font

	cell := &Cell{Value: "123", style: style}
	style = GetStyle()
	c.Assert(style, NotNil)
	c.Assert(Size, Equals, 10)
	c.Assert(Name, Equals, "Calibra")
}

// Test that SetStyle correctly translates into a xlsxFont element
func (s *CellSuite) TestSetStyleWithFonts(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Test")
	row := AddRow()
	cell := AddCell()
	font := NewFont(12, "Calibra")
	style := NewStyle()
	Font = *font
	SetStyle(style)
	style = GetStyle()
	xFont, _, _, _ := makeXLSXStyleElements()
	c.Assert(Val, Equals, "12")
	c.Assert(Val, Equals, "Calibra")
}

// Test that GetStyle correctly converts the xlsxStyle.Fills.
func (s *CellSuite) TestGetStyleWithFills(c *C) {
	fill := *NewFill("solid", "FF000000", "00FF0000")
	style := NewStyle()
	Fill = fill
	cell := &Cell{Value: "123", style: style}
	style = GetStyle()
	_, xFill, _, _ := makeXLSXStyleElements()
	c.Assert(PatternType, Equals, "solid")
	c.Assert(RGB, Equals, "00FF0000")
	c.Assert(RGB, Equals, "FF000000")
}

// Test that SetStyle correctly updates xlsxStyle.Fills.
func (s *CellSuite) TestSetStyleWithFills(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Test")
	row := AddRow()
	cell := AddCell()
	fill := NewFill("solid", "00FF0000", "FF000000")
	style := NewStyle()
	Fill = *fill
	SetStyle(style)
	style = GetStyle()
	_, xFill, _, _ := makeXLSXStyleElements()
	xPatternFill := PatternFill
	c.Assert(PatternType, Equals, "solid")
	c.Assert(RGB, Equals, "00FF0000")
	c.Assert(RGB, Equals, "FF000000")
}

// Test that GetStyle correctly converts the xlsxStyle.Borders.
func (s *CellSuite) TestGetStyleWithBorders(c *C) {
	border := *NewBorder("thin", "thin", "thin", "thin")
	style := NewStyle()
	Border = border
	cell := Cell{Value: "123", style: style}
	style = GetStyle()
	_, _, xBorder, _ := makeXLSXStyleElements()
	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "thin")
}

// We can return a string representation of the formatted data
func (l *CellSuite) TestSetFloatWithFormat(c *C) {
	cell := Cell{}
	SetFloatWithFormat(37947.75334343, "yyyy/mm/dd")
	c.Assert(Value, Equals, "37947.75334343")
	c.Assert(NumFmt, Equals, "yyyy/mm/dd")
	c.Assert(Type(), Equals, CellTypeNumeric)
}

func (l *CellSuite) TestSetFloat(c *C) {
	cell := Cell{}
	SetFloat(0)
	c.Assert(Value, Equals, "0")
	SetFloat(0.000005)
	c.Assert(Value, Equals, "5e-06")
	SetFloat(100.0)
	c.Assert(Value, Equals, "100")
	SetFloat(37947.75334343)
	c.Assert(Value, Equals, "37947.75334343")
}

func (s *CellSuite) TestGetTime(c *C) {
	cell := Cell{}
	SetFloat(0)
	date, err := GetTime(false)
	c.Assert(err, Equals, nil)
	c.Assert(date, Equals, time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC))
	SetFloat(39813.0)
	date, err = GetTime(true)
	c.Assert(err, Equals, nil)
	c.Assert(date, Equals, time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC))
	Value = "d"
	_, err = GetTime(false)
	c.Assert(err, NotNil)
}

// FormattedValue returns an error for formatting errors
func (l *CellSuite) TestFormattedValueErrorsOnBadFormat(c *C) {
	cell := Cell{Value: "Fudge Cake"}
	NumFmt = "#,##0 ;(#,##0)"
	value, err := FormattedValue()
	c.Assert(value, Equals, "Fudge Cake")
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "strconv.ParseFloat: parsing \"Fudge Cake\": invalid syntax")
}

// FormattedValue returns a string containing error text for formatting errors
func (l *CellSuite) TestFormattedValueReturnsErrorAsValueForBadFormat(c *C) {
	cell := Cell{Value: "Fudge Cake"}
	NumFmt = "#,##0 ;(#,##0)"
	_, err := FormattedValue()
	c.Assert(err.Error(), Equals, "strconv.ParseFloat: parsing \"Fudge Cake\": invalid syntax")
}

// formattedValueChecker removes all the boilerplate for testing Cell.FormattedValue
// after its change from returning one value (a string) to two values (string, error)
// This allows all the old one-line asserts in the test to continue to be one
// line, instead of multi-line with error checking.
type formattedValueChecker struct {
	c *C
}

func (fvc *formattedValueChecker) Equals(cell Cell, expected string) {
	val, err := FormattedValue()
	if err != nil {
		fvc.c.Error(err)
	}
	fvc.c.Assert(val, Equals, expected)
}

// We can return a string representation of the formatted data
func (l *CellSuite) TestFormattedValue(c *C) {
	// XXX TODO, this test should probably be split down, and made
	// in terms of SafeFormattedValue, as FormattedValue wraps
	// that function now.
	cell := Cell{Value: "37947.7500001"}
	negativeCell := Cell{Value: "-37947.7500001"}
	smallCell := Cell{Value: "0.007"}
	earlyCell := Cell{Value: "2.1"}

	fvc := formattedValueChecker{c: c}

	NumFmt = "general"
	fvc.Equals(cell, "37947.7500001")
	NumFmt = "general"
	fvc.Equals(negativeCell, "-37947.7500001")

	// TODO: This test is currently broken.  For a string type cell, I
	// don't think FormattedValue() should be doing a numeric conversion on the value
	// before returning the string.
	NumFmt = "0"
	fvc.Equals(cell, "37947")

	NumFmt = "#,##0" // For the time being we're not doing
	// this comma formatting, so it'll fall back to the related
	// non-comma form.
	fvc.Equals(cell, "37947")

	NumFmt = "#,##0.00;(#,##0.00)"
	fvc.Equals(cell, "37947.75")

	NumFmt = "0.00"
	fvc.Equals(cell, "37947.75")

	NumFmt = "#,##0.00" // For the time being we're not doing
	// this comma formatting, so it'll fall back to the related
	// non-comma form.
	fvc.Equals(cell, "37947.75")

	NumFmt = "#,##0 ;(#,##0)"
	fvc.Equals(cell, "37947")
	NumFmt = "#,##0 ;(#,##0)"
	fvc.Equals(negativeCell, "(37947)")

	NumFmt = "#,##0 ;[red](#,##0)"
	fvc.Equals(cell, "37947")
	NumFmt = "#,##0 ;[red](#,##0)"
	fvc.Equals(negativeCell, "(37947)")

	NumFmt = "#,##0.00;(#,##0.00)"
	fvc.Equals(negativeCell, "(-37947.75)")

	NumFmt = "0%"
	fvc.Equals(cell, "3794775%")

	NumFmt = "0.00%"
	fvc.Equals(cell, "3794775.00%")

	NumFmt = "0.00e+00"
	fvc.Equals(cell, "3.794775e+04")

	NumFmt = "##0.0e+0" // This is wrong, but we'll use it for now.
	fvc.Equals(cell, "3.794775e+04")

	NumFmt = "mm-dd-yy"
	fvc.Equals(cell, "11-22-03")

	NumFmt = "d-mmm-yy"
	fvc.Equals(cell, "22-Nov-03")
	NumFmt = "d-mmm-yy"
	fvc.Equals(earlyCell, "1-Jan-00")

	NumFmt = "d-mmm"
	fvc.Equals(cell, "22-Nov")
	NumFmt = "d-mmm"
	fvc.Equals(earlyCell, "1-Jan")

	NumFmt = "mmm-yy"
	fvc.Equals(cell, "Nov-03")

	NumFmt = "h:mm am/pm"
	fvc.Equals(cell, "6:00 pm")
	NumFmt = "h:mm am/pm"
	fvc.Equals(smallCell, "12:10 am")

	NumFmt = "h:mm:ss am/pm"
	fvc.Equals(cell, "6:00:00 pm")
	NumFmt = "hh:mm:ss"
	fvc.Equals(cell, "18:00:00")
	NumFmt = "h:mm:ss am/pm"
	fvc.Equals(smallCell, "12:10:04 am")

	NumFmt = "h:mm"
	fvc.Equals(cell, "6:00")
	NumFmt = "h:mm"
	fvc.Equals(smallCell, "12:10")
	NumFmt = "hh:mm"
	fvc.Equals(smallCell, "00:10")

	NumFmt = "h:mm:ss"
	fvc.Equals(cell, "6:00:00")
	NumFmt = "hh:mm:ss"
	fvc.Equals(cell, "18:00:00")

	NumFmt = "hh:mm:ss"
	fvc.Equals(smallCell, "00:10:04")
	NumFmt = "h:mm:ss"
	fvc.Equals(smallCell, "12:10:04")

	NumFmt = "m/d/yy h:mm"
	fvc.Equals(cell, "11/22/03 6:00")
	NumFmt = "m/d/yy hh:mm"
	fvc.Equals(cell, "11/22/03 18:00")
	NumFmt = "m/d/yy h:mm"
	fvc.Equals(smallCell, "12/30/99 12:10")
	NumFmt = "m/d/yy hh:mm"
	fvc.Equals(smallCell, "12/30/99 00:10")
	NumFmt = "m/d/yy hh:mm"
	fvc.Equals(earlyCell, "1/1/00 02:24")
	NumFmt = "m/d/yy h:mm"
	fvc.Equals(earlyCell, "1/1/00 2:24")

	NumFmt = "mm:ss"
	fvc.Equals(cell, "00:00")
	NumFmt = "mm:ss"
	fvc.Equals(smallCell, "10:04")

	NumFmt = "[hh]:mm:ss"
	fvc.Equals(cell, "18:00:00")
	NumFmt = "[h]:mm:ss"
	fvc.Equals(cell, "6:00:00")
	NumFmt = "[h]:mm:ss"
	fvc.Equals(smallCell, "10:04")

	const (
		expect1 = "0000.0086"
		expect2 = "1004.8000"
		format  = "mmss.0000"
		tlen    = len(format)
	)

	for i := 0; i < 3; i++ {
		tfmt := format[0 : tlen-i]
		NumFmt = tfmt
		fvc.Equals(cell, expect1[0:tlen-i])
		NumFmt = tfmt
		fvc.Equals(smallCell, expect2[0:tlen-i])
	}

	NumFmt = "yyyy\\-mm\\-dd"
	fvc.Equals(cell, "2003\\-11\\-22")

	NumFmt = "dd/mm/yyyy hh:mm:ss"
	fvc.Equals(cell, "22/11/2003 18:00:00")

	NumFmt = "dd/mm/yy"
	fvc.Equals(cell, "22/11/03")
	NumFmt = "dd/mm/yy"
	fvc.Equals(earlyCell, "01/01/00")

	NumFmt = "hh:mm:ss"
	fvc.Equals(cell, "18:00:00")
	NumFmt = "hh:mm:ss"
	fvc.Equals(smallCell, "00:10:04")

	NumFmt = "dd/mm/yy\\ hh:mm"
	fvc.Equals(cell, "22/11/03\\ 18:00")

	NumFmt = "yyyy/mm/dd"
	fvc.Equals(cell, "2003/11/22")

	NumFmt = "yy-mm-dd"
	fvc.Equals(cell, "03-11-22")

	NumFmt = "d-mmm-yyyy"
	fvc.Equals(cell, "22-Nov-2003")
	NumFmt = "d-mmm-yyyy"
	fvc.Equals(earlyCell, "1-Jan-1900")

	NumFmt = "m/d/yy"
	fvc.Equals(cell, "11/22/03")
	NumFmt = "m/d/yy"
	fvc.Equals(earlyCell, "1/1/00")

	NumFmt = "m/d/yyyy"
	fvc.Equals(cell, "11/22/2003")
	NumFmt = "m/d/yyyy"
	fvc.Equals(earlyCell, "1/1/1900")

	NumFmt = "dd-mmm-yyyy"
	fvc.Equals(cell, "22-Nov-2003")

	NumFmt = "dd/mm/yyyy"
	fvc.Equals(cell, "22/11/2003")

	NumFmt = "mm/dd/yy hh:mm am/pm"
	fvc.Equals(cell, "11/22/03 18:00 pm")
	NumFmt = "mm/dd/yy h:mm am/pm"
	fvc.Equals(cell, "11/22/03 6:00 pm")

	NumFmt = "mm/dd/yyyy hh:mm:ss"
	fvc.Equals(cell, "11/22/2003 18:00:00")
	NumFmt = "mm/dd/yyyy hh:mm:ss"
	fvc.Equals(smallCell, "12/30/1899 00:10:04")

	NumFmt = "yyyy-mm-dd hh:mm:ss"
	fvc.Equals(cell, "2003-11-22 18:00:00")
	NumFmt = "yyyy-mm-dd hh:mm:ss"
	fvc.Equals(smallCell, "1899-12-30 00:10:04")

	NumFmt = "mmmm d, yyyy"
	fvc.Equals(cell, "November 22, 2003")
	NumFmt = "mmmm d, yyyy"
	fvc.Equals(smallCell, "December 30, 1899")

	NumFmt = "dddd, mmmm dd, yyyy"
	fvc.Equals(cell, "Saturday, November 22, 2003")
	NumFmt = "dddd, mmmm dd, yyyy"
	fvc.Equals(smallCell, "Saturday, December 30, 1899")
}

// test setters and getters
func (s *CellSuite) TestSetterGetters(c *C) {
	cell := Cell{}

	SetString("hello world")
	if val, err := String(); err != nil {
		c.Error(err)
	} else {
		c.Assert(val, Equals, "hello world")
	}
	c.Assert(Type(), Equals, CellTypeString)

	SetInt(1024)
	intValue, _ := Int()
	c.Assert(intValue, Equals, 1024)
	c.Assert(NumFmt, Equals, builtInNumFmt[builtInNumFmtIndex_GENERAL])
	c.Assert(Type(), Equals, CellTypeGeneral)

	SetInt64(1024)
	int64Value, _ := Int64()
	c.Assert(int64Value, Equals, int64(1024))
	c.Assert(NumFmt, Equals, builtInNumFmt[builtInNumFmtIndex_GENERAL])
	c.Assert(Type(), Equals, CellTypeGeneral)

	SetFloat(1.024)
	float, _ := Float()
	intValue, _ = Int() // convert
	c.Assert(float, Equals, 1.024)
	c.Assert(intValue, Equals, 1)
	c.Assert(NumFmt, Equals, builtInNumFmt[builtInNumFmtIndex_GENERAL])
	c.Assert(Type(), Equals, CellTypeGeneral)

	SetFormula("10+20")
	c.Assert(Formula(), Equals, "10+20")
	c.Assert(Type(), Equals, CellTypeFormula)
}

// TestOddInput is a regression test for #101. When the number format
// was "@" (string), the input below caused a crash in strconv.ParseFloat.
// The solution was to check if cell.Value was both a CellTypeString and
// had a NumFmt of "general" or "@" and short-circuit FormattedValue() if so.
func (s *CellSuite) TestOddInput(c *C) {
	cell := Cell{}
	odd := `[1],[12,"DATE NOT NULL DEFAULT '0000-00-00'"]`
	Value = odd
	NumFmt = "@"
	if val, err := String(); err != nil {
		c.Error(err)
	} else {
		c.Assert(val, Equals, odd)
	}
}

// TestBool tests basic Bool getting and setting booleans.
func (s *CellSuite) TestBool(c *C) {
	cell := Cell{}
	SetBool(true)
	c.Assert(Value, Equals, "1")
	c.Assert(Bool(), Equals, true)
	SetBool(false)
	c.Assert(Value, Equals, "0")
	c.Assert(Bool(), Equals, false)
}

// TestStringBool tests calling Bool on a non CellTypeBool value.
func (s *CellSuite) TestStringBool(c *C) {
	cell := Cell{}
	SetInt(0)
	c.Assert(Bool(), Equals, false)
	SetInt(1)
	c.Assert(Bool(), Equals, true)
	SetString("")
	c.Assert(Bool(), Equals, false)
	SetString("0")
	c.Assert(Bool(), Equals, true)
}

// TestSetValue tests whether SetValue handle properly for different type values.
func (s *CellSuite) TestSetValue(c *C) {
	cell := Cell{}

	// int
	for _, i := range []interface{}{1, int8(1), int16(1), int32(1), int64(1)} {
		SetValue(i)
		val, err := Int64()
		c.Assert(err, IsNil)
		c.Assert(val, Equals, int64(1))
	}

	// float
	for _, i := range []interface{}{1.11, float32(1.11), float64(1.11)} {
		SetValue(i)
		val, err := Float()
		c.Assert(err, IsNil)
		c.Assert(val, Equals, 1.11)
	}

	// time
	SetValue(time.Unix(0, 0))
	val, err := Float()
	c.Assert(err, IsNil)
	c.Assert(math.Floor(val), Equals, 25569.0)

	// string and nil
	for _, i := range []interface{}{nil, "", []byte("")} {
		SetValue(i)
		c.Assert(Value, Equals, "")
	}

	// others
	SetValue([]string{"test"})
	c.Assert(Value, Equals, "[test]")
}
