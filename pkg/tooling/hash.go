package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/spf13/viper"
)

func HashPassword(pwd string) (string, error) {
	var PasswordSalt = viper.GetString("tools.PasswordSalt")

	hash := sha256.New()
	_, err := hash.Write([]byte(pwd))
	if err != nil {
		return "", err
	}
	PwdHash := hex.EncodeToString(hash.Sum([]byte(PasswordSalt)))
	return PwdHash, nil
}
