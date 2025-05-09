package dbInteractions

type DBRowMeasurement[T string | float64] struct {
	Timestamp int
	Topic string
	Value T
}

type User struct {
	Username	string
	HashedPwd	string
	CreatedAt	string
	PermissionLevel int
}
