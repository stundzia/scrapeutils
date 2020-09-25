package db

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type ProxyDB struct {
	dbConn *sql.DB
	logger *zap.Logger
	table string
}

type ProxyRow struct {
	Host string
	Port int
	Username string
	Password string
}

func NewDbConn(user string, password string, dbName string) (db *sql.DB) {
	db, err := sql.Open("mysql", user+":"+password+"@/"+dbName)
	if err != nil {
		log.Fatalf("could not obtain db connection: %s", err.Error())
	}
	return db
}

func NewProxyDB(user string, password string, dbName string, table string) *ProxyDB {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("could not initiate zap logger: %s", err)
	}
	return &ProxyDB{
		dbConn: NewDbConn(user, password, dbName),
		logger: logger,
		table: table,
	}
}

func (pdb *ProxyDB) GetProxies(basesource string, proxyType string, limit int) ([]*ProxyRow, error) {
	if len(proxyType) > 0 {
		proxyType = fmt.Sprintf(" AND type='%s'", proxyType)
	}
	rows, err := pdb.dbConn.Query(
		fmt.Sprintf(
			"SELECT host, port, username, password" +
				" FROM %s WHERE basesource='%s' %s ORDER BY RAND() LIMIT %d",
				pdb.table, basesource, proxyType, limit))
	if err != nil {
		pdb.logger.Error("failed to query proxies from db", zap.String("error", err.Error()))
		return nil, err
	}
	var proxies []*ProxyRow
	for rows.Next() {
		proxy := &ProxyRow{}
		err = rows.Scan(&proxy.Host, &proxy.Port, &proxy.Username, &proxy.Password)
		if err != nil {
			pdb.logger.Error("failed scanning proxy row", zap.String("error", err.Error()))
			err = nil
			continue
		}

		proxies = append(proxies, proxy)
	}
	return proxies, nil
}