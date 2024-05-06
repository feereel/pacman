package utility

type Number interface {
	int | uint | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Vector2D[T Number] struct {
	X T
	Y T
}

func (vect Vector2D[T]) InBound(minX, maxX, minY, maxY T) bool {
	return vect.X >= minX && vect.X < maxX && vect.Y >= minY && vect.Y < maxY
}

func (A Vector2D[T]) Add(B Vector2D[T]) Vector2D[T] {
	return Vector2D[T]{
		X: A.X + B.X,
		Y: A.Y + B.Y,
	}
}

func (A Vector2D[T]) Sub(B Vector2D[T]) Vector2D[T] {
	return Vector2D[T]{
		X: A.X - B.X,
		Y: A.Y - B.Y,
	}
}

func ReverseBytes(s []byte) (o []byte) {
	o = make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		o[i] = s[len(s)-i-1]
	}
	return o
}
