package xlsx

import "strconv"

// Style is a high level structure intended to provide user access to
// the contents of Style within an XLSX file.
type Style struct {
	Border          Border
	Fill            Fill
	Font            Font
	ApplyBorder     bool
	ApplyFill       bool
	ApplyFont       bool
	ApplyAlignment  bool
	Alignment       Alignment
	NamedStyleIndex *int
}

// Return a new Style structure initialised with the default values.
func NewStyle() *Style {
	return &Style{
		Alignment: *DefaultAlignment(),
		Border:    *DefaultBorder(),
		Fill:      *DefaultFill(),
		Font:      *DefaultFont(),
	}
}

// Generate the underlying XLSX style elements that correspond to the Style.
func (style *Style) makeXLSXStyleElements() (xFont xlsxFont, xFill xlsxFill, xBorder xlsxBorder, xCellXf xlsxXf) {
	xFont = xlsxFont{}
	xFill = xlsxFill{}
	xBorder = xlsxBorder{}
	xCellXf = xlsxXf{}
	Val = strconv.Itoa(style.Font.Size)
	Val = style.Font.Name
	Val = strconv.Itoa(style.Font.Family)
	Val = strconv.Itoa(style.Font.Charset)
	RGB = style.Font.Color
	if style.Font.Bold {
		B = &xlsxVal{}
	} else {
		B = nil
	}
	if style.Font.Italic {
		I = &xlsxVal{}
	} else {
		I = nil
	}
	if style.Font.Underline {
		U = &xlsxVal{}
	} else {
		U = nil
	}
	xPatternFill := xlsxPatternFill{}
	PatternType = style.Fill.PatternType
	RGB = style.Fill.FgColor
	RGB = style.Fill.BgColor
	PatternFill = xPatternFill
	Left = xlsxLine{
		Style: style.Border.Left,
		Color: xlsxColor{RGB: style.Border.LeftColor},
	}
	Right = xlsxLine{
		Style: style.Border.Right,
		Color: xlsxColor{RGB: style.Border.RightColor},
	}
	Top = xlsxLine{
		Style: style.Border.Top,
		Color: xlsxColor{RGB: style.Border.TopColor},
	}
	Bottom = xlsxLine{
		Style: style.Border.Bottom,
		Color: xlsxColor{RGB: style.Border.BottomColor},
	}
	xCellXf = makeXLSXCellElement()
	ApplyBorder = style.ApplyBorder
	ApplyFill = style.ApplyFill
	ApplyFont = style.ApplyFont
	ApplyAlignment = style.ApplyAlignment
	if style.NamedStyleIndex != nil {
		XfId = style.NamedStyleIndex
	}
	return
}

func makeXLSXCellElement() (xCellXf xlsxXf) {
	NumFmtId = 0
	return
}

// Border is a high level structure intended to provide user access to
// the contents of Border Style within an Sheet.
type Border struct {
	Left        string
	LeftColor   string
	Right       string
	RightColor  string
	Top         string
	TopColor    string
	Bottom      string
	BottomColor string
}

func NewBorder(left, right, top, bottom string) *Border {
	return &Border{
		Left:   left,
		Right:  right,
		Top:    top,
		Bottom: bottom,
	}
}

// Fill is a high level structure intended to provide user access to
// the contents of background and foreground color index within an Sheet.
type Fill struct {
	PatternType string
	BgColor     string
	FgColor     string
}

func NewFill(patternType, fgColor, bgColor string) *Fill {
	return &Fill{
		PatternType: patternType,
		FgColor:     fgColor,
		BgColor:     bgColor,
	}
}

type Font struct {
	Size      int
	Name      string
	Family    int
	Charset   int
	Color     string
	Bold      bool
	Italic    bool
	Underline bool
}

func NewFont(size int, name string) *Font {
	return &Font{Size: size, Name: name}
}

type Alignment struct {
	Horizontal   string
	Indent       int
	ShrinkToFit  bool
	TextRotation int
	Vertical     string
	WrapText     bool
}

var defaultFontSize = 12
var defaultFontName = "Verdana"

func SetDefaultFont(size int, name string) {
	defaultFontSize = size
	defaultFontName = name
}

func DefaultFont() *Font {
	return NewFont(defaultFontSize, defaultFontName)
}

func DefaultFill() *Fill {
	return NewFill("none", "FFFFFFFF", "00000000")

}

func DefaultBorder() *Border {
	return NewBorder("none", "none", "none", "none")
}

func DefaultAlignment() *Alignment {
	return &Alignment{
		Horizontal: "general",
		Vertical:   "bottom",
	}
}
