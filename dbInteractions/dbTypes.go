package dbInteractions

type DBRowMeasurement[T string | float64] struct {
	Timestamp int
	Topic string
	Value T
}

type User struct {
	ID		int	`json:"ID"`
	CreatedAt	string	`json:"created_at"`
	Email		string	`json:"email"`
	PermissionLevel int	`json:"permission_level"`
	Username	string	`json:"username"`
}

type UserPreAuth struct {
	Email		string	`json:"email,omitempty"`
	UnhashedPwd	string	`json:"password,omitempty"`
	Username	string	`json:"username"`
}

type OTPVerifyObj struct {
	User	UserPreAuth	`json:"user"`
	OTP	string		`json:"otp"`
}
