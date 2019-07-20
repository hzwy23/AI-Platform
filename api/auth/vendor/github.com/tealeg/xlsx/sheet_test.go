package xlsx

import (
	"bytes"
	"encoding/xml"
)

type SheetSuite struct{}

var _ = Suite(&SheetSuite{})

// Test we can add a Row to a Sheet
func (s *SheetSuite) TestAddRow(c *C) {
	var f *File
	f = NewFile()
	sheet, _ := AddSheet("MySheet")
	row := AddRow()
	c.Assert(row, NotNil)
	c.Assert(len(Rows), Equals, 1)
}

func (s *SheetSuite) TestMakeXLSXSheetFromRows(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	row := AddRow()
	cell := AddCell()
	Value = "A cell!"
	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	xSheet := makeXLSXSheet(refTable, styles)
	c.Assert(Ref, Equals, "A1")
	c.Assert(Row, HasLen, 1)
	xRow := Row[0]
	c.Assert(R, Equals, 1)
	c.Assert(Spans, Equals, "")
	c.Assert(C, HasLen, 1)
	xC := C[0]
	c.Assert(R, Equals, "A1")
	c.Assert(S, Equals, 0)
	c.Assert(T, Equals, "s") // Shared string type
	c.Assert(V, Equals, "0") // reference to shared string
	xSST := makeXLSXSST()
	c.Assert(xSST.Count, Equals, 1)
	c.Assert(xSST.UniqueCount, Equals, 1)
	c.Assert(xSST.SI, HasLen, 1)
	xSI := xSST.SI[0]
	c.Assert(xSI.T, Equals, "A cell!")
}

// Test if the NumFmts assigned properly according the FormatCode in cell.
func (s *SheetSuite) TestMakeXLSXSheetWithNumFormats(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	row := AddRow()

	cell1 := AddCell()
	Value = "A cell!"
	NumFmt = "general"

	cell2 := AddCell()
	Value = "37947.7500001"
	NumFmt = "0"

	cell3 := AddCell()
	Value = "37947.7500001"
	NumFmt = "mm-dd-yy"

	cell4 := AddCell()
	Value = "37947.7500001"
	NumFmt = "hh:mm:ss"

	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	worksheet := makeXLSXSheet(refTable, styles)

	c.Assert(CellStyleXfs, IsNil)

	c.Assert(Count, Equals, 5)
	c.Assert(NumFmtId, Equals, 0)
	c.Assert(NumFmtId, Equals, 0)
	c.Assert(NumFmtId, Equals, 1)
	c.Assert(NumFmtId, Equals, 14)
	c.Assert(NumFmtId, Equals, 164)
	c.Assert(Count, Equals, 1)
	c.Assert(NumFmtId, Equals, 164)
	c.Assert(FormatCode, Equals, "hh:mm:ss")

	// Finally we check that the cell points to the right CellXf /
	// CellStyleXf.
	c.Assert(S, Equals, 1)
	c.Assert(S, Equals, 2)
}

// When we create the xlsxSheet we also populate the xlsxStyles struct
// with style information.
func (s *SheetSuite) TestMakeXLSXSheetAlsoPopulatesXLSXSTyles(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	row := AddRow()

	cell1 := AddCell()
	Value = "A cell!"
	style1 := NewStyle()
	Font = *NewFont(10, "Verdana")
	Fill = *NewFill("solid", "FFFFFFFF", "00000000")
	Border = *NewBorder("none", "thin", "none", "thin")
	SetStyle(style1)

	// We need a second style to check that Xfs are populated correctly.
	cell2 := AddCell()
	Value = "Another cell!"
	style2 := NewStyle()
	Font = *NewFont(10, "Verdana")
	Fill = *NewFill("solid", "FFFFFFFF", "00000000")
	Border = *NewBorder("none", "thin", "none", "thin")
	SetStyle(style2)

	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	worksheet := makeXLSXSheet(refTable, styles)

	c.Assert(Count, Equals, 2)
	c.Assert(Val, Equals, "12")
	c.Assert(Val, Equals, "Verdana")
	c.Assert(Val, Equals, "10")
	c.Assert(Val, Equals, "Verdana")

	c.Assert(Count, Equals, 3)
	c.Assert(PatternType, Equals, "none")
	c.Assert(RGB, Equals, "FFFFFFFF")
	c.Assert(RGB, Equals, "00000000")

	c.Assert(Count, Equals, 2)
	c.Assert(Style, Equals, "none")
	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "none")
	c.Assert(Style, Equals, "thin")

	c.Assert(CellStyleXfs, IsNil)

	c.Assert(Count, Equals, 2)
	c.Assert(FontId, Equals, 0)
	c.Assert(FillId, Equals, 0)
	c.Assert(BorderId, Equals, 0)

	// Finally we check that the cell points to the right CellXf /
	// CellStyleXf.
	c.Assert(S, Equals, 1)
	c.Assert(S, Equals, 1)
}

