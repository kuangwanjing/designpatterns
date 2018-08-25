package factory

import (
	"format/method"
)

func XmlFactory() method.Format {
	return method.XmlFormat{Prefix: "  ", Indent: "    "}
}
