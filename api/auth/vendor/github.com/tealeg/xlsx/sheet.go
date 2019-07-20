package xlsx

import (
	"fmt"
	"strconv"
)

// Sheet is a high level structure intended to provide user access to
// the contents of a particular sheet within an XLSX file.
type Sheet struct {
	Name        string
	File        *File
	Rows        []*Row
	Cols        []*Col
	MaxRow      int
	MaxCol      int
	Hidden      bool
	Selected    bool
	SheetViews  []SheetView
	SheetFormat SheetFormat
	AutoFilter  *AutoFilter
}

type SheetView struct {
	Pane *Pane
}

type Pane struct {
	XSplit      float64
	YSplit      float64
	TopLeftCell string
	ActivePane  string
	State       string // Either "split" or "frozen"
}

type SheetFormat struct {
	DefaultColWidth  float64
	DefaultRowHeight float64
	OutlineLevelCol  uint8
	OutlineLevelRow  uint8
}

type AutoFilter struct {
	TopLeftCell     string
	BottomRightCell string
}

// Add a new Row to a Sheet
func (s *Sheet) AddRow() *Row {
	row := &Row{Sheet: s}
	s.Rows = append(s.Rows, row)
	if len(s.Rows) > s.MaxRow {
		s.MaxRow = len(s.Rows)
	}
	return row
}

// Make sure we always have as many Cols as we do cells.
func (s *Sheet) maybeAddCol(cellCount int) {
	if cellCount > s.MaxCol {
		col := &Col{
			style:     NewStyle(),
			Min:       cellCount,
			Max:       cellCount,
			Hidden:    false,
			Collapsed: false}
		s.Cols = append(s.Cols, col)
		s.MaxCol = cellCount
	}
}

// Make sure we always have as many Cols as we do cells.
func (s *Sheet) Col(idx int) *Col {
	s.maybeAddCol(idx + 1)
	return s.Cols[idx]
}

// GetDetails a Cell by passing it's cartesian coordinates (zero based) as
// row and column integer indexes.
//
// For example:
//
//    cell := sheet.Cell(0,0)
//
// ... would set the variable "cell" to contain a Cell struct
// containing the data from the field "A1" on the spreadsheet.
func (sh *Sheet) Cell(row, col int) *Cell {

	// If the user requests a row beyond what we have, then extend.
	for len(sh.Rows) <= row {
		sh.AddRow()
	}

	r := sh.Rows[row]
	for len(Cells) <= col {
		AddCell()
	}

	return Cells[col]
}

//Set the width of a single column or multiple columns.
func (s *Sheet) SetColWidth(startcol, endcol int, width float64) error {
	if startcol > endcol {
		return fmt.Errorf("Could not set width for range %d-%d: startcol must be less than endcol.", startcol, endcol)
	}
	col := &Col{
		style:     NewStyle(),
		Min:       startcol + 1,
		Max:       endcol + 1,
		Hidden:    false,
		Collapsed: false,
		Width:     width}
	s.Cols = append(s.Cols, col)
	if endcol+1 > s.MaxCol {
		s.MaxCol = endcol + 1
	}
	return nil
}

