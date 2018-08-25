package method

import (
	"encoding/xml"
)

type XmlFormat struct {
	Format
	Prefix string
	Indent string
}

func (f XmlFormat) Stringify(data interface{}) string {
	bs, _ := xml.MarshalIndent(data, f.Prefix, f.Indent)
	return string(bs)
}
