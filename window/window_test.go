package window

import (
	"bytes"
	"io"
	"testing"

	prompt "github.com/c-bata/go-prompt"
)

func TestWindow_PrintFileContents(t *testing.T) {
	type fields struct {
		Size         Size
		Output       io.Writer
		FileContents [][]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "file row + 2 == window row",
			fields: fields{
				Size:         Size{Row: 4, Column: 100},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
			},
			want: []byte("\033[H\033[2JHello World!\nI am bob\n\n\033[H"),
		},
		{
			name: "file row + 1 == window row",
			fields: fields{
				Size:         Size{Row: 3, Column: 100},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
			},
			want: []byte("\033[H\033[2JHello World!\nI am bob\n\033[H"),
		},
		{
			name: "file row  == window row",
			fields: fields{
				Size:         Size{Row: 2, Column: 100},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
			},
			want: []byte("\033[H\033[2JHello World!\n\033[H"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			w := &Window{
				Size:         tt.fields.Size,
				Output:       out,
				FileContents: tt.fields.FileContents,
			}

			if w.PrintFileContents(); out.String() != string(tt.want) {
				t.Errorf("got: %v, want:  %v", out.Bytes(), tt.want)
			}
		})
	}
}

func TestWindow_SetFileContents(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name   string
		args   args
		wantFc [][]byte
	}{
		{
			name:   "normal case",
			args:   args{fileName: "../testdata/test.txt"},
			wantFc: [][]byte{[]byte("11111"), []byte("2222"), []byte("333"), []byte("44"), []byte(""), []byte("5")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{FileContents: nil}
			if err := w.SetFileContents(tt.args.fileName); err == nil {
				for i, bs := range w.FileContents {
					if !bytes.Equal(bs, tt.wantFc[i]) {
						t.Errorf("want[%d] = %s, got[%d] = %s", i, string(tt.wantFc[i]), i, string(bs))
					}
				}
			}
		})
	}
}

func TestWindow_SetFileContentsFileExist(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "file exists", args: args{fileName: "../testdata/test.txt"}, wantErr: false},
		{name: "file not exists", args: args{fileName: "../testdata/non_test.txt"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{}
			if err := w.SetFileContents(tt.args.fileName); tt.wantErr != (err != nil) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err != nil)
			}
		})
	}

}

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
	p := Prompt{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.GetKey(tt.args.b); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
