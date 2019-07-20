package xlsx

type StyleSuite struct{}

var _ = Suite(&StyleSuite{})

func (s *StyleSuite) TestNewStyle(c *C) {
	style := NewStyle()
	c.Assert(style, NotNil)
}

func (s *StyleSuite) TestNewStyleDefaultts(c *C) {
	style := NewStyle()
	c.Assert(Font, Equals, *DefaultFont())
	c.Assert(Fill, Equals, *DefaultFill())
	c.Assert(Border, Equals, *DefaultBorder())
}

func (s *StyleSuite) TestMakeXLSXStyleElements(c *C) {
	style := NewStyle()
	font := *NewFont(12, "Verdana")
	Bold = true
	Italic = true
	Underline = true
	Font = font
	fill := *NewFill("solid", "00FF0000", "FF000000")
	Fill = fill
	border := *NewBorder("thin", "thin", "thin", "thin")
	Border = border
	ApplyBorder = true
	ApplyFill = true

	ApplyFont = true
	xFont, xFill, xBorder, xCellXf := makeXLSXStyleElements()
	c.Assert(Val, Equals, "12")
	c.Assert(Val, Equals, "Verdana")
	c.Assert(B, NotNil)
	c.Assert(I, NotNil)
	c.Assert(U, NotNil)
	c.Assert(PatternType, Equals, "solid")
	c.Assert(RGB, Equals, "00FF0000")
	c.Assert(RGB, Equals, "FF000000")
	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "thin")
	c.Assert(Style, Equals, "thin")
	c.Assert(ApplyBorder, Equals, true)
	c.Assert(ApplyFill, Equals, true)
	c.Assert(ApplyFont, Equals, true)

}

type FontSuite struct{}

var _ = Suite(&FontSuite{})

func (s *FontSuite) TestNewFont(c *C) {
	font := NewFont(12, "Verdana")
	c.Assert(font, NotNil)
	c.Assert(Name, Equals, "Verdana")
	c.Assert(Size, Equals, 12)
}
