package db

import (
	"database/sql"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/delongw/go-int-cipher"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

var db *sql.DB
var once sync.Once
var DB_PATH string = ""

func GetPrivateKey(addr int) (key []byte, ok bool, err error) {
	stmt, _ := db.Prepare("SELECT encrypted_private_key FROM teleport WHERE encrypted_addr = ?")
	defer stmt.Close()

	var encrypted_private_key string
	encrypted_addr := int_cipher.Encrypt(uint(addr), RC4_KEY)
	err = stmt.QueryRow(encrypted_addr).Scan(&encrypted_private_key)
	if err != nil {
		return
	}

	key, err = decrypt_private_key(encrypted_private_key)
	if err != nil {
		return
	}
	return key, true, nil
}

func SetPrivateKey(addr int, private_key string) error {
	if len(private_key) != 32 {
		return errors.New("private length is wrong: " + private_key)
	}
	key, ok, _ := GetPrivateKey(addr)
	if ok && fmt.Sprintf("%X", key) == private_key {
		return nil
	}

	stmt, _ := db.Prepare("INSERT INTO teleport (encrypted_addr, encrypted_private_key) VALUES (?, ?)")
	defer stmt.Close()

	encrypted_addr := int_cipher.Encrypt(uint(addr), RC4_KEY)
	encrypted_key := encrypt_private_key(private_key)
	_, err := stmt.Exec(encrypted_addr, encrypted_key)
	return err
}

func retry_connect() {
	if db != nil {
		db.Close()
	}
	sqlite_db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	db = sqlite_db
}

// this is initial func
func SetSqlite3Path(path string) {
	if DB_PATH != "" {
		panic("db path can not be set twice")
	}
	DB_PATH = path
	once.Do(retry_connect)
}
