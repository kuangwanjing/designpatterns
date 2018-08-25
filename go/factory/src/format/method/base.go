package method

type Format interface {
	Stringify(interface{}) string
}
