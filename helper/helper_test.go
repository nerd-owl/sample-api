package helper

import "testing"

func TestCheckName(t *testing.T) {
	result := CheckName("Savez")

	if !result {
		t.Fail()
	}
}

func TestPhone(t *testing.T) {
	result := CheckPhone("7408963464")

	if !result {
		t.Fail()
	}
}
