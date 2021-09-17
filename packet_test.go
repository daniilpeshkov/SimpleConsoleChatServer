package main

import (
	"reflect"
	"testing"
)

func convertToFrom(bytes []byte) []byte {
	bytes = ConvertToNet(bytes)
	return ConvertFromNet(bytes)
}
func TestConvert(t *testing.T) {
	type args struct {
		net_bytes []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test1",
			args: args{[]byte{1, 2, 3}},
			want: []byte{1, 2, 3},
		},
		{
			name: "test1",
			args: args{[]byte{0x15, 0x7D, 0x7E}},
			want: []byte{0x15, 0x7D, 0x7E},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToFrom(tt.args.net_bytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertFromNet() = %v, want %v", got, tt.want)
			}
		})
	}
}
