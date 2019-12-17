package main

import (
	"testing"

	prompt "github.com/c-bata/go-prompt"
)

func TestGetKey(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want prompt.Key
	}{
		{name: "type Ctr + C", args: args{b: []byte{0x1b}}, want: prompt.Escape},
		{name: "type Up", args: args{b: []byte{0x1b, 0x5b, 0x41}}, want: prompt.Up},
		{name: "type Down", args: args{b: []byte{0x1b, 0x5b, 0x42}}, want: prompt.Down},
		{name: "type Right", args: args{b: []byte{0x1b, 0x5b, 0x43}}, want: prompt.Right},
		{name: "type Left", args: args{b: []byte{0x1b, 0x5b, 0x44}}, want: prompt.Left},
		{name: "type ControlC", args: args{b: []byte{0x3}}, want: prompt.ControlC},
		{name: "type other(A)", args: args{b: []byte("A")}, want: prompt.NotDefined},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKey(tt.args.b); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
