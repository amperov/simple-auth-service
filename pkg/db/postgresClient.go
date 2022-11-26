package db

// package db implements some atomic methods for working with DB

import (
	hash "authService/pkg/tooling"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"log"
)

type Client struct {
	conn *pgx.Conn
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func InitDBConfig() *Config {
	cfg := Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
	return &cfg
}

func GetClient() (*Client, error) {
	dbCfg := InitDBConfig()
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.DBName, dbCfg.Password, dbCfg.SSLMode))
	if err != nil {
		log.Printf("[WARN] from GetClient: %s\n", err.Error())
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Printf("[WARN] from GetClient: %s\n", err.Error())
		return nil, err
	}
	return &Client{conn: conn}, nil

}

//TODO Update Method

func (c *Client) IsExist(Username, Email, Password string, ctx context.Context) (bool, error) {
	db := c.conn
	var exist bool

	query := fmt.Sprintf("SELECT exists(id) FROM users WHERE username=$1 AND password_hash=$2")
	PasswordHash, err := hash.HashPassword(Password)
	if err != nil {
		return exist, err
	}
	row := db.QueryRow(ctx, query, Username, PasswordHash)
	err = row.Scan(&exist)
	if err != nil {
		return exist, err
	}
	if exist != true {
		query = fmt.Sprintf("SELECT exists(id) FROM users WHERE email=$1 AND password_hash=$2")

		row = db.QueryRow(ctx, query, Email, PasswordHash)
		err = row.Scan(&exist)
		if err != nil {
			return exist, err
		}
	}
	return exist, nil
}
func (c *Client) Insert(Username, Email, Password string, ctx context.Context) (int, error) {
	db := c.conn
	var id int

	query := fmt.Sprintf("INSERT INTO users ( username, email, password_hash) values ($1, $2, $3) returning id")
	PasswordHash, err := hash.HashPassword(Password)
	if err != nil {
		return 0, err
	}
	row := db.QueryRow(ctx, query, Username, Email, PasswordHash)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (c *Client) Get(Username, Email, Password string, ctx context.Context) (int, error) {
	db := c.conn
	var id int

	query := fmt.Sprintf("SELECT id FROM users WHERE (email=$1 AND password_hash=$2) OR (username=$3 AND password_hash=$2)")
	PasswordHash, err := hash.HashPassword(Password)
	if err != nil {
		return 0, err
	}
	row := db.QueryRow(ctx, query, Username, PasswordHash, Email)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *Client) Delete(Username, Email, Password string, ctx context.Context) error {
	db := c.conn

	query := fmt.Sprintf("DELETE FROM users WHERE (email=$1 AND password_hash=$2) OR (username=$3 AND password_hash=$2)")
	PasswordHash, err := hash.HashPassword(Password)
	if err != nil {
		return err
	}
	db.QueryRow(ctx, query, Username, PasswordHash, Email)

	return nil
}