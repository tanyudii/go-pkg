package bcrypt

import "golang.org/x/crypto/bcrypt"

func Hash(val string, c ...int) (string, error) {
	cost := bcrypt.DefaultCost
	if len(c) > 0 {
		cost = c[0]
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(val), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Compare(hashed string, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return false, err
	}
	return true, nil
}
