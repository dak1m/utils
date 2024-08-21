package utils

import (
	"encoding/json"
	"testing"
)

func TestAmount_UnmarshalJSON(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		result  MarshalObj
		wantErr bool
	}{
		{
			args: args{
				bytes: []byte("{\"amount\": \"200.1\"}"),
			},
			result:  MarshalObj{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal(tt.args.bytes, &tt.result); (err != nil) != tt.wantErr {
				t.Errorf("MapConvertStructByTag() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("%v", tt.result)
			}
		})
	}
}

func TestAmount_MarshalJSON(t *testing.T) {
	type args struct {
		obj MarshalObj
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				obj: MarshalObj{
					Amount: "200.000000",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result, err := json.Marshal(&tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("MapConvertStructByTag() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("%v", string(result))
			}
		})
	}
}

type MarshalObj struct {
	Amount Amount `json:"amount"`
}
