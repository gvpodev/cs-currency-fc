package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Execute applies all migrations in the migrations folder
func Execute() {
	db, err := sql.Open("sqlite3", "currency.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir("./migrations")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Erro ao ler diretório de migrações: %v\n", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			fmt.Printf("Aplicando migração: %s\n", file.Name())

			content, err := ioutil.ReadFile(filepath.Join("./migrations", file.Name()))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stdout, "Erro ao ler arquivo de migração: %v\n", err)
			}

			tx, err := db.Begin()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stdout, "Erro ao iniciar transação: %v\n", err)
			}

			if _, err := tx.Exec(string(content)); err != nil {
				tx.Rollback()
				_, _ = fmt.Fprintf(os.Stdout, "Erro ao executar migração: %v\n", err)
			}

			err = tx.Commit()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stdout, "Erro ao confirmar transação: %v\n", err)
			}
		}
	}
}
