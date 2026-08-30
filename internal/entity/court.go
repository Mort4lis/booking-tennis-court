package entity

// Court is a physical tennis court that can be booked.
type Court struct {
	ID     int64
	Name   string
	Indoor bool
	Active bool
}
