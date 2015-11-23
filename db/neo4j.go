package db

/***

import (
	"fmt"
	"github.com/jmcvetta/neoism"
	"sync"
)

var neo4j *neoism.Database
var once sync.Once

func GetPrivateKeyNeo4j(addr int) (key string, err error) {
	//return "55294d59b1f1db94f848fd2364ebc979", nil
	//encrypted_private_key := "55294d59b1f1db94f848fd2364ebc979"
	//return encrypted_private_key, nil
	res := []struct {
		key string `json:t.k`
	}{}
	cq := &neoism.CypherQuery{
		Statement: `
			MATCH (t)
			WHERE t.addr = {addr}
			RETURN t.k
		`,
		Parameters: neoism.Props{"addr": 8444},
		Result:     &res,
	}
	err = neo4j.Cypher(cq)
	return res[0].key, nil
}

func init() {
	once.Do(func() {
		db, err := neoism.Connect("http://localhost:7474")
		if err != nil {
			fmt.Println("Can not connect to neo4j")
			return
		}
		neo4j = db
	})
}

***/