// If the column width is not customised, the xslxCol.CustomWidth field is set to 0.
func (s *SheetSuite) TestMakeXLSXSheetDefaultsCustomColWidth(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	row := AddRow()
	cell1 := AddCell()
	Value = "A cell!"

	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	worksheet := makeXLSXSheet(refTable, styles)
	c.Assert(CustomWidth, Equals, false)
}

// If the column width is customised, the xslxCol.CustomWidth field is set to 1.
func (s *SheetSuite) TestMakeXLSXSheetSetsCustomColWidth(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	row := AddRow()
	cell1 := AddCell()
	Value = "A cell!"
	err := SetColWidth(0, 1, 10.5)
	c.Assert(err, IsNil)

	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	worksheet := makeXLSXSheet(refTable, styles)
	c.Assert(CustomWidth, Equals, true)
}

func (s *SheetSuite) TestMarshalSheet(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	row := AddRow()
	cell := AddCell()
	Value = "A cell!"
	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	xSheet := makeXLSXSheet(refTable, styles)

	output := bytes.NewBufferString(xml.Header)
	body, err := xml.Marshal(xSheet)
	c.Assert(err, IsNil)
	c.Assert(body, NotNil)
	_, err = output.Write(body)
	c.Assert(err, IsNil)

	expectedXLSXSheet := `<?xml version="1.0" encoding="UTF-8"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetPr filterMode="false"><pageSetUpPr fitToPage="false"></pageSetUpPr></sheetPr><dimension ref="A1"></dimension><sheetViews><sheetView windowProtection="false" showFormulas="false" showGridLines="true" showRowColHeaders="true" showZeros="true" rightToLeft="false" tabSelected="true" showOutlineSymbols="true" defaultGridColor="true" view="normal" topLeftCell="A1" colorId="64" zoomScale="100" zoomScaleNormal="100" zoomScalePageLayoutView="100" workbookViewId="0"><selection pane="topLeft" activeCell="A1" activeCellId="0" sqref="A1"></selection></sheetView></sheetViews><sheetFormatPr defaultRowHeight="12.85"></sheetFormatPr><cols><col collapsed="false" hidden="false" max="1" min="1" style="0" width="9.5"></col></cols><sheetData><row r="1"><c r="A1" t="s"><v>0</v></c></row></sheetData><printOptions headings="false" gridLines="false" gridLinesSet="true" horizontalCentered="false" verticalCentered="false"></printOptions><pageMargins left="0.7875" right="0.7875" top="1.05277777777778" bottom="1.05277777777778" header="0.7875" footer="0.7875"></pageMargins><pageSetup paperSize="9" scale="100" firstPageNumber="1" fitToWidth="1" fitToHeight="1" pageOrder="downThenOver" orientation="portrait" usePrinterDefaults="false" blackAndWhite="false" draft="false" cellComments="none" useFirstPageNumber="true" horizontalDpi="300" verticalDpi="300" copies="1"></pageSetup><headerFooter differentFirst="false" differentOddEven="false"><oddHeader>&amp;C&amp;&#34;Times New Roman,Regular&#34;&amp;12&amp;A</oddHeader><oddFooter>&amp;C&amp;&#34;Times New Roman,Regular&#34;&amp;12Page &amp;P</oddFooter></headerFooter></worksheet>`

	c.Assert(output.String(), Equals, expectedXLSXSheet)
}

