package postgres

import (
	"fmt"
	"os"
	"testing"
)

var testDbName = "meirltest"

func TestMain(m *testing.M) {
	result, err := RunWithTestDB(testDbName, false, func() int {
		return m.Run()
	}, "../../resources/sql/create_schema.sql")
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(result)
}
