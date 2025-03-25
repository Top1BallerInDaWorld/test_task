package service

type Clicker interface {
	AddClick(int)
}
type ClickCounter struct {
	clicker Clicker
}

func NewClickCounter(clicker Clicker) *ClickCounter {
	return &ClickCounter{clicker: clicker}
}

func (c *ClickCounter) AddClick(i int) {
	c.clicker.AddClick(i)
}
