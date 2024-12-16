package defaults

const (
	Rows       = 1000
	Cols       = 1000
	Iterations = 512
)

var (
	Directions = [...]struct{ DX, DY int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	Threads = [...]int{2, 4, 8, 16}
)