func (s *SheetSuite) TestMarshalSheetWithMultipleCells(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	row := AddRow()
	cell := AddCell()
	Value = "A cell (with value 1)!"
	cell = AddCell()
	Value = "A cell (with value 2)!"
	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	xSheet := makeXLSXSheet(refTable, styles)

	output := bytes.NewBufferString(xml.Header)
	body, err := xml.Marshal(xSheet)
	c.Assert(err, IsNil)
	c.Assert(body, NotNil)
	_, err = output.Write(body)
	c.Assert(err, IsNil)

	expectedXLSXSheet := `<?xml version="1.0" encoding="UTF-8"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetPr filterMode="false"><pageSetUpPr fitToPage="false"></pageSetUpPr></sheetPr><dimension ref="A1:B1"></dimension><sheetViews><sheetView windowProtection="false" showFormulas="false" showGridLines="true" showRowColHeaders="true" showZeros="true" rightToLeft="false" tabSelected="true" showOutlineSymbols="true" defaultGridColor="true" view="normal" topLeftCell="A1" colorId="64" zoomScale="100" zoomScaleNormal="100" zoomScalePageLayoutView="100" workbookViewId="0"><selection pane="topLeft" activeCell="A1" activeCellId="0" sqref="A1"></selection></sheetView></sheetViews><sheetFormatPr defaultRowHeight="12.85"></sheetFormatPr><cols><col collapsed="false" hidden="false" max="1" min="1" style="0" width="9.5"></col><col collapsed="false" hidden="false" max="2" min="2" style="0" width="9.5"></col></cols><sheetData><row r="1"><c r="A1" t="s"><v>0</v></c><c r="B1" t="s"><v>1</v></c></row></sheetData><printOptions headings="false" gridLines="false" gridLinesSet="true" horizontalCentered="false" verticalCentered="false"></printOptions><pageMargins left="0.7875" right="0.7875" top="1.05277777777778" bottom="1.05277777777778" header="0.7875" footer="0.7875"></pageMargins><pageSetup paperSize="9" scale="100" firstPageNumber="1" fitToWidth="1" fitToHeight="1" pageOrder="downThenOver" orientation="portrait" usePrinterDefaults="false" blackAndWhite="false" draft="false" cellComments="none" useFirstPageNumber="true" horizontalDpi="300" verticalDpi="300" copies="1"></pageSetup><headerFooter differentFirst="false" differentOddEven="false"><oddHeader>&amp;C&amp;&#34;Times New Roman,Regular&#34;&amp;12&amp;A</oddHeader><oddFooter>&amp;C&amp;&#34;Times New Roman,Regular&#34;&amp;12Page &amp;P</oddFooter></headerFooter></worksheet>`
	c.Assert(output.String(), Equals, expectedXLSXSheet)
}

func (s *SheetSuite) TestSetColWidth(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	_ = SetColWidth(0, 0, 10.5)
	_ = SetColWidth(1, 5, 11)

	c.Assert(Width, Equals, 10.5)
	c.Assert(Max, Equals, 1)
	c.Assert(Min, Equals, 1)
	c.Assert(Width, Equals, float64(11))
	c.Assert(Max, Equals, 6)
	c.Assert(Min, Equals, 2)
}

func (s *SheetSuite) TestSetRowHeightCM(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	row := AddRow()
	SetHeightCM(1.5)
	c.Assert(Height, Equals, 42.51968505)
}

