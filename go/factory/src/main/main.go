package main

import (
	"encoding/xml"
	"factory"
	"format"
)

type Person struct {
	XMLName xml.Name `json:"-";xml:"person"`
	Name    string   `json:"name";xml:"name"`
	Age     int      `json:"age";xml:"age"`
}

func RegisterFormatFactories() {
	format.Register("json", factory.JsonFactory)
	format.Register("xml", factory.XmlFactory)
}

func main() {
	//initialize the factories
	RegisterFormatFactories()
	p1 := Person{Name: "John", Age: 21}
	s1, _ := format.Stringify("json", p1)
	println(s1)
	p2 := Person{Name: "Mary", Age: 30}
	s2, _ := format.Stringify("xml", p2)
	println(s2)
}
