package repository

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (Mock) Close() error {
	return nil
}
