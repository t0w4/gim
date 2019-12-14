package position

type Position struct {
	X int
	Y int
}

func (p *Position) MoveDown(num int) {
	p.Y += num
}

func (p *Position) MoveUp(num int) {
	if p.Y == 0 {
		return
	}
	p.Y -= num
}

func (p *Position) MoveRight(num int) {
	p.X += num
}

func (p *Position) MoveLeft(num int) {
	if p.X == 0 {
		return
	}
	p.X -= num
}
