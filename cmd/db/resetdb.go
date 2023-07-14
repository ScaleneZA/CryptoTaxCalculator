package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ResetDB(dbc *sql.DB) {
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	projectDir := strings.TrimSpace(string(output))
	log.Println("Project directory:", projectDir)

	sqlBytes, err := os.ReadFile("cmd/db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	sqlQuery := string(sqlBytes)

	_, err = dbc.Exec(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Schema reset and all data deleted.")
}
