package extra_data

import (
	"encoding/csv"
	"errors"
	"fmt"
	sqlite "gateway/sqlite3"
	"io"
	"os"
	"strconv"
)

// csv file format should be:
// addr,private_key
// ...
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

		fmt.Println(fmt.Sprintf("%s : %b : %s", record, err == nil, err)) // record has the type []string
		if err == nil {
			success_count += 1
		} else {
			fail_count += 1
			fmt.Println(err)
		}
	}

	fmt.Println("import complete")
	fmt.Println("success count: " + fmt.Sprintf("%d", success_count))
	fmt.Println("fail count: " + fmt.Sprintf("%d", fail_count))
	fmt.Println("total count: " + fmt.Sprintf("%d", (success_count+fail_count)))
}

func insert_row(record []string) error {
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
