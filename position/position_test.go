package position

import "testing"

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
