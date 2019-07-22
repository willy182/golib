package golib

import (
	"testing"
)

func TestParseToFormValue(t *testing.T) {
	type Source struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Gender string
		Number int `json:"number"`
	}
	type args struct {
		source interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantForm string
		wantErr  bool
	}{
		{
			name: "Testcase #1: Positive",
			args: args{
				source: &Source{ID: "10", Name: "agungdp", Gender: "L", Number: 28},
			},
			wantForm: "Gender=L&id=10&name=agungdp&number=28",
			wantErr:  false,
		},
		{
			name: "Testcase #2: Negative, source is not struct",
			args: args{
				source: "test",
			},
			wantForm: "Gender=L&id=10&name=agungdp&number=28",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotForm, err := ParseToFormValue(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToFormValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && gotForm.Encode() != tt.wantForm {
				t.Errorf("ParseToFormValue() = %v, want %v", gotForm.Encode(), tt.wantForm)
			}
		})
	}
}
