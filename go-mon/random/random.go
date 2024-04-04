package random

import (
	crand "crypto/rand"
	"errors"
	"math"
	"math/big"
	"math/rand"
)

var (
	ErrRandomTooMuch = errors.New("[ERROR]: The length of randomness is too much")
)

// StringRand generates a random string of fixed size
func StringRand(size int) (string, error) {
	if size > 32 {
		return "", ErrRandomTooMuch
	}
	alpha := "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(buf), nil
}

// NumberRand generates a random number of fixed size
func NumberRand(size int) (int, error) {
	if size > 19 {
		return 0, ErrRandomTooMuch
	}
	maxLimit := int64(int(math.Pow10(size)) - 1)
	lowLimit := int(math.Pow10(size - 1))
	randNumb, err := crand.Int(crand.Reader, big.NewInt(maxLimit))
	if err != nil {
		return 0, err
	}
	randNumbInt := int(randNumb.Int64())
	// Handling integers between 0, 10^(n-1) .. for n=4, handling cases between (0, 999)
	if randNumbInt <= lowLimit {
		randNumbInt += lowLimit
	}
	// Never likely to occur, must for safe side.
	if randNumbInt > int(maxLimit) {
		randNumbInt = int(maxLimit)
	}
	return randNumbInt, nil
}
