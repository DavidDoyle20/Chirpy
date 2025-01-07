package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestValidateLength(t *testing.T) {
	testCases := []struct{
		field reflect.Value
		want bool
		wantErr bool
	} {
		{reflect.ValueOf("1234"), true, false},
		{reflect.ValueOf(""), true, false},
		{reflect.ValueOf("sassssssssssssssssssssssassssssssssssssssssssssassssssssssssssssssssssassssssssssssssssssssssassssssssssssssssssssssassssssssssssssssssssssssssssssssssssssssssssssssssssss"), false, false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.field), func(t *testing.T) {
			if got := validateLength(tc.field.String()); got != tc.want {
				t.Errorf("validateLength(%s) = %t", tc.field.String(), got)
			}
		})
	}
}

func TestReplaceBadWords(t *testing.T) {
	testCases := []struct{
		a string
		want string
		wantErr bool
	} {
		{"", "", false},
		{"body", "body", false},
		{"kerfuffle", "****", false},
		{"body kerfuffle", "body ****", false},
		{"kerfuffle body", "**** body", false},
		{"sharbert", "****", false},
		{"fornax", "****", false},
		{"kerfuffle sharbert", "**** ****", false},
		{"kerfufflesharbert", "kerfufflesharbert", false},
		{"kerfuffle!", "kerfuffle!", false},
		{"KerFufflE", "****", false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.a), func(t *testing.T) {
			if got := replaceBadWords(tc.a); got != tc.want {
				t.Errorf("replaceBadWords(%s) => %s", tc.a, got)
			}
		})
	}
}