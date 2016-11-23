package postgres

import (
	"os"
	"os/exec"

	"github.com/jmoiron/sqlx"
)

const (
	exitOk int = iota
	exitFail
)

// RunWithTestDB runs a test function, ensuring that
// a test database with the given name exists and that the
// database is torn down after testFn finishes running. Errors
// are returned along with the result of the test function. Paths to
// SQL scripts may optionally be passed in to be executed on the test
// database instance before testFn is run.
// *NOTE*: The script parsing function is VERY simple, it just delineates based on ';'.
// Keep that in mind when passing in scripts to execute
func RunWithTestDB(dbName string, verbose bool, testFn func() int, scripts ...string) (result int, err error) {
	result = exitFail

	var db *sqlx.DB
	db, err = InitDB("postgres", "", "localhost", "5432", "")
	if err != nil {
		return
	}
	defer db.Close()

	// drop now in case panic caused defer to not run
	db.Exec("DROP DATABASE " + dbName)

	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		return
	}

	// delete database at end
	defer db.Exec("DROP DATABASE " + dbName)

	for _, script := range scripts {
		err = execSimpleScript(dbName, script, verbose)
		if err != nil {
			return
		}
	}

	result = testFn()
	return
}

func whileConnectedToTestDb(dbName string, do func(*sqlx.DB) error) error {
	db, err := InitDB("postgres", "", "localhost", "5432", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	return do(db)
}

func execSimpleScript(dbName, path string, verbose bool) error {
	cmd := exec.Command("psql", "-U", "postgres", "-d", dbName, "-f", path)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
