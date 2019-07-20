package xlsx

import (
	"bytes"
	"encoding/xml"
)

type ThemeSuite struct{}

var _ = Suite(&ThemeSuite{})

func (s *ThemeSuite) TestThemeColors(c *C) {
	themeXmlBytes := bytes.NewBufferString(`
<?xml version="1.0"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office Theme">
<a:themeElements>
  <a:clrScheme name="Office">
    <a:dk1>
      <a:sysClr val="windowText" lastClr="000000"/>
    </a:dk1>
    <a:lt1>
      <a:sysClr val="window" lastClr="FFFFFF"/>
    </a:lt1>
    <a:dk2>
      <a:srgbClr val="1F497D"/>
    </a:dk2>
    <a:lt2>
      <a:srgbClr val="EEECE1"/>
    </a:lt2>
    <a:accent1>
      <a:srgbClr val="4F81BD"/>
    </a:accent1>
    <a:accent2>
      <a:srgbClr val="C0504D"/>
    </a:accent2>
    <a:accent3>
      <a:srgbClr val="9BBB59"/>
    </a:accent3>
    <a:accent4>
      <a:srgbClr val="8064A2"/>
    </a:accent4>
    <a:accent5>
      <a:srgbClr val="4BACC6"/>
    </a:accent5>
    <a:accent6>
      <a:srgbClr val="F79646"/>
    </a:accent6>
    <a:hlink>
      <a:srgbClr val="0000FF"/>
    </a:hlink>
    <a:folHlink>
      <a:srgbClr val="800080"/>
    </a:folHlink>
  </a:clrScheme>
</a:themeElements>
</a:theme>
	`)
	var themeXml xlsxTheme
	err := xml.NewDecoder(themeXmlBytes).Decode(&themeXml)
	c.Assert(err, IsNil)

	clrSchemes := Children
	c.Assert(len(clrSchemes), Equals, 12)

	dk1Scheme := clrSchemes[0]
	c.Assert(XMLName.Local, Equals, "dk1")
	c.Assert(SrgbClr, IsNil)
	c.Assert(SysClr, NotNil)
	c.Assert(Val, Equals, "windowText")
	c.Assert(LastClr, Equals, "000000")

	dk2Scheme := clrSchemes[2]
	c.Assert(XMLName.Local, Equals, "dk2")
	c.Assert(SysClr, IsNil)
	c.Assert(SrgbClr, NotNil)
	c.Assert(Val, Equals, "1F497D")

	theme := newTheme(themeXml)
	c.Assert(themeColor(0, 0), Equals, "FFFFFFFF")
	c.Assert(themeColor(2, 0), Equals, "FFEEECE1")
}
