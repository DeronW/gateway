package db

import (
	"database/sql"
	"errors"
	"fmt"
	"gateway/lib/misc"
	log "github.com/Sirupsen/logrus"
	"github.com/delongw/go-int-cipher"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

//var db *driver.Conn
var db *sql.DB
var once sync.Once

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

func encrypt_private_key(private_key string) string {
	rand_key := misc.Rand8byte()
	key_bytes, _ := misc.Str2byte(private_key)
	secret := misc.Rc4xor(key_bytes, misc.BytesXor(rand_key, []byte(RC4_KEY)))
	return fmt.Sprintf("%X", append(rand_key, secret...))
}

func decrypt_private_key(ekey string) ([]byte, error) {
	if len(ekey) != 48 {
		return nil, errors.New("illegal encrypted private key: " + ekey)
	}
	rand_key, _ := misc.Str2byte(ekey[:16])
	secret, _ := misc.Str2byte(ekey[16:])
	return misc.Rc4xor(secret, misc.BytesXor(rand_key, []byte(RC4_KEY))), nil
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
