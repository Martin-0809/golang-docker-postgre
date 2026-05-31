package repository

import (
	"context"
	"database/sql"
	"golang-docker-postgres/domain"
)

// postgresUserRepository 建立一個結構體，裡面持有資料庫連線 (sql.DB)
type postgresUserRepository struct {
	Conn *sql.DB
}

// NewPostgresUserRepository 初始化這個 Repository 的函式（給外部組裝用）
func NewPostgresUserRepository(conn *sql.DB) domain.UserRepository {
	return &postgresUserRepository{Conn: conn}
}

// GetByID 實作 GetByID 方法，真正去對 PostgreSQL 下 SQL 撈資料
func (p *postgresUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, name FROM users WHERE id = $1`
	
	user := &domain.User{}
	// 執行查詢並將結果寫入 user 結構體
	err := p.Conn.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}
// ... 原本的 GetByID 保持不動 ...

// Store 實作將資料寫入 PostgreSQL
func (p *postgresUserRepository) Store(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name) VALUES ($1) RETURNING id`
	
	// 執行寫入，並將自動生成的序列 ID 刷回 user 結構體中
	err := p.Conn.QueryRowContext(ctx, query, user.Name).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}