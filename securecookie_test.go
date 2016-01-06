package securecookie

import (
	"fmt"
	"testing"
)

func TestFormatFiledNormal(t *testing.T) {
	field := FormatField("Hello")
	if field != "5:Hello" {
		t.Error("")
	}
}

func TestFormatFiledZero(t *testing.T) {
	field := FormatField("")
	if field != "0:" {
		t.Error("")
	}
}

func TestConsumeFieldNormal(t *testing.T) {
	value, err := ConsumeField("5:Hello")
	if err != nil {
		t.Error(err)
	}
	if value != "Hello" {
		t.Error("")
	}
}

func TestConsumeFieldZero(t *testing.T) {
	value, err := ConsumeField("0:")
	if err != nil {
		t.Error(err)
	}
	if value != "" {
		t.Error("")
	}
}

func TestConsumeFieldError(t *testing.T) {
	sequence := []string{"0:11", "4:Hello", ":adf", "sdf8", "8:"}
	for _, content := range sequence {
		_, err := ConsumeField(content)
		if err == nil {
			t.Error("")
		}
	}
}

func TestCreateSignedValue(t *testing.T) {
	r := CreateSignedValue("He", "pwd", "121212")
	fmt.Println(r)
}
