package app

import "testing"

func TestNewPagination(t *testing.T) {
	pagination := NewPagination(5, 20, 200)

	if pagination.HasNext != true {
		t.Fail()
	}

	if pagination.HasPrev != true {
		t.Fail()
	}
}
