package window

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"

	prompt "github.com/c-bata/go-prompt"
)

func TestNewWindow(t *testing.T) {
	tests := []struct {
		name     string
		wantX    int
		wantY    int
		wantMode int
	}{
		{name: "normal test", wantX: 1, wantY: 1, wantMode: normalMode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWindow(os.Stdin, os.Stdout)
			if got.position.X != tt.wantX || got.position.Y != tt.wantY {
				t.Errorf(
					"got: X=%d, Y=%d, want: X=%d, Y=%d", got.position.X, got.position.Y, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestWindow_IsInsertMode(t *testing.T) {
	type fields struct {
		mode int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "normal mode", fields: fields{mode: normalMode}, want: false},
		{name: "insert mode", fields: fields{mode: insertMode}, want: true},
		{name: "command mode", fields: fields{mode: commandMode}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{
				mode: tt.fields.mode,
			}
			if got := w.IsInsertMode(); got != tt.want {
				t.Errorf("IsInsertMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow_IsNormalMode(t *testing.T) {
	type fields struct {
		mode int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "normal mode", fields: fields{mode: normalMode}, want: true},
		{name: "insert mode", fields: fields{mode: insertMode}, want: false},
		{name: "command mode", fields: fields{mode: commandMode}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{
				mode: tt.fields.mode,
			}
			if got := w.IsNormalMode(); got != tt.want {
				t.Errorf("IsNormalMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow_IsCommandMode(t *testing.T) {
	type fields struct {
		mode int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "normal mode", fields: fields{mode: normalMode}, want: false},
		{name: "insert mode", fields: fields{mode: insertMode}, want: false},
		{name: "command mode", fields: fields{mode: commandMode}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{
				mode: tt.fields.mode,
			}
			if got := w.IsCommandMode(); got != tt.want {
				t.Errorf("IsCommandMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow_IsCommandNotTyped(t *testing.T) {
	type fields struct {
		command []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "not typed", fields: fields{command: []byte{}}, want: true},
		{name: "typed", fields: fields{command: []byte("wq")}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{
				command: tt.fields.command,
			}
			if got := w.IsCommandNotTyped(); got != tt.want {
				t.Errorf("IsCommandNotTyped() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow_SetInsertMode(t *testing.T) {
	type fields struct {
		mode int
	}
	tests := []struct {
		name     string
		fields   fields
		wantMode int
	}{
		{name: "normal case", fields: fields{mode: normalMode}, wantMode: insertMode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{mode: tt.fields.mode}
			if w.SetInsertMode(); tt.wantMode != w.mode {
				t.Errorf("got: mode=%d, want: mode=%d", w.mode, tt.wantMode)
			}
		})
	}
}

func TestWindow_SetNormalMode(t *testing.T) {
	tests := []struct {
		name     string
		wantMode int
	}{
		{name: "normal case", wantMode: normalMode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{}
			if w.SetNormalMode(); tt.wantMode != w.mode {
				t.Errorf("got: mode=%d, want: mode=%d", w.mode, tt.wantMode)
			}
		})
	}
}

func TestWindow_SetCommandMode(t *testing.T) {
	tests := []struct {
		name     string
		wantMode int
	}{
		{name: "normal case", wantMode: commandMode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{}
			if w.SetCommandMode(); tt.wantMode != w.mode {
				t.Errorf("got: mode=%d, want: mode=%d", w.mode, tt.wantMode)
			}
		})
	}
}

func TestWindow_AddCommand(t *testing.T) {
	type fields struct {
		command []byte
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantCommand []byte
	}{
		{name: "typed", fields: fields{command: []byte("w")}, args: args{[]byte("q")}, wantCommand: []byte("wq")},
		{name: "not typed", fields: fields{command: []byte{}}, args: args{[]byte("q")}, wantCommand: []byte("q")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{
				command: tt.fields.command,
			}
			if w.AddCommand(tt.args.b); !reflect.DeepEqual(tt.wantCommand, w.command) {
				t.Errorf("got: command=%s, want: command=%s", w.command, tt.wantCommand)
			}
		})
	}
}

func TestWindow_RemoveCommand(t *testing.T) {
	type fields struct {
		command []byte
	}
	tests := []struct {
		name        string
		fields      fields
		wantCommand []byte
	}{
		{name: "2 typed", fields: fields{command: []byte("wq")}, wantCommand: []byte("w")},
		{name: "1 typed", fields: fields{command: []byte("q")}, wantCommand: []byte{}},
		{name: "not typed", fields: fields{command: []byte{}}, wantCommand: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{
				command: tt.fields.command,
			}
			if w.RemoveCommand(); !reflect.DeepEqual(tt.wantCommand, w.command) {
				t.Errorf("got: command=%s, want: command=%s", w.command, tt.wantCommand)
			}
		})
	}
}

func TestWindow_ResetCommand(t *testing.T) {
	type fields struct {
		command []byte
	}
	tests := []struct {
		name        string
		fields      fields
		wantCommand []byte
	}{
		{name: "typed", fields: fields{command: []byte("wq")}, wantCommand: []byte{}},
		{name: "not typed", fields: fields{command: []byte{}}, wantCommand: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{
				command: tt.fields.command,
			}
			if w.ResetCommand(); !reflect.DeepEqual(tt.wantCommand, w.command) {
				t.Errorf("got: command=%s, want: command=%s", w.command, tt.wantCommand)
			}
		})
	}
}

func TestWindow_TypedCommand(t *testing.T) {
	type fields struct {
		command []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "typed", fields: fields{command: []byte("wq")}, want: "wq"},
		{name: "blank", fields: fields{command: []byte{}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Window{
				command: tt.fields.command,
			}
			if got := w.TypedCommand(); got != tt.want {
				t.Errorf("TypedCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow_InputtedUp(t *testing.T) {
	type fields struct {
		Size         Size
		FileContents [][]byte
		position     Position
		mode         int
	}
	tests := []struct {
		name    string
		fields  fields
		wantX   int
		wantY   int
		wantOut []byte
	}{
		{
			name: "Y=1",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 1, Y: 1},
				mode:         normalMode,
			},
			wantX:   1,
			wantY:   1,
			wantOut: []byte(""),
		},
		{
			name: "Upper character length is greater than current X",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 7, Y: 3},
				mode:         normalMode,
			},
			wantX:   7,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 7, Y: 2, Up    \033[2;7H"),
		},
		{
			name: "Upper character length is equal current X",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 8, Y: 3},
				mode:         normalMode,
			},
			wantX:   8,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 8, Y: 2, Up    \033[2;8H"),
		},
		{
			name: "Upper character length is less than current X (not zero), normal mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 10, Y: 3},
				mode:         normalMode,
			},
			wantX:   8,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 8, Y: 2, Up    \033[2;8H"),
		},
		{
			name: "Upper character length is less than current X (not zero), insert mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 10, Y: 3},
				mode:         insertMode,
			},
			wantX:   9,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 9, Y: 2, Up    \033[2;9H"),
		},
		{
			name: "Upper character length is less than current X (zero)",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte(""), []byte("This is a pen!")},
				position:     Position{X: 10, Y: 3},
				mode:         normalMode,
			},
			wantX:   1,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 1, Y: 2, Up    \033[2;1H"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			w := &Window{
				Size:         tt.fields.Size,
				Output:       out,
				FileContents: tt.fields.FileContents,
				position:     tt.fields.position,
				mode:         tt.fields.mode,
			}
			w.InputtedUp()
			if tt.wantX != w.position.X || tt.wantY != w.position.Y {
				t.Errorf("got: X= %d, Y=%d  want: X=%d, Y=%d", w.position.X, w.position.Y, tt.wantX, tt.wantY)
			}
			if out.String() != string(tt.wantOut) {
				t.Errorf("got: %v, want:  %v", out.Bytes(), tt.wantOut)
			}
		})
	}
}

func TestWindow_InputtedDown(t *testing.T) {
	type fields struct {
		Size         Size
		FileContents [][]byte
		position     Position
		mode         int
	}
	tests := []struct {
		name    string
		fields  fields
		wantX   int
		wantY   int
		wantOut []byte
	}{
		{
			name: "Y= File lines",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 1, Y: 2},
				mode:         normalMode,
			},
			wantX:   1,
			wantY:   2,
			wantOut: []byte(""),
		},
		{
			name: "Lower character length is greater than current X",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 7, Y: 1},
				mode:         normalMode,
			},
			wantX:   7,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 7, Y: 2, Down  \033[2;7H"),
		},
		{
			name: "Lower character length is equal current X",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 8, Y: 1},
				mode:         normalMode,
			},
			wantX:   8,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 8, Y: 2, Down  \033[2;8H"),
		},
		{
			name: "Lower character length is less than current X (not zero), normal mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 10, Y: 1},
				mode:         normalMode,
			},
			wantX:   8,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 8, Y: 2, Down  \033[2;8H"),
		},
		{
			name: "Lower character length is less than current X (not zero), insert mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 10, Y: 1},
				mode:         insertMode,
			},
			wantX:   9,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 9, Y: 2, Down  \033[2;9H"),
		},
		{
			name: "Upper character length is less than current X (zero)",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte(""), []byte("This is a pen!")},
				position:     Position{X: 10, Y: 1},
				mode:         normalMode,
			},
			wantX:   1,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 1, Y: 2, Down  \033[2;1H"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			w := &Window{
				Size:         tt.fields.Size,
				Output:       out,
				FileContents: tt.fields.FileContents,
				position:     tt.fields.position,
				mode:         tt.fields.mode,
			}
			w.InputtedDown()
			if tt.wantX != w.position.X || tt.wantY != w.position.Y {
				t.Errorf("got: X= %d, Y=%d  want: X=%d, Y=%d", w.position.X, w.position.Y, tt.wantX, tt.wantY)
			}
			if out.String() != string(tt.wantOut) {
				t.Errorf("got: %v, want:  %v", out.Bytes(), tt.wantOut)
			}
		})
	}
}

func TestWindow_InputtedLeft(t *testing.T) {
	type fields struct {
		Size         Size
		FileContents [][]byte
		position     Position
	}
	tests := []struct {
		name    string
		fields  fields
		wantX   int
		wantY   int
		wantOut []byte
	}{
		{
			name: "X=1",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob"), []byte("OK ?")},
				position:     Position{X: 1, Y: 3},
			},
			wantX:   1,
			wantY:   3,
			wantOut: []byte("\033[100;0H> X: 1, Y: 3, Left  \033[3;1H"),
		},
		{
			name: "X>1",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob"), []byte("OK ?")},
				position:     Position{X: 3, Y: 3},
			},
			wantX:   2,
			wantY:   3,
			wantOut: []byte("\033[100;0H> X: 2, Y: 3, Left  \033[3;2H"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			w := &Window{
				Size:         tt.fields.Size,
				Output:       out,
				FileContents: tt.fields.FileContents,
				position:     tt.fields.position,
			}
			w.InputtedLeft()
			if tt.wantX != w.position.X || tt.wantY != w.position.Y {
				t.Errorf("got: X= %d, Y=%d  want: X=%d, Y=%d", w.position.X, w.position.Y, tt.wantX, tt.wantY)
			}
			if out.String() != string(tt.wantOut) {
				t.Errorf("got: %v, want:  %v", out.Bytes(), tt.wantOut)
			}
		})
	}
}

