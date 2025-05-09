package dbInteractions

type DBRowMeasurement[T string | float64] struct {
	Timestamp int
	Topic string
	Value T
}

type User struct {
	ID		int
	CreatedAt	string
	Email		string
	PermissionLevel int
	Username	string
}

type UserPreAuth struct {
	Email		string	`json:"email"`
	UnhashedPwd	string	`json:"password"`
	Username	string	`json:"username"`
}
