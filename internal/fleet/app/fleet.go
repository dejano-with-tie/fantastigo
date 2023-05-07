package app

type Fleet struct {
	Id       string
	Name     string
	Capacity int
}

type FleetRepo interface {
	GetById(id string) (Fleet, error)
	Save(f Fleet) (Fleet, error)
}
