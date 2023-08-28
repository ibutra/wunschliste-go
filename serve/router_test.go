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

	if !match("/", "GET", r) {
		t.Fail()
	}
	if match("/index", "GET", r) {
		t.Fail()
	}
	if match("/", "POST", r) {
		t.Fail()
	}
	
	r,err = http.NewRequest("POST", "/index", nil)
	if err != nil {
		t.Fatal()
	}
	if !match("/index", "POST", r) {
		t.Fail()
	}
	if match("/index", "PUT", r) {
		t.Fail()
	}
	if match("/index", "POST", r, &id) {
		t.Fail()
	}
	if match("/index", "POST", r, &txt) {
		t.Fail()
	}

	r,err = http.NewRequest("GET", "/book/12", nil)
	if err != nil {
		t.Fatal()
	}
	if !match("/book/:id", "GET", r, &id) {
		t.Fail()
	}
	if id != 12 {
		t.Fail()
	}
	if !match("/book/:id", "GET", r, &txt) {
		t.Fail()
	}
	if txt != "12" {
		t.Fail()
	}

	r, err = http.NewRequest("GET", "/book/identifier42", nil)
	if err != nil {
		t.Fatal()
	}
	if !match("/book/:string", "GET", r, &txt) {
		t.Fail()
	}
	if txt != "identifier42" {
		t.Fail()
	}
	if match("/book/:id", "GET", r, &id) {
		t.Fail()
	}
}
