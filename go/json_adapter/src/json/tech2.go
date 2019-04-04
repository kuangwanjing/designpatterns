/*
introduce the second technique of customize json encoding and decoding:
extract similar objects with minor differences from the json text.
https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
*/
package json

import (
	"encoding/json"
	"errors"
	"fmt"
)

type BankAccount struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	RoutingNumber string `json:"routing_number"`
}

type Card struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Last4  string `json:"last4"`
}

// data contains two types of objects
type Data struct {
	*Card
	*BankAccount
}

func (d Data) MarshalJSON() ([]byte, error) {
	if d.Card != nil {
		return json.Marshal(d.Card)
	} else if d.BankAccount != nil {
		return json.Marshal(d.BankAccount)
	} else {
		return json.Marshal(nil)
	}
}

// twice decoding: first use a temporary struct to extract the data type then decode with corresponding template
func (d *Data) UnmarshalJSON(data []byte) error {
	temp := struct {
		Object string `json:"object"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Object == "card" {
		var c Card
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		d.Card = &c
		d.BankAccount = nil
	} else if temp.Object == "bank_account" {
		var ba BankAccount
		if err := json.Unmarshal(data, &ba); err != nil {
			return err
		}
		d.BankAccount = &ba
		d.Card = nil
	} else {
		return errors.New("Invalid object value")
	}
	return nil
}

func TestTech2() {
	jsonStr := `
		  [{
			"object": "bank_account",
			"id": "bank_123",
			"routing_number": "4243"
		  },{
			"object": "card",
			"id": "card_123",
			"last4": "4242"
		  }]
	`
	var data []Data
	json.Unmarshal([]byte(jsonStr), &data)
	for _, d := range data {
		if d.Card != nil {
			fmt.Println(d.Card)
		} else if d.BankAccount != nil {
			fmt.Println(d.BankAccount)
		}
	}
}
