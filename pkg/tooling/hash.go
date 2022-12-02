package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/spf13/viper"
	"log"
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

func Hash(access string, refresh string) string {
	hash := sha1.New()
	_, err := hash.Write([]byte(access + refresh))
	if err != nil {
		log.Println(err.Error())
	}
	Symbols := hex.EncodeToString(hash.Sum([]byte("22231313")))
	return string(Symbols[:26])

}
