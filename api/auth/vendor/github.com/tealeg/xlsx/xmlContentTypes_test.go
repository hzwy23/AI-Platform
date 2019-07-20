package xlsx

import (
	"encoding/xml"
)

type ContentTypesSuite struct{}

var _ = Suite(&ContentTypesSuite{})

func (l *ContentTypesSuite) TestMarshalContentTypes(c *C) {
	var types xlsxTypes = xlsxTypes{}
	Overrides = make([]xlsxOverride, 1)
	Overrides[0] = xlsxOverride{PartName: "/_rels/.rels", ContentType: "application/vnd.openxmlformats-package.relationships+xml"}
	output, err := xml.Marshal(types)
	stringOutput := xml.Header + string(output)
	c.Assert(err, IsNil)
	expectedContentTypes := `<?xml version="1.0" encoding="UTF-8"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Override PartName="/_rels/.rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"></Override></Types>`
	c.Assert(stringOutput, Equals, expectedContentTypes)
}

func (l *ContentTypesSuite) TestMakeDefaultContentTypes(c *C) {
	var types xlsxTypes = MakeDefaultContentTypes()
	c.Assert(len(Overrides), Equals, 8)
	c.Assert(PartName, Equals, "/_rels/.rels")
	c.Assert(ContentType, Equals, "application/vnd.openxmlformats-package.relationships+xml")
	c.Assert(PartName, Equals, "/docProps/app.xml")
	c.Assert(ContentType, Equals, "application/vnd.openxmlformats-officedocument.extended-properties+xml")
	c.Assert(PartName, Equals, "/docProps/core.xml")
	c.Assert(ContentType, Equals, "application/vnd.openxmlformats-package.core-properties+xml")
	c.Assert(PartName, Equals, "/xl/_rels/workbook.xml.rels")
	c.Assert(ContentType, Equals, "application/vnd.openxmlformats-package.relationships+xml")
	c.Assert(PartName, Equals, "/xl/sharedStrings.xml")
	c.Assert(ContentType, Equals, "application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml")
	c.Assert(PartName, Equals, "/xl/styles.xml")
	c.Assert(ContentType, Equals, "application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml")
	c.Assert(PartName, Equals, "/xl/workbook.xml")
	c.Assert(ContentType, Equals, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml")
	c.Assert(PartName, Equals, "/xl/theme/theme1.xml")
	c.Assert(ContentType, Equals, "application/vnd.openxmlformats-officedocument.theme+xml")

	c.Assert(Extension, Equals, "rels")
	c.Assert(ContentType, Equals, "application/vnd.openxmlformats-package.relationships+xml")
	c.Assert(Extension, Equals, "xml")
	c.Assert(ContentType, Equals, "application/xml")

}
