package secrets

import (
	"database/sql"
	"fmt"

	"vault-cli/internal/config"
)

func Rotate(database *sql.DB, cfg *config.Config) (int, error) {
	rows, err := database.Query(`SELECT category, name FROM secrets`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var c, n string
		if err := rows.Scan(&c, &n); err != nil {
			return count, err
		}

		plain, err := Get(database, cfg, c, n)
		if err != nil {
			return count, fmt.Errorf("get secret %s/%s: %w", c, n, err)
		}

		if err := Add(database, cfg, Secret{Category: c, Name: n, Value: plain}); err != nil {
			return count, fmt.Errorf("re-encrypt %s/%s: %w", c, n, err)
		}

		count++
	}

	return count, rows.Err()
}
