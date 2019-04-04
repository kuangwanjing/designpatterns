// polymorphism and json
// refer to: https://medium.com/@alain.drolet.0/how-to-unmarshal-an-array-of-json-objects-of-different-types-into-a-go-struct-10eab5f9a3a2
package json

import (
	"encoding/json"
	"fmt"
)

// use an alias for Wheel to prevent infinite recursion on unmarshalling Wheel
type Wheel InnerWheel

type InnerWheel struct {
	Id     string `json:"id"`
	Radius int    `json:"radius"`
	Type   string `json:"type"`
	attr   interface{}
}

type SolidAttributes struct {
	Material string `json:"material"`
}

type SpokedAttributes struct {
	SpokeCount  int `json:"spoke_count"`
	InnerRadius int `json:"inner_radius"`
}

type TireAttributes struct {
	Fill        string `json:"fill_type"`
	InnerRadius int    `json:"inner_radius"`
}

func (w *Wheel) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, (*InnerWheel)(w))
	if err != nil {
		return err
	}
	var err1 error = nil
	var subTypeAttr interface{}
	switch w.Type {
	case "solid":
		subTypeAttr = &SolidAttributes{}
		err1 = json.Unmarshal(b, subTypeAttr)
		break
	case "spoked":
		subTypeAttr = &SpokedAttributes{}
		err1 = json.Unmarshal(b, subTypeAttr)
		break
	case "tire":
		subTypeAttr = &TireAttributes{}
		err1 = json.Unmarshal(b, subTypeAttr)
		break
	default:
		return fmt.Errorf("Wheel.UnmarshalJSON: unexpected type; type = %s", w.Type)
	}
	if err == nil {
		w.attr = subTypeAttr
	}
	return err1
}

func TestTech3() {
	var wheelJson = `
		{
		  "id": "bicycle_wheel_8450",
		  "radius": 20,
		  "type": "spoked",
		  "spoke_count": 30,
		  "inner_radius": 19
		}
	`
	wheelObj := &Wheel{}
	err := json.Unmarshal([]byte(wheelJson), wheelObj)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(wheelObj.attr)
	}

	input := `[{"id":"1234","radius":17,"type":"solid","material":"wood"},{"id":"5432","radius":20,"type":"spoked","spoken_count":30,"inner_radius":19},{"id":"7898","raidus":20,"type":"tire","inner_radius":17,"fill_type":"inflated"}]`

	var wheels []Wheel

	err = json.Unmarshal([]byte(input), &wheels)

	if err == nil {
		for _, w := range wheels {
			fmt.Println(w.attr)
		}
	}
}
