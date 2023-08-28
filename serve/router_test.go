package serve

import (
	"net/http"
	"testing"
)

func TestMatch(t *testing.T) {
	var id int
	var txt string
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal()
	}

	if !match("/", []string{"GET"}, r) {
		t.Fail()
	}
	if match("/index", []string{"GET"}, r) {
		t.Fail()
	}
	if match("/", []string{"POST"}, r) {
		t.Fail()
	}
	
	r,err = http.NewRequest("POST", "/index", nil)
	if err != nil {
		t.Fatal()
	}
	if !match("/index", []string{"POST"}, r) {
		t.Fail()
	}
	if match("/index", []string{"PUT"}, r) {
		t.Fail()
	}
	if match("/index", []string{"POST"}, r, &id) {
		t.Fail()
	}
	if match("/index", []string{"POST"}, r, &txt) {
		t.Fail()
	}

	r,err = http.NewRequest("GET", "/book/12", nil)
	if err != nil {
		t.Fatal()
	}
	if !match("/book/:id", []string{"GET"}, r, &id) {
		t.Fail()
	}
	if id != 12 {
		t.Fail()
	}
	if !match("/book/:id", []string{"GET"}, r, &txt) {
		t.Fail()
	}
	if txt != "12" {
		t.Fail()
	}

	r, err = http.NewRequest("GET", "/book/identifier42", nil)
	if err != nil {
		t.Fatal()
	}
	if !match("/book/:string", []string{"GET"}, r, &txt) {
		t.Fail()
	}
	if txt != "identifier42" {
		t.Fail()
	}
	if match("/book/:id", []string{"GET"}, r, &id) {
		t.Fail()
	}
}
