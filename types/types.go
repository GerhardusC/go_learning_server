package types

type Person struct {
	Name string		`json:"name"`
	Surname string		`json:"surname"`
	Siblings []string	`json:"siblings"`
}

type Measurement struct {
	Timestamp int		`json:"timestamp"`
	Topic string            `json:"topic"`
	Value float64           `json:"value"`
}
