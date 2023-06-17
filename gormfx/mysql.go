package gormfx

import (
	"github.com/guoyk93/ufx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLParams struct {
	DSN string `json:"dsn" default:"root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local" validate:"required"`
}

func DecodeMySQLParams(conf ufx.Conf) (params MySQLParams, err error) {
	err = conf.Bind(&params, "mysql")
	return
}

func NewMySQLDialector(params MySQLParams) gorm.Dialector {
	return mysql.Open(params.DSN)
}
