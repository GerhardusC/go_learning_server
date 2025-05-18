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
	OTP		string	`json:"otp"`
	SessionID	string	`json:"session_id"`
}

type UserWithHashedPwd struct {
	Email		string	`json:"email"`
	HashedPwd	string	`json:"password"`
	Username	string	`json:"username"`
}
