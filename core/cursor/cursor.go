package cursor

type Cursor struct {
	Ascending  bool
	Limit      int
	Cursor     int64
	NextCursor int64
	HasMore    bool
}

func New(ascending bool, limit int, cursor int64) Cursor {
	if limit <= 0 || limit > 100 {
		limit = 5
	}
	return Cursor{
		Ascending:  ascending,
		Limit:      limit,
		Cursor:     cursor,
		NextCursor: cursor,
		HasMore:    false,
	}
}

func NewOldest(limit int) Cursor {
	return New(false, limit, 0)
}

func (c *Cursor) SetHasMore(hasMore bool) {
	c.HasMore = hasMore
}

func (c *Cursor) SetNextCursor(nextCursor int64) {
	c.NextCursor = nextCursor
}
