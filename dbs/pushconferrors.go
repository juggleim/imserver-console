package dbs

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrPushConfConflict = errors.New("push configuration already exists")
	ErrPushConfNotFound = errors.New("push configuration not found")
)

func normalizePushConfError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrPushConfNotFound
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrPushConfConflict
	}
	return err
}
