package dbInteractions

import (
	"context"
	"crypto/sha256"
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

func (user UserPreAuth) AuthenticateUsernamePwd () (string, error) {
	authorisedUser, err := user.GetFromDB()
	if err != nil {
		return "", err
	}

	return GenerateToken(&authorisedUser)
}

func (user UserPreAuth) SendOTP () error {
	// TODO: Create session token flow.
	otp := strings.TrimPrefix(fmt.Sprintf("%d", rand.Intn(1000000)), "1")
	err := user.saveOTP(otp)

	if err != nil {
		log.Println("Could not save OTP for the session, aborting")
		return err
	}

	err = utils.SendOTP(otp, user.Email)
	
	if err != nil {
		log.Println("Could not send OTP for the session.")
		// TODO: CLear the current OTP in the redis instance.
		return err
	}
	return nil
}

func (user UserPreAuth) VerifyOTP (otp string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})
	defer rdb.Close()
	
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

func (user UserPreAuth) saveOTP (otp string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})
	defer rdb.Close()

	hashedOTP := fmt.Sprintf("%x", sha256.Sum256([]byte(otp)))
	hashedUsername := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Username + user.Email)))

	err := rdb.Set(redisContext, hashedUsername, hashedOTP, time.Second * 60).Err()

	if err != nil {
		log.Println("Something went wrong while setting the OTP in redis", err)
		return err
	}

	// TODO: Currently pwd can change halfway through auth, fix this.

	err = rdb.Set(redisContext, hashedUsername, hashedOTP, time.Second * 80).Err()

	return nil
}

func GenerateToken (authorisedUser *User) (string, error) {
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "test-secret"
	}

	// TODO: Add expiery date for token.
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_details": *authorisedUser,
		},
	)
	
	return token.SignedString([]byte(sec))
}
