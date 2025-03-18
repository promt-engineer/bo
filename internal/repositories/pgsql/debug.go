package pgsql

import (
	"gorm.io/gorm"
)

type DebugRepository struct {
	conn   *gorm.DB
	dbName string
}

func NewDebugRepository(conn *gorm.DB, dbName string) *DebugRepository {
	return &DebugRepository{conn: conn, dbName: dbName}
}

func (r *DebugRepository) SizeMB() (float64, error) {
	var size float64

	r.conn.Raw("select pg_database_size(?) as size", r.dbName).Scan(&size)

	return size / 1024 / 1024, nil
}