func TestWindow_InputtedRight(t *testing.T) {
	type fields struct {
		Size         Size
		FileContents [][]byte
		position     Position
		mode         int
	}
	tests := []struct {
		name    string
		fields  fields
		wantX   int
		wantY   int
		wantOut []byte
	}{
		{
			name: "X=character length, normal mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 8, Y: 2},
				mode:         normalMode,
			},
			wantX:   8,
			wantY:   2,
			wantOut: []byte("\033[2;8H"),
		},
		{
			name: "X=character length, insert mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 8, Y: 2},
				mode:         insertMode,
			},
			wantX:   9,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 9, Y: 2, Right\033[2;9H"),
		},
		{
			name: "X<character length",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 3, Y: 2},
				mode:         normalMode,
			},
			wantX:   4,
			wantY:   2,
			wantOut: []byte("\033[100;0H> X: 4, Y: 2, Right\033[2;4H"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			w := &Window{
				Size:         tt.fields.Size,
				Output:       out,
				FileContents: tt.fields.FileContents,
				position:     tt.fields.position,
				mode:         tt.fields.mode,
			}
			w.InputtedRight()
			if tt.wantX != w.position.X || tt.wantY != w.position.Y {
				t.Errorf("got: X= %d, Y=%d  want: X=%d, Y=%d", w.position.X, w.position.Y, tt.wantX, tt.wantY)
			}
			if out.String() != string(tt.wantOut) {
				t.Errorf("got: %v, want:  %v", out.Bytes(), tt.wantOut)
			}
		})
	}
}

