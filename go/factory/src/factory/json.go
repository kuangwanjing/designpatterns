package factory

import (
	"format/method"
)

func JsonFactory() method.Format {
	return method.JsonFormat{}
}
