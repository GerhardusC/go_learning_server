package dbInteractions

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing-server/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

var redisContext = context.Background()
const SESSION_DURATION = time.Second * 90

func (user UserPreAuth) AuthenticateUsernamePwd () (string, error) {
	authorisedUser, err := user.GetFromDB()
	if err != nil {
		return "", err
	}

	return authorisedUser.GenerateToken()
}

/** 
**				-> saveTemporaryUser 
**			       /		       ) <- SessionID finally returned
** Saving the OTP: (   SendOTP --> saveOTP
** Returns the session ID after sending the OTP. */
func (user *UserPreAuth) SendOTP () (string, error) {
	otp := strings.TrimPrefix(fmt.Sprintf("%d", rand.Intn(1000000) + 1000000), "1")

	go utils.SendOTPEmail(otp, user.Email)
	
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})
	defer rdb.Close()

	sessionID := fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String() + user.Username + user.Email)))

	user.saveOTP(otp, rdb)
	user.saveTemporaryUser(sessionID, rdb)


	return sessionID, nil
}

/** 
** Returns the session ID used to retrieve the user if user is verified. */
func (user *UserPreAuth) saveTemporaryUser (sessionID string, rdb *redis.Client) error {
	hashedPwd := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Username + user.UnhashedPwd)))

	userToSave := UserWithHashedPwd{
		Email: user.Email,
		HashedPwd: hashedPwd,
		Username: user.Username,
	}

	userBytes, err := json.Marshal(&userToSave)

	if err != nil {
		fmt.Println("Failed unmarshal user")
		return err
	}

	err = rdb.Set(
		redisContext,
		sessionID,
		userBytes,
		SESSION_DURATION,
	).Err()

	if err != nil {
		fmt.Println("Failed to set session ID in Redis")
		return err
	}

	return nil
}

func (user *UserPreAuth) saveOTP (otp string, rdb *redis.Client) error {
	hashedOTP := fmt.Sprintf("%x", sha256.Sum256([]byte(otp)))
	redisKey := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Username + user.Email)))

	err := rdb.Set(redisContext, redisKey, hashedOTP, SESSION_DURATION).Err()

	if err != nil {
		log.Println("Something went wrong while saving the OTP", err)
		return err
	}
	return err
}

func (user UserWithHashedPwd) verifyOTP (otp string, rdb *redis.Client) error {
	if user.Username == "" || otp == "" {
		return errors.New("Username or otp can't be empty")
	}

	hashedOTP := fmt.Sprintf("%x", sha256.Sum256([]byte(otp)))
	hashedUsername := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Username + user.Email)))

	retrievedOTP, err := rdb.Get(redisContext, hashedUsername).Result()

	if err != nil {
		log.Println("Something went wrong while getting the OTP from redis", err)
		return err
	}

	if retrievedOTP != hashedOTP {
		return errors.New("Failed to verify OTP")
	}

	return nil
}

func (authorisedUser *User) GenerateToken()  (string, error) {
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "test-secret"
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_details": *authorisedUser,
			"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	)
	
	return token.SignedString([]byte(sec))
}

func (session OTPVerifyObj) GetUser() (*UserWithHashedPwd, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})
	defer rdb.Close()

	userString, err := rdb.Get(redisContext, session.SessionID).Bytes()

	if err != nil {
		fmt.Println("Cannot get user from context")
		return nil, err
	}

	var user UserWithHashedPwd

	err = json.Unmarshal(userString, &user)

	if err != nil {
		fmt.Println("Failed to unmarshal user retrieved from redis")
		return nil, err
	}

	err = user.verifyOTP(session.OTP, rdb)

	if err != nil {
		fmt.Println("Cannot verify OTP")
		return nil, err
	}

	return &user, nil
}

