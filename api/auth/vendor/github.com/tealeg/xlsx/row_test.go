package xlsx

type RowSuite struct{}

var _ = Suite(&RowSuite{})

// Test we can add a new Cell to a Row
func (r *RowSuite) TestAddCell(c *C) {
	var f *File
	f = NewFile()
	sheet, _ := AddSheet("MySheet")
	row := AddRow()
	cell := AddCell()
	c.Assert(cell, NotNil)
	c.Assert(len(Cells), Equals, 1)
}
