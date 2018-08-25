package method

import (
	"encoding/json"
)

type JsonFormat struct {
	Format
}

func (f JsonFormat) Stringify(data interface{}) string {
	bs, _ := json.Marshal(data)
	return string(bs)
}
