package window

import (
	"bytes"
	"io"
	"testing"
)

func TestWindow_PrintFileContents(t *testing.T) {
	type fields struct {
		Size   Size
		Output io.Writer
	}
	type args struct {
		fc [][]byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "file row + 2 == window row",
			fields: fields{
				Size: Size{Row: 4, Column: 100},
			},
			args: args{fc: [][]byte{[]byte("Hello World!"), []byte("I am bob")}},
			want: []byte("\033[H\033[2JHello World!\nI am bob\n\n\033[H"),
		},
		{
			name: "file row + 1 == window row",
			fields: fields{
				Size: Size{Row: 3, Column: 100},
			},
			args: args{fc: [][]byte{[]byte("Hello World!"), []byte("I am bob")}},
			want: []byte("\033[H\033[2JHello World!\nI am bob\n\033[H"),
		},
		{
			name: "file row  == window row",
			fields: fields{
				Size: Size{Row: 2, Column: 100},
			},
			args: args{fc: [][]byte{[]byte("Hello World!"), []byte("I am bob")}},
			want: []byte("\033[H\033[2JHello World!\n\033[H"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			w := &Window{
				Size:   tt.fields.Size,
				Output: out,
			}

			if w.PrintFileContents(tt.args.fc); out.String() != string(tt.want) {
				t.Errorf("got: %v, want:  %v", out.Bytes(), tt.want)
			}
		})
	}
}