// When merging cells, the cell may be the 'original' or the 'covered'.
// First, figure out which cells are merge starting points. Then create
// the necessary cells underlying the merge area.
// Then go through all the underlying cells and apply the appropriate
// border, based on the original cell.
func (s *Sheet) handleMerged() {
	merged := make(map[string]*Cell)

	for r, row := range s.Rows {
		for c, cell := range Cells {
			if HMerge > 0 || VMerge > 0 {
				coord := fmt.Sprintf("%s%d", numericToLetters(c), r+1)
				merged[coord] = cell
			}
		}
	}

	// This loop iterates over all cells that should be merged and applies the correct
	// borders to them depending on their position. If any cells required by the merge
	// are missing, they will be allocated by s.Cell().
	for key, cell := range merged {
		mainstyle := GetStyle()

		top := Top
		left := Left
		right := Right
		bottom := Bottom

		// When merging cells, the upper left cell does not maintain
		// the original borders
		Top = "none"
		Left = "none"
		Right = "none"
		Bottom = "none"

		maincol, mainrow, _ := GetCoordsFromCellIDString(key)
		for rownum := 0; rownum <= VMerge; rownum++ {
			for colnum := 0; colnum <= HMerge; colnum++ {
				tmpcell := s.Cell(mainrow+rownum, maincol+colnum)
				style := GetStyle()
				ApplyBorder = true

				if rownum == 0 {
					Top = top
				}

				if rownum == (VMerge) {
					Bottom = bottom
				}

				if colnum == 0 {
					Left = left
				}

				if colnum == (HMerge) {
					Right = right
				}
			}
		}
	}
}

// Dump sheet to its XML representation, intended for internal use only
func (s *Sheet) makeXLSXSheet(refTable *RefTable, styles *xlsxStyleSheet) *xlsxWorksheet {
	worksheet := newXlsxWorksheet()
	xSheet := xlsxSheetData{}
	maxRow := 0
	maxCell := 0
	var maxLevelCol, maxLevelRow uint8

	// Scan through the sheet and see if there are any merged cells. If there
	// are, we may need to extend the size of the sheet. There needs to be
	// phantom cells underlying the area covered by the merged cell
	s.handleMerged()

	for index, sheetView := range s.SheetViews {
		if sheetView.Pane != nil {
			Pane = &xlsxPane{
				XSplit:      sheetView.Pane.XSplit,
				YSplit:      sheetView.Pane.YSplit,
				TopLeftCell: sheetView.Pane.TopLeftCell,
				ActivePane:  sheetView.Pane.ActivePane,
				State:       sheetView.Pane.State,
			}

		}
	}

	if s.Selected {
		TabSelected = true
	}

	if s.SheetFormat.DefaultRowHeight != 0 {
		DefaultRowHeight = s.SheetFormat.DefaultRowHeight
	}
	DefaultColWidth = s.SheetFormat.DefaultColWidth

	colsXfIdList := make([]int, len(s.Cols))
	Cols = &xlsxCols{Col: []xlsxCol{}}
	for c, col := range s.Cols {
		XfId := 0
		if Min == 0 {
			Min = 1
		}
		if Max == 0 {
			Max = 1
		}
		style := GetStyle()
		//col's style always not nil
		if style != nil {
			xNumFmt := newNumFmt(numFmt)
			XfId = handleStyleForXLSX(style, NumFmtId, styles)
		}
		colsXfIdList[c] = XfId

		var customWidth bool
		if Width == 0 {
			Width = ColWidth
			customWidth = false

		} else {
			customWidth = true
		}
		Col = append(Col,
			xlsxCol{Min: Min,
				Max:          Max,
				Hidden:       Hidden,
				Width:        Width,
				CustomWidth:  customWidth,
				Collapsed:    Collapsed,
				OutlineLevel: OutlineLevel,
				Style:        XfId,
			})

		if OutlineLevel > maxLevelCol {
			maxLevelCol = OutlineLevel
		}
	}

	for r, row := range s.Rows {
		if r > maxRow {
			maxRow = r
		}
		xRow := xlsxRow{}
		R = r + 1
		if isCustom {
			CustomHeight = true
			Ht = fmt.Sprintf("%g", Height)
		}
		OutlineLevel = OutlineLevel
		if OutlineLevel > maxLevelRow {
			maxLevelRow = OutlineLevel
		}
		for c, cell := range Cells {
			XfId := colsXfIdList[c]

			// generate NumFmtId and add new NumFmt
			xNumFmt := newNumFmt(NumFmt)

			style := style
			if style != nil {
				XfId = handleStyleForXLSX(style, NumFmtId, styles)
			} else if len(NumFmt) > 0 && numFmt != NumFmt {
				XfId = handleNumFmtIdForXLSX(NumFmtId, styles)
			}

			if c > maxCell {
				maxCell = c
			}
			xC := xlsxC{}
			R = fmt.Sprintf("%s%d", numericToLetters(c), r+1)
			switch cellType {
			case CellTypeString:
				if len(Value) > 0 {
					V = strconv.Itoa(AddString(Value))
				}
				T = "s"
				S = XfId
			case CellTypeBool:
				V = Value
				T = "b"
				S = XfId
			case CellTypeNumeric:
				V = Value
				S = XfId
			case CellTypeDate:
				V = Value
				S = XfId
			case CellTypeFormula:
				V = Value
				F = &xlsxF{Content: formula}
				S = XfId
			case CellTypeError:
				V = Value
				F = &xlsxF{Content: formula}
				T = "e"
				S = XfId
			case CellTypeGeneral:
				V = Value
				S = XfId
			}

			C = append(C, xC)

			if HMerge > 0 || VMerge > 0 {
				// r == rownum, c == colnum
				mc := xlsxMergeCell{}
				start := fmt.Sprintf("%s%d", numericToLetters(c), r+1)
				endcol := c + HMerge
				endrow := r + VMerge + 1
				end := fmt.Sprintf("%s%d", numericToLetters(endcol), endrow)
				Ref = start + ":" + end
				if MergeCells == nil {
					MergeCells = &xlsxMergeCells{}
				}
				Cells = append(Cells, mc)
			}
		}
		Row = append(Row, xRow)
	}

	// Post sheet format with the freshly determined max levels
	s.SheetFormat.OutlineLevelCol = maxLevelCol
	s.SheetFormat.OutlineLevelRow = maxLevelRow
	// .. and then also apply this to the xml worksheet
	OutlineLevelCol = s.SheetFormat.OutlineLevelCol
	OutlineLevelRow = s.SheetFormat.OutlineLevelRow

	if MergeCells != nil {
		Count = len(Cells)
	}

	if s.AutoFilter != nil {
		AutoFilter = &xlsxAutoFilter{Ref: fmt.Sprintf("%v:%v", s.AutoFilter.TopLeftCell, s.AutoFilter.BottomRightCell)}
	}

	SheetData = xSheet
	dimension := xlsxDimension{}
	Ref = fmt.Sprintf("A1:%s%d",
		numericToLetters(maxCell), maxRow+1)
	if Ref == "A1:A1" {
		Ref = "A1"
	}
	Dimension = dimension
	return worksheet
}

