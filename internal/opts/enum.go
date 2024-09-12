package opts

import "fmt"

type SQLEngine string

const (
	SQLEnginePostgresql SQLEngine = "postgresql"
	SQLEngineMySQL      SQLEngine = "mysql"
	SQLEngineSQLite     SQLEngine = "sqlite"
)

type SQLDriver string

const (
	SQLPackagePGXV4    string = "pgx/v4"
	SQLPackagePGXV5    string = "pgx/v5"
	SQLPackageStandard string = "database/sql"
)

var validPackages = map[string]struct{}{
	string(SQLPackagePGXV4):    {},
	string(SQLPackagePGXV5):    {},
	string(SQLPackageStandard): {},
}

func validatePackage(sqlPackage string) error {
	if _, found := validPackages[sqlPackage]; !found {
		return fmt.Errorf("unknown SQL package: %s", sqlPackage)
	}
	return nil
}

const (
	SQLDriverPGXV4            SQLDriver = "github.com/jackc/pgx/v4"
	SQLDriverPGXV5            SQLDriver = "github.com/jackc/pgx/v5"
	SQLDriverLibPQ            SQLDriver = "github.com/lib/pq"
	SQLDriverGoSQLDriverMySQL SQLDriver = "github.com/go-sql-driver/mysql"
)

func NewSQLDriver(sqlPackage string) SQLDriver {
	switch sqlPackage {
	case SQLPackagePGXV4:
		return SQLDriverPGXV4
	case SQLPackagePGXV5:
		return SQLDriverPGXV5
	default:
		return SQLDriverLibPQ
	}
}

var validDrivers = map[string]struct{}{
	string(SQLDriverPGXV4):            {},
	string(SQLDriverPGXV5):            {},
	string(SQLDriverLibPQ):            {},
	string(SQLDriverGoSQLDriverMySQL): {},
}

func validateDriver(sqlDriver string) error {
	if _, found := validDrivers[sqlDriver]; !found {
		return fmt.Errorf("unknown SQL driver: %s", sqlDriver)
	}
	return nil
}

func (d SQLDriver) IsPGX() bool {
	return d == SQLDriverPGXV4 || d == SQLDriverPGXV5
}

func (d SQLDriver) IsLibPQ() bool {
	return d == SQLDriverLibPQ
}

func (d SQLDriver) IsGoSQLDriverMySQL() bool {
	return d == SQLDriverGoSQLDriverMySQL
}

func (d SQLDriver) Package() string {
	switch d {
	case SQLDriverPGXV4:
		return SQLPackagePGXV4
	case SQLDriverPGXV5:
		return SQLPackagePGXV5
	default:
		return SQLPackageStandard
	}
}
