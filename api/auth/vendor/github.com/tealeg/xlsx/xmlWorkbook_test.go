package xlsx

import (
	"bytes"
	"encoding/xml"
)

type WorkbookSuite struct{}

var _ = Suite(&WorkbookSuite{})

// Test we can succesfully unmarshal the workbook.xml file from within
// an XLSX file and return a xlsxWorkbook struct (and associated
// children).
func (w *WorkbookSuite) TestUnmarshallWorkbookXML(c *C) {
	var buf = bytes.NewBufferString(
		`<?xml version="1.0"
        encoding="UTF-8"
        standalone="yes"?>
        <workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"
                  xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
          <fileVersion appName="xl"
                       lastEdited="4"
                       lowestEdited="4"
                       rupBuild="4506"/>
          <workbookPr defaultThemeVersion="124226" date1904="true"/>
          <bookViews>
            <workbookView xWindow="120"
                          yWindow="75"
                          windowWidth="15135"
                          windowHeight="7620"/>
          </bookViews>
          <sheets>
            <sheet name="Sheet1"
                   sheetId="1"
                   r:id="rId1"
                   state="visible"/>
            <sheet name="Sheet2"
                   sheetId="2"
                   r:id="rId2"
                   state="hidden"/>
            <sheet name="Sheet3"
                   sheetId="3"
                   r:id="rId3"
                   state="veryHidden"/>
          </sheets>
          <definedNames>
            <definedName name="monitors" comment="this is the comment"
                         description="give cells a name"
                         localSheetId="0">Sheet1!$A$1533</definedName>
          </definedNames>
          <calcPr calcId="125725"/>
          </workbook>`)
	var workbook *xlsxWorkbook
	workbook = new(xlsxWorkbook)
	err := xml.NewDecoder(buf).Decode(workbook)
	c.Assert(err, IsNil)
	c.Assert(AppName, Equals, "xl")
	c.Assert(LastEdited, Equals, "4")
	c.Assert(LowestEdited, Equals, "4")
	c.Assert(RupBuild, Equals, "4506")
	c.Assert(DefaultThemeVersion, Equals, "124226")
	c.Assert(Date1904, Equals, true)
	c.Assert(WorkBookView, HasLen, 1)
	workBookView := WorkBookView[0]
	c.Assert(XWindow, Equals, "120")
	c.Assert(YWindow, Equals, "75")
	c.Assert(WindowWidth, Equals, 15135)
	c.Assert(WindowHeight, Equals, 7620)
	c.Assert(Sheet, HasLen, 3)
	sheet := Sheet[0]
	c.Assert(Id, Equals, "rId1")
	c.Assert(Name, Equals, "Sheet1")
	c.Assert(SheetId, Equals, "1")
	c.Assert(State, Equals, "visible")
	c.Assert(DefinedName, HasLen, 1)
	dname := DefinedName[0]
	c.Assert(Data, Equals, "Sheet1!$A$1533")
	c.Assert(LocalSheetID, Equals, 0)
	c.Assert(Name, Equals, "monitors")
	c.Assert(Comment, Equals, "this is the comment")
	c.Assert(Description, Equals, "give cells a name")
	c.Assert(CalcId, Equals, "125725")
}

// Test we can marshall a Workbook to xml
func (w *WorkbookSuite) TestMarshallWorkbook(c *C) {
	var workbook *xlsxWorkbook
	workbook = new(xlsxWorkbook)
	FileVersion = xlsxFileVersion{}
	AppName = "xlsx"
	WorkbookPr = xlsxWorkbookPr{BackupFile: false}
	BookViews = xlsxBookViews{}
	WorkBookView = make([]xlsxWorkBookView, 1)
	WorkBookView[0] = xlsxWorkBookView{}
	Sheets = xlsxSheets{}
	Sheet = make([]xlsxSheet, 1)
	Sheet[0] = xlsxSheet{Name: "sheet1", SheetId: "1", Id: "rId2"}

	body, err := xml.Marshal(workbook)
	c.Assert(err, IsNil)
	expectedWorkbook := `<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><fileVersion appName="xlsx"></fileVersion><workbookPr date1904="false"></workbookPr><workbookProtection></workbookProtection><bookViews><workbookView></workbookView></bookViews><sheets><sheet name="sheet1" sheetId="1" xmlns:relationships="http://schemas.openxmlformats.org/officeDocument/2006/relationships" relationships:id="rId2"></sheet></sheets><definedNames></definedNames><calcPr></calcPr></workbook>`
	c.Assert(string(body), Equals, expectedWorkbook)
}
