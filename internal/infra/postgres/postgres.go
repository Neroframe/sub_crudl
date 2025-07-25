package postgres

import (
	"fmt"

	"github.com/Neroframe/sub_crudl/config"
)

func BuildDSN(sql config.Postgres) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		sql.Host, sql.Port, sql.User, sql.Password, sql.DBName)
}
