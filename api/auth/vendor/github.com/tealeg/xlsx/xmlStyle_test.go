package xlsx

type XMLStyleSuite struct{}

var _ = Suite(&XMLStyleSuite{})

// Test we produce valid output for an empty style file.
func (x *XMLStyleSuite) TestMarshalEmptyXlsxStyleSheet(c *C) {
	styles := newXlsxStyleSheet(nil)
	result, err := Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, `<?xml version="1.0" encoding="UTF-8"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"></styleSheet>`)
}

// Test we produce valid output for a style file with one font definition.
func (x *XMLStyleSuite) TestMarshalXlsxStyleSheetWithAFont(c *C) {
	styles := newXlsxStyleSheet(nil)
	Fonts = xlsxFonts{}
	Count = 1
	Font = make([]xlsxFont, 1)
	font := xlsxFont{}
	Val = "10"
	Val = "Andale Mono"
	B = &xlsxVal{}
	I = &xlsxVal{}
	U = &xlsxVal{}
	Font[0] = font

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><fonts count="1"><font><sz val="10"/><name val="Andale Mono"/><b/><i/><u/></font></fonts></styleSheet>`
	result, err := Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, expected)
}

// Test we produce valid output for a style file with one fill definition.
func (x *XMLStyleSuite) TestMarshalXlsxStyleSheetWithAFill(c *C) {
	styles := newXlsxStyleSheet(nil)
	Fills = xlsxFills{}
	Count = 1
	Fill = make([]xlsxFill, 1)
	fill := xlsxFill{}
	patternFill := xlsxPatternFill{
		PatternType: "solid",
		FgColor:     xlsxColor{RGB: "#FFFFFF"},
		BgColor:     xlsxColor{RGB: "#000000"}}
	PatternFill = patternFill
	Fill[0] = fill

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><fills count="1"><fill><patternFill patternType="solid"><fgColor rgb="#FFFFFF"/><bgColor rgb="#000000"/></patternFill></fill></fills></styleSheet>`
	result, err := Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, expected)
}

// Test we produce valid output for a style file with one border definition.
// Empty elements are required to accommodate for Excel quirks.
func (x *XMLStyleSuite) TestMarshalXlsxStyleSheetWithABorder(c *C) {
	styles := newXlsxStyleSheet(nil)
	Borders = xlsxBorders{}
	Count = 1
	Border = make([]xlsxBorder, 1)
	border := xlsxBorder{}
	Style = "solid"
	Style = ""
	Border[0] = border
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><borders count="1"><border><left style="solid"></left><right style=""></right><top style=""></top><bottom style=""></bottom></border></borders></styleSheet>`

	result, err := Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, expected)
}

// Test we produce valid output for a style file with one cellStyleXf definition.
func (x *XMLStyleSuite) TestMarshalXlsxStyleSheetWithACellStyleXf(c *C) {
	styles := newXlsxStyleSheet(nil)
	CellStyleXfs = &xlsxCellStyleXfs{}
	Count = 1
	Xf = make([]xlsxXf, 1)
	xf := xlsxXf{}
	ApplyAlignment = true
	ApplyBorder = true
	ApplyFont = true
	ApplyFill = true
	ApplyProtection = true
	BorderId = 0
	FillId = 0
	FontId = 0
	NumFmtId = 0
	Alignment = xlsxAlignment{
		Horizontal:   "left",
		Indent:       1,
		ShrinkToFit:  true,
		TextRotation: 0,
		Vertical:     "middle",
		WrapText:     false}
	Xf[0] = xf

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><cellStyleXfs count="1"><xf applyAlignment="1" applyBorder="1" applyFont="1" applyFill="1" applyNumberFormat="0" applyProtection="1" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="left" indent="1" shrinkToFit="1" textRotation="0" vertical="middle" wrapText="0"/></xf></cellStyleXfs></styleSheet>`
	result, err := Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, expected)
}

