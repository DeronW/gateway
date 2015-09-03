package teleport

import (
//"encoding/base64"
)

type Command struct {
	encrypted      bool
	w_wncrypted    bool
	addr           int
	op             int
	params         string
	user_key_index string
	src_cost       string
	src_seq        int
	version        int
}
