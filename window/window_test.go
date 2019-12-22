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
	w := &Window{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := w.GetKey(tt.args.b); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_MoveDown(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	type args struct {
		num int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{name: "normal test", fields: fields{X: 1, Y: 1}, args: args{num: 1}, want: fields{X: 1, Y: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Position{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if p.MoveDown(tt.args.num); p.X != tt.want.X || p.Y != tt.want.Y {
				t.Errorf("got: X = %d, Y = %d; want: X = %d, Y = %d", p.X, p.Y, tt.want.X, tt.want.Y)
			}
		})
	}
}

func TestPosition_MoveUp(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	type args struct {
		num int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{name: "normal test", fields: fields{X: 2, Y: 2}, args: args{num: 1}, want: fields{X: 2, Y: 1}},
		{name: "no change if Y is 0", fields: fields{X: 2, Y: 1}, args: args{num: 1}, want: fields{X: 2, Y: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Position{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if p.MoveUp(tt.args.num); p.X != tt.want.X || p.Y != tt.want.Y {
				t.Errorf("got: X = %d, Y = %d; want: X = %d, Y = %d", p.X, p.Y, tt.want.X, tt.want.Y)
			}
		})
	}
}

func TestPosition_MoveRight(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	type args struct {
		num int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{name: "normal test", fields: fields{X: 1, Y: 1}, args: args{num: 1}, want: fields{X: 2, Y: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Position{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if p.MoveRight(tt.args.num); p.X != tt.want.X || p.Y != tt.want.Y {
				t.Errorf("got: X = %d, Y = %d; want: X = %d, Y = %d", p.X, p.Y, tt.want.X, tt.want.Y)
			}
		})
	}
}

func TestPosition_MoveLeft(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	type args struct {
		num int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{name: "normal test", fields: fields{X: 2, Y: 2}, args: args{num: 1}, want: fields{X: 1, Y: 2}},
		{name: "no change if X is 0", fields: fields{X: 1, Y: 2}, args: args{num: 1}, want: fields{X: 1, Y: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Position{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if p.MoveLeft(tt.args.num); p.X != tt.want.X || p.Y != tt.want.Y {
				t.Errorf("got: X = %d, Y = %d; want: X = %d, Y = %d", p.X, p.Y, tt.want.X, tt.want.Y)
			}
		})
	}
}