// Test we produce valid output for a style file with one cellStyle definition.
func (x *XMLStyleSuite) TestMarshalXlsxStyleSheetWithACellStyle(c *C) {
	var builtInId int
	styles := newXlsxStyleSheet(nil)
	CellStyles = &xlsxCellStyles{Count: 1}
	CellStyle = make([]xlsxCellStyle, 1)

	builtInId = 31
	CellStyle[0] = xlsxCellStyle{
		Name:      "Bob",
		BuiltInId: &builtInId, // XXX Todo - work out built-ins!
		XfId:      0,
	}
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><cellStyles count="1"><cellStyle builtInId="31" name="Bob" xfId="0"></cellStyle></cellStyles></styleSheet>`
	result, err := Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, expected)
}

// Test we produce valid output for a style file with one cellXf
// definition.
func (x *XMLStyleSuite) TestMarshalXlsxStyleSheetWithACellXf(c *C) {
	styles := newXlsxStyleSheet(nil)
	CellXfs = xlsxCellXfs{}
	Count = 1
	Xf = make([]xlsxXf, 1)
	xf := xlsxXf{}
	ApplyAlignment = true
	ApplyBorder = true
	ApplyFont = true
	ApplyFill = true
	ApplyNumberFormat = true
	ApplyProtection = true
	BorderId = 0
	FillId = 0
	FontId = 0
	NumFmtId = 0
	Alignment = xlsxAlignment{
		Horizontal:   "left",
		Indent:       1,
		ShrinkToFit:  true,
		TextRotation: 0,
		Vertical:     "middle",
		WrapText:     false}
	Xf[0] = xf

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><cellXfs count="1"><xf applyAlignment="1" applyBorder="1" applyFont="1" applyFill="1" applyNumberFormat="1" applyProtection="1" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="left" indent="1" shrinkToFit="1" textRotation="0" vertical="middle" wrapText="0"/></xf></cellXfs></styleSheet>`
	result, err := Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, expected)
}

