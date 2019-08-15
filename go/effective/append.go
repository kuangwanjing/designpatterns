package main

import (
	"fmt"
)

type ByteSlice []byte

func (slice ByteSlice) Append(data []byte) []byte {
	// Body exactly the same as the Append function defined above.
	l := len(slice)
	if l+len(data) > cap(slice) { // reallocate
		// Allocate double what's needed, for future growth.
		newSlice := make([]byte, (l+len(data))*2)
		// The copy function is predeclared and works for any slice type.
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : l+len(data)]
	copy(slice[l:], data)
	return slice
}

func (p *ByteSlice) Append2(data []byte) {
	slice := *p
	l := len(slice)
	if l+len(data) > cap(slice) { // reallocate
		// Allocate double what's needed, for future growth.
		newSlice := make([]byte, (l+len(data))*2)
		// The copy function is predeclared and works for any slice type.
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : l+len(data)]
	copy(slice[l:], data)
	*p = slice
}

// the ByteSlice can be used as a writer (fmt, bytes.Buffer...)
func (p *ByteSlice) Write(data []byte) (n int, err error) {
	slice := *p
	l := len(slice)
	if l+len(data) > cap(slice) { // reallocate
		// Allocate double what's needed, for future growth.
		newSlice := make([]byte, (l+len(data))*2)
		// The copy function is predeclared and works for any slice type.
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : l+len(data)]
	copy(slice[l:], data)
	*p = slice
	return len(data), nil
}

func main() {
	var bs ByteSlice
	bs = bs.Append([]byte("abc"))
	fmt.Printf("%v\n", string(bs))
	(&bs).Append2([]byte("def"))
	fmt.Printf("%v\n", string(bs))
	var b ByteSlice
	fmt.Fprintf(&b, "This hour has %d days\n", 7) // since Write method is binded on the pointer
	fmt.Printf("%v\n", string(b))
}
