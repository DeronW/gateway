package extra_data

import (
	"encoding/csv"
	"errors"
	"fmt"
	sqlite "gateway/db"
	"io"
	"os"
	"strconv"
)

func ImportTeleportData(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	success_count, fail_count := 0, 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = insert_row(record)
		if err == nil {
			fail_count += 1
		} else {
			success_count += 1
			fmt.Println(err)
		}
	}

	fmt.Println("import complete")
	fmt.Println("success count: " + fmt.Sprintf("%d", success_count))
	fmt.Println("fail count: " + fmt.Sprintf("%d", fail_count))
	fmt.Println("total count: " + fmt.Sprintf("%d", (success_count+fail_count)))
}

func insert_row(record []string) error {
	fmt.Println(record) // record has the type []string
	if len(record) != 2 {
		return errors.New("record format error")
	}

	addr, err := strconv.ParseInt(record[0], 10, 32)
	if err != nil {
		return err
	}
	//private_key, err := rails_private_key_to_bytes(record[1])
	private_key := record[1]

	if err != nil {
		return err
	}
	err = sqlite.SetPrivateKey(int(addr), private_key)

	if err != nil {
		return err
	}
	return nil
}

//func rails_private_key_to_bytes(s string) (string, error) {
//if len(s) != 32 {
//return "", errors.New("wrong private, lenght is not 32: " + s)
//}
//var b []byte
//for i := 0; i < 32; i += 2 {
//a, err := strconv.ParseUint(string(s[i:i+2]), 16, 8)
//if err != nil {
//return "", err
//}
//b = append(b, byte(a))
//}
//return fmt.Sprintf("%X", b), nil
//}
