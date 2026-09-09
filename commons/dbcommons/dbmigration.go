package dbcommons

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"sort"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
	utils "github.com/juggleim/imserver-console/commons/tools"
	"gorm.io/gorm"
)

//go:embed sqls/*
var sqlFs embed.FS

const (
	JChatDbVersionKey = "jchatdb_version"
)

type CountResult struct {
	Count int64 `gorm:"count"`
}

func Upgrade() error {
	// upgrade jchat db
	var currVersion int64 = 0
	dao := GlobalConfDao{}
	conf, err := dao.FindByKey(JChatDbVersionKey)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		var mysqlErr *mysqlDriver.MySQLError
		// Only the initial version-table lookup may bootstrap a missing table.
		if !errors.As(err, &mysqlErr) || mysqlErr.Number != 1146 {
			return fmt.Errorf("read database migration version: %w", err)
		}
	}
	if err == nil && conf != nil {
		ver, err := utils.String2Int64(conf.ConfValue)
		if err != nil || ver < 0 {
			return fmt.Errorf("invalid database migration version")
		}
		currVersion = ver
	}
	fmt.Println("[JChatDbMigration]current version:", currVersion)
	sqlFiles, err := sqlFs.ReadDir("sqls")
	if err != nil {
		return fmt.Errorf("read database migration files: %w", err)
	}
	neededVers := []int64{}
	for _, sqlFile := range sqlFiles {
		fileName := sqlFile.Name()
		if len(fileName) == 12 {
			fileName = fileName[:8]
		}
		ver, err := utils.String2Int64(fileName)
		if err == nil && ver > 0 {
			neededVers = append(neededVers, ver)
		}
	}
	//sort
	sort.Slice(neededVers, func(i, j int) bool {
		return neededVers[i] < neededVers[j]
	})
	for _, ver := range neededVers {
		if ver > currVersion {
			sqlFileName := fmt.Sprintf("sqls/%d.sql", ver)
			fmt.Println("[DbMigration]start to execute sql file:", sqlFileName)
			if err := executeSqlFile(sqlFileName); err != nil {
				return fmt.Errorf("execute database migration %d: %w", ver, err)
			}
			if err := dao.Upsert(GlobalConfDao{
				ConfKey:   JChatDbVersionKey,
				ConfValue: fmt.Sprintf("%d", ver),
			}); err != nil {
				return fmt.Errorf("write database migration version %d: %w", ver, err)
			}
			fmt.Println("[DbMigration]execute sql file success:", sqlFileName)
		}
	}
	return nil
}

func executeSqlFile(fileName string) error {
	sqlFile, err := sqlFs.Open(fileName)
	if err != nil {
		fmt.Println("[DbMigration_Err]Read sql file err:", err, "file_name:", fileName)
		return err
	}
	defer sqlFile.Close()

	// MySQL user variables and PREPARE statements are connection-local.
	return GetDb().Connection(func(conn *gorm.DB) error {
		scanner := bufio.NewScanner(sqlFile)
		var queryBuilder strings.Builder
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "--") {
				continue
			}
			queryBuilder.WriteString(line)
			queryBuilder.WriteByte(' ')
			if strings.HasSuffix(line, ";") {
				query := strings.TrimSpace(queryBuilder.String())
				if query != "" {
					if err := conn.Exec(query).Error; err != nil {
						fmt.Println("[DbMigration_Err]Execute sql error:", err, query)
						return err
					}
				}
				queryBuilder.Reset()
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("[DbMigration_Err]Scan sql file err:", err, "file_name:", fileName)
			return err
		}
		if strings.TrimSpace(queryBuilder.String()) != "" {
			return fmt.Errorf("sql file contains an unterminated statement: %s", fileName)
		}
		return nil
	})
}
