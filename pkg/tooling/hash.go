package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/spf13/viper"
	"log"
)

func HashPassword(pwd string) (string, error) {
	var PasswordSalt = viper.GetString("tools.PasswordSalt")

	log.Println("salt:", PasswordSalt)

	hash := sha256.New()
	n, err := hash.Write([]byte(pwd))
	if err != nil {
		log.Println(n, err.Error())
		return "", err
	}
	PwdHash := hex.EncodeToString(hash.Sum([]byte(PasswordSalt)))
	return PwdHash, nil
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