func handleStyleForXLSX(style *Style, NumFmtId int, styles *xlsxStyleSheet) (XfId int) {
	xFont, xFill, xBorder, xCellXf := makeXLSXStyleElements()
	fontId := addFont(xFont)
	fillId := addFill(xFill)

	// HACK - adding light grey fill, as in OO and Google
	greyfill := xlsxFill{}
	PatternType = "lightGray"
	addFill(greyfill)

	borderId := addBorder(xBorder)
	FontId = fontId
	FillId = fillId
	BorderId = borderId
	NumFmtId = NumFmtId
	// apply the numFmtId when it is not the default cellxf
	if NumFmtId > 0 {
		ApplyNumberFormat = true
	}

	Horizontal = Horizontal
	Indent = Indent
	ShrinkToFit = ShrinkToFit
	TextRotation = TextRotation
	Vertical = Vertical
	WrapText = WrapText

	XfId = addCellXf(xCellXf)
	return
}

func handleNumFmtIdForXLSX(NumFmtId int, styles *xlsxStyleSheet) (XfId int) {
	xCellXf := makeXLSXCellElement()
	NumFmtId = NumFmtId
	if NumFmtId > 0 {
		ApplyNumberFormat = true
	}
	XfId = addCellXf(xCellXf)
	return
}
