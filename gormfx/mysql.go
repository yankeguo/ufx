package gormfx

import (
	"flag"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLParams struct {
	DSN string
}

func DecodeMySQLParams(fset *flag.FlagSet) *MySQLParams {
	p := &MySQLParams{}
	fset.StringVar(&p.DSN, "mysql.dsn", "root:root@tcp(127.0.0.3:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", "mysql dsn")
	return p
}

func NewMySQLDialector(p *MySQLParams) gorm.Dialector {
	return mysql.Open(p.DSN)
}
