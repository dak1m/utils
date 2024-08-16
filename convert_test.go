package utils

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMapConvertStructByTag(t *testing.T) {
	type args struct {
		input map[string]string
		obj   interface{}
		tag   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestMapConvertStructByTag",
			args: args{
				input: map[string]string{
					"account":  "[{\"balance\": \"200\",\"other_balance\": \"\"}]",
					"birthday": "2020-01-01",
					"ptr":      "111",
				},
				obj: &TestStruct{},
				tag: "json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := MapConvertStructByTag(tt.args.input, tt.args.obj, tt.args.tag); (err != nil) != tt.wantErr {
					t.Errorf("MapConvertStructByTag() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					t.Logf("%v", tt.args.obj)
				}
			},
		)
	}
}

type TestStruct struct {
	Account  *CustomStructList `json:"account"`
	Birthday Birthday          `json:"birthday"`
	Ptr      *string           `json:"ptr"`
}

type CustomStruct struct {
	Balance      string `json:"balance"`
	OtherBalance string `json:"other_balance"`
}

type Birthday string

var _ ConversionFrom = (*CustomStructList)(nil)

type CustomStructList []CustomStruct

func (c *CustomStructList) FromSource(val string) error {
	var list CustomStructList
	err := json.Unmarshal([]byte(val), &list)

	for i := range list {
		if list[i].OtherBalance == "" {
			list[i].OtherBalance = "0"
		}
	}

	*c = list

	return err
}

func TestStructConvertMapByTag(t *testing.T) {
	type args struct {
		obj interface{}
		tag string
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "TestStructConvertMapByTag",
			args: args{
				obj: TestStruct{
					Account: &CustomStructList{
						{
							Balance:      "200",
							OtherBalance: "0",
						},
					},
					Ptr: nil,
				},
				tag: "json",
			},
		}}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				//ptr := "111"
				//a := tt.args.obj.(TestStruct)
				//a.Ptr = &ptr
				//tt.args.obj = a
				if got := StructConvertMapByTag(tt.args.obj, tt.args.tag); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("StructConvertMapByTag() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
