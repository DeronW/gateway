package db

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/delongw/go-int-cipher"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

//var db *driver.Conn
var db *sql.DB
var once sync.Once

func GetPrivateKey(addr int) (key string, err error) {
	stmt, _ := db.Prepare("select encrypted_private_key from teleport where addr = ?")
	defer stmt.Close()

	var encrypted_private_key string
	encrypted_key := int_cipher.Encrypt(addr, RC4_KEY)
	err = stmt.QueryRow(encrypted_key).Scan(&encrypted_private_key)
	if err != nil {
		log.WithFields(log.Fields{
			"error":         err,
			"encrypted_key": encrypted_key,
		}).Info("no find encrypted teleport private key")
	}
	return encrypted_private_key, nil
}

func retry_connect() {
	if db != nil {
		db.Close()
	}
	sqlite_db, err := sql.Open(
		"sqlite3",
		"/home/delong/phantom-go/src/gateway/sqlite3.db",
	)
	if err != nil {
		log.Fatal(err)
	}
	db = sqlite_db
}

func init() {
	once.Do(retry_connect)
}
