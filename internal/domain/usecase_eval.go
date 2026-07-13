package domain

type Customer struct {
	ID      string
	Country string
	Tier    string
}

type EvaluateUsecaseInput struct {
	ProjectID int64
	Customer  Customer
	Cart      Cart
}