func (s *SheetSuite) TestAlignment(c *C) {
	leftalign := *DefaultAlignment()
	Horizontal = "left"
	centerHalign := *DefaultAlignment()
	Horizontal = "center"
	rightalign := *DefaultAlignment()
	Horizontal = "right"

	file := NewFile()
	sheet, _ := AddSheet("Sheet1")

	style := NewStyle()

	hrow := AddRow()

	// Horizontals
	cell := AddCell()
	Value = "left"
	Alignment = leftalign
	ApplyAlignment = true
	SetStyle(style)

	style = NewStyle()
	cell = AddCell()
	Value = "centerH"
	Alignment = centerHalign
	ApplyAlignment = true
	SetStyle(style)

	style = NewStyle()
	cell = AddCell()
	Value = "right"
	Alignment = rightalign
	ApplyAlignment = true
	SetStyle(style)

	// Verticals
	topalign := *DefaultAlignment()
	Vertical = "top"
	centerValign := *DefaultAlignment()
	Vertical = "center"
	bottomalign := *DefaultAlignment()
	Vertical = "bottom"

	style = NewStyle()
	vrow := AddRow()
	cell = AddCell()
	Value = "top"
	Alignment = topalign
	ApplyAlignment = true
	SetStyle(style)

	style = NewStyle()
	cell = AddCell()
	Value = "centerV"
	Alignment = centerValign
	ApplyAlignment = true
	SetStyle(style)

	style = NewStyle()
	cell = AddCell()
	Value = "bottom"
	Alignment = bottomalign
	ApplyAlignment = true
	SetStyle(style)

	parts, err := MarshallParts()
	c.Assert(err, IsNil)
	obtained := parts["xl/styles.xml"]

	shouldbe := `<?xml version="1.0" encoding="UTF-8"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><fonts count="1"><font><sz val="12"/><name val="Verdana"/><family val="0"/><charset val="0"/></font></fonts><fills count="2"><fill><patternFill patternType="none"><fgColor rgb="FFFFFFFF"/><bgColor rgb="00000000"/></patternFill></fill><fill><patternFill patternType="lightGray"/></fill></fills><borders count="1"><border><left style="none"></left><right style="none"></right><top style="none"></top><bottom style="none"></bottom></border></borders><cellXfs count="8"><xf applyAlignment="0" applyBorder="0" applyFont="0" applyFill="0" applyNumberFormat="0" applyProtection="0" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="general" indent="0" shrinkToFit="0" textRotation="0" vertical="bottom" wrapText="0"/></xf><xf applyAlignment="0" applyBorder="0" applyFont="0" applyFill="0" applyNumberFormat="0" applyProtection="0" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="general" indent="0" shrinkToFit="0" textRotation="0" vertical="bottom" wrapText="0"/></xf><xf applyAlignment="1" applyBorder="0" applyFont="0" applyFill="0" applyNumberFormat="0" applyProtection="0" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="left" indent="0" shrinkToFit="0" textRotation="0" vertical="bottom" wrapText="0"/></xf><xf applyAlignment="1" applyBorder="0" applyFont="0" applyFill="0" applyNumberFormat="0" applyProtection="0" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="center" indent="0" shrinkToFit="0" textRotation="0" vertical="bottom" wrapText="0"/></xf><xf applyAlignment="1" applyBorder="0" applyFont="0" applyFill="0" applyNumberFormat="0" applyProtection="0" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="right" indent="0" shrinkToFit="0" textRotation="0" vertical="bottom" wrapText="0"/></xf><xf applyAlignment="1" applyBorder="0" applyFont="0" applyFill="0" applyNumberFormat="0" applyProtection="0" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="general" indent="0" shrinkToFit="0" textRotation="0" vertical="top" wrapText="0"/></xf><xf applyAlignment="1" applyBorder="0" applyFont="0" applyFill="0" applyNumberFormat="0" applyProtection="0" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="general" indent="0" shrinkToFit="0" textRotation="0" vertical="center" wrapText="0"/></xf><xf applyAlignment="1" applyBorder="0" applyFont="0" applyFill="0" applyNumberFormat="0" applyProtection="0" borderId="0" fillId="0" fontId="0" numFmtId="0"><alignment horizontal="general" indent="0" shrinkToFit="0" textRotation="0" vertical="bottom" wrapText="0"/></xf></cellXfs></styleSheet>`

	expected := bytes.NewBufferString(shouldbe)

	c.Assert(obtained, Equals, expected.String())
}

func (s *SheetSuite) TestBorder(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")
	row := AddRow()

	cell1 := AddCell()
	Value = "A cell!"
	style1 := NewStyle()
	Border = *NewBorder("thin", "thin", "thin", "thin")
	ApplyBorder = true
	SetStyle(style1)

	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	worksheet := makeXLSXSheet(refTable, styles)

	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "thin")

	c.Assert(S, Equals, 1)
}

func (s *SheetSuite) TestOutlineLevels(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")

	r1 := AddRow()
	c11 := AddCell()
	Value = "A1"
	c12 := AddCell()
	Value = "B1"

	r2 := AddRow()
	c21 := AddCell()
	Value = "A2"
	c22 := AddCell()
	Value = "B2"

	r3 := AddRow()
	c31 := AddCell()
	Value = "A3"
	c32 := AddCell()
	Value = "B3"

	// Add some groups
	OutlineLevel = 1
	OutlineLevel = 2
	OutlineLevel = 1

	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	worksheet := makeXLSXSheet(refTable, styles)

	c.Assert(OutlineLevelCol, Equals, uint8(1))
	c.Assert(OutlineLevelRow, Equals, uint8(2))

	c.Assert(OutlineLevel, Equals, uint8(1))
	c.Assert(OutlineLevel, Equals, uint8(0))
	c.Assert(OutlineLevel, Equals, uint8(1))
	c.Assert(OutlineLevel, Equals, uint8(2))
	c.Assert(OutlineLevel, Equals, uint8(0))
}

func (s *SheetSuite) TestAutoFilter(c *C) {
	file := NewFile()
	sheet, _ := AddSheet("Sheet1")

	r1 := AddRow()
	AddCell()
	AddCell()
	AddCell()

	r2 := AddRow()
	AddCell()
	AddCell()
	AddCell()

	r3 := AddRow()
	AddCell()
	AddCell()
	AddCell()

	// Define a filter area
	AutoFilter = &AutoFilter{TopLeftCell: "B2", BottomRightCell: "C3"}

	refTable := NewSharedStringRefTable()
	styles := newXlsxStyleSheet(nil)
	worksheet := makeXLSXSheet(refTable, styles)

	c.Assert(AutoFilter, NotNil)
	c.Assert(Ref, Equals, "B2:C3")
}
