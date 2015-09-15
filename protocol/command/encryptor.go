package command

func Decrypt(s []byte) (cnt []byte) {
	ext := len(s) % 16
	blocks := len(s) / 16
	//plain := ""

	for i := 0; i < blocks; i++ {
		//plain += s[i*16:(i+1)*16] ^ iv_str
	}

	if ext != 0 {

	}
	return
}

func decrypt_iv_str() {

}