// Test we produce valid output for a style file with one NumFmt
// definition.
func (x *XMLStyleSuite) TestMarshalXlsxStyleSheetWithANumFmt(c *C) {
	styles := &xlsxStyleSheet{}
	NumFmts = xlsxNumFmts{}
	NumFmt = make([]xlsxNumFmt, 0)
	numFmt := xlsxNumFmt{NumFmtId: 164, FormatCode: "GENERAL"}
	addNumFmt(numFmt)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><numFmts count="1"><numFmt numFmtId="164" formatCode="GENERAL"/></numFmts></styleSheet>`
	result, err := Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, expected)
}

func (x *XMLStyleSuite) TestFontEquals(c *C) {
	fontA := xlsxFont{Sz: xlsxVal{Val: "11"},
		Color:  xlsxColor{RGB: "FFFF0000"},
		Name:   xlsxVal{Val: "Calibri"},
		Family: xlsxVal{Val: "2"},
		B:      &xlsxVal{},
		I:      &xlsxVal{},
		U:      &xlsxVal{}}
	fontB := xlsxFont{Sz: xlsxVal{Val: "11"},
		Color:  xlsxColor{RGB: "FFFF0000"},
		Name:   xlsxVal{Val: "Calibri"},
		Family: xlsxVal{Val: "2"},
		B:      &xlsxVal{},
		I:      &xlsxVal{},
		U:      &xlsxVal{}}

	c.Assert(Equals(fontB), Equals, true)
	Val = "12"
	c.Assert(Equals(fontB), Equals, false)
	Val = "11"
	RGB = "12345678"
	c.Assert(Equals(fontB), Equals, false)
	RGB = "FFFF0000"
	Val = "Arial"
	c.Assert(Equals(fontB), Equals, false)
	Val = "Calibri"
	Val = "1"
	c.Assert(Equals(fontB), Equals, false)
	Val = "2"
	B = nil
	c.Assert(Equals(fontB), Equals, false)
	B = &xlsxVal{}
	I = nil
	c.Assert(Equals(fontB), Equals, false)
	I = &xlsxVal{}
	U = nil
	c.Assert(Equals(fontB), Equals, false)
	U = &xlsxVal{}
	// For sanity
	c.Assert(Equals(fontB), Equals, true)
}

func (x *XMLStyleSuite) TestFillEquals(c *C) {
	fillA := xlsxFill{PatternFill: xlsxPatternFill{
		PatternType: "solid",
		FgColor:     xlsxColor{RGB: "FFFF0000"},
		BgColor:     xlsxColor{RGB: "0000FFFF"}}}
	fillB := xlsxFill{PatternFill: xlsxPatternFill{
		PatternType: "solid",
		FgColor:     xlsxColor{RGB: "FFFF0000"},
		BgColor:     xlsxColor{RGB: "0000FFFF"}}}
	c.Assert(Equals(fillB), Equals, true)
	PatternType = "gray125"
	c.Assert(Equals(fillB), Equals, false)
	PatternType = "solid"
	RGB = "00FF00FF"
	c.Assert(Equals(fillB), Equals, false)
	RGB = "FFFF0000"
	RGB = "12456789"
	c.Assert(Equals(fillB), Equals, false)
	RGB = "0000FFFF"
	// For sanity
	c.Assert(Equals(fillB), Equals, true)
}

func (x *XMLStyleSuite) TestBorderEquals(c *C) {
	borderA := xlsxBorder{Left: xlsxLine{Style: "none"},
		Right:  xlsxLine{Style: "none"},
		Top:    xlsxLine{Style: "none"},
		Bottom: xlsxLine{Style: "none"}}
	borderB := xlsxBorder{Left: xlsxLine{Style: "none"},
		Right:  xlsxLine{Style: "none"},
		Top:    xlsxLine{Style: "none"},
		Bottom: xlsxLine{Style: "none"}}
	c.Assert(Equals(borderB), Equals, true)
	Style = "thin"
	c.Assert(Equals(borderB), Equals, false)
	Style = "none"
	Style = "thin"
	c.Assert(Equals(borderB), Equals, false)
	Style = "none"
	Style = "thin"
	c.Assert(Equals(borderB), Equals, false)
	Style = "none"
	Style = "thin"
	c.Assert(Equals(borderB), Equals, false)
	Style = "none"
	// for sanity
	c.Assert(Equals(borderB), Equals, true)
}

func (x *XMLStyleSuite) TestXfEquals(c *C) {
	xfA := xlsxXf{
		ApplyAlignment:  true,
		ApplyBorder:     true,
		ApplyFont:       true,
		ApplyFill:       true,
		ApplyProtection: true,
		BorderId:        0,
		FillId:          0,
		FontId:          0,
		NumFmtId:        0}
	xfB := xlsxXf{
		ApplyAlignment:  true,
		ApplyBorder:     true,
		ApplyFont:       true,
		ApplyFill:       true,
		ApplyProtection: true,
		BorderId:        0,
		FillId:          0,
		FontId:          0,
		NumFmtId:        0}
	c.Assert(Equals(xfB), Equals, true)
	ApplyAlignment = false
	c.Assert(Equals(xfB), Equals, false)
	ApplyAlignment = true
	ApplyBorder = false
	c.Assert(Equals(xfB), Equals, false)
	ApplyBorder = true
	ApplyFont = false
	c.Assert(Equals(xfB), Equals, false)
	ApplyFont = true
	ApplyFill = false
	c.Assert(Equals(xfB), Equals, false)
	ApplyFill = true
	ApplyProtection = false
	c.Assert(Equals(xfB), Equals, false)
	ApplyProtection = true
	BorderId = 1
	c.Assert(Equals(xfB), Equals, false)
	BorderId = 0
	FillId = 1
	c.Assert(Equals(xfB), Equals, false)
	FillId = 0
	FontId = 1
	c.Assert(Equals(xfB), Equals, false)
	FontId = 0
	NumFmtId = 1
	c.Assert(Equals(xfB), Equals, false)
	NumFmtId = 0
	// for sanity
	c.Assert(Equals(xfB), Equals, true)

	var i1 int = 1

	XfId = &i1
	c.Assert(Equals(xfB), Equals, false)

	XfId = &i1
	c.Assert(Equals(xfB), Equals, true)

	var i2 int = 1
	XfId = &i2
	c.Assert(Equals(xfB), Equals, true)

	i2 = 2
	c.Assert(Equals(xfB), Equals, false)
}

func (s *CellSuite) TestNewNumFmt(c *C) {
	styles := newXlsxStyleSheet(nil)
	NumFmts = xlsxNumFmts{}
	NumFmt = make([]xlsxNumFmt, 0)

	c.Assert(newNumFmt("0"), DeepEquals, xlsxNumFmt{1, "0"})
	c.Assert(newNumFmt("0.00e+00"), DeepEquals, xlsxNumFmt{11, "0.00e+00"})
	c.Assert(newNumFmt("mm-dd-yy"), DeepEquals, xlsxNumFmt{14, "mm-dd-yy"})
	c.Assert(newNumFmt("hh:mm:ss"), DeepEquals, xlsxNumFmt{164, "hh:mm:ss"})
	c.Assert(len(NumFmt), Equals, 1)
}

func (s *CellSuite) TestAddNumFmt(c *C) {
	styles := &xlsxStyleSheet{}
	NumFmts = xlsxNumFmts{}
	NumFmt = make([]xlsxNumFmt, 0)

	addNumFmt(xlsxNumFmt{1, "0"})
	c.Assert(Count, Equals, 0)
	addNumFmt(xlsxNumFmt{14, "mm-dd-yy"})
	c.Assert(Count, Equals, 0)
	addNumFmt(xlsxNumFmt{164, "hh:mm:ss"})
	c.Assert(Count, Equals, 1)
	addNumFmt(xlsxNumFmt{165, "yyyy/mm/dd"})
	c.Assert(Count, Equals, 2)
	addNumFmt(xlsxNumFmt{165, "yyyy/mm/dd"})
	c.Assert(Count, Equals, 2)
}
