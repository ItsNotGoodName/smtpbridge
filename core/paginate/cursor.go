package paginate

import "math"

type Cursor struct {
	Cursor     int64
	Limit      int
	Ascending  bool
	BackCursor int64
	NextCursor int64
}

func NewCursor(cursor int64, limit int, ascending bool) Cursor {
	if limit <= 0 || limit > 100 {
		limit = 5
	}
	if !ascending && cursor == 0 {
		cursor = math.MaxInt64
	}
	return Cursor{
		Cursor:     cursor,
		Limit:      limit,
		Ascending:  ascending,
		NextCursor: cursor,
	}
}

func NewCursorOldest(limit int) Cursor {
	return NewCursor(0, limit, false)
}

func (c *Cursor) SetNextCursor(nextCursor int64) {
	c.NextCursor = nextCursor
}

func (c *Cursor) SetBackCursor(backCursor int64) {
	c.BackCursor = backCursor
}