func TestWindow_InputtedOther(t *testing.T) {
	type fields struct {
		Size         Size
		FileContents [][]byte
		position     Position
		mode         int
	}
	tests := []struct {
		name     string
		fields   fields
		input    []byte
		wantX    int
		wantY    int
		wantOut  []byte
		wantMode int
	}{
		{
			name: "inputted i and not insert mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 3, Y: 2},
				mode:         normalMode,
			},
			input:    []byte("i"),
			wantX:    3,
			wantY:    2,
			wantOut:  []byte(""),
			wantMode: insertMode,
		},
		{
			name: "inputted i and insert mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 3, Y: 2},
				mode:         insertMode,
			},
			input:    []byte("i"),
			wantX:    4,
			wantY:    2,
			wantOut:  []byte("\033[2;0HI iam bob\033[2;4H"),
			wantMode: insertMode,
		},
		{
			name: "inputted : and not insert mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 3, Y: 2},
				mode:         normalMode,
			},
			input:    []byte(":"),
			wantX:    3,
			wantY:    2,
			wantOut:  []byte("\033[100;0H:"),
			wantMode: commandMode,
		},
		{
			name: "inputted : and insert mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 3, Y: 2},
				mode:         insertMode,
			},
			input:    []byte(":"),
			wantX:    4,
			wantY:    2,
			wantOut:  []byte("\033[2;0HI :am bob\033[2;4H"),
			wantMode: insertMode,
		},
		{
			name: "inputted not i and insert mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 3, Y: 2},
				mode:         insertMode,
			},
			input:    []byte("A"),
			wantX:    4,
			wantY:    2,
			wantOut:  []byte("\033[2;0HI Aam bob\033[2;4H"),
			wantMode: insertMode,
		},
		{
			name: "inputted not i and not insert mode",
			fields: fields{
				Size: Size{
					Row:    100,
					Column: 150,
				},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position:     Position{X: 3, Y: 2},
				mode:         normalMode,
			},
			input:    []byte("A"),
			wantX:    3,
			wantY:    2,
			wantOut:  []byte("\033[100;0H> X: 3, Y: 2, input: A     \033[2;3H"),
			wantMode: normalMode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			w := &Window{
				Size:         tt.fields.Size,
				Output:       out,
				FileContents: tt.fields.FileContents,
				position:     tt.fields.position,
				mode:         tt.fields.mode,
			}
			w.InputtedOther(tt.input)
			if tt.wantX != w.position.X || tt.wantY != w.position.Y {
				t.Errorf("got: X= %d, Y=%d  want: X=%d, Y=%d", w.position.X, w.position.Y, tt.wantX, tt.wantY)
			}
			if out.String() != string(tt.wantOut) {
				t.Errorf("got: %v, want:  %v", out.Bytes(), tt.wantOut)
			}
			if tt.wantMode != w.mode {
				t.Errorf("got: mode=%d, want: mode=%d", w.mode, tt.wantMode)
			}
		})
	}
}

func TestWindow_PrintFileContents(t *testing.T) {
	type fields struct {
		Size         Size
		Output       io.Writer
		FileContents [][]byte
		position     Position
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
				position: Position{
					X: 2,
					Y: 3,
				},
			},
			want: []byte("\033[H\033[2JHello World!\nI am bob\n\n\033[3;2H"),
		},
		{
			name: "file row + 1 == window row",
			fields: fields{
				Size:         Size{Row: 3, Column: 100},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position: Position{
					X: 1,
					Y: 2,
				},
			},
			want: []byte("\033[H\033[2JHello World!\nI am bob\n\033[2;1H"),
		},
		{
			name: "file row  == window row",
			fields: fields{
				Size:         Size{Row: 2, Column: 100},
				FileContents: [][]byte{[]byte("Hello World!"), []byte("I am bob")},
				position: Position{
					X: 3,
					Y: 2,
				},
			},
			want: []byte("\033[H\033[2JHello World!\n\033[2;3H"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			w := &Window{
				Size:         tt.fields.Size,
				Output:       out,
				FileContents: tt.fields.FileContents,
				position:     tt.fields.position,
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
