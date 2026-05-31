// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	// 引入你剛剛建立的三個層級（零件）
	"golang-docker-postgres/delivery"
	"golang-docker-postgres/repository"
	"golang-docker-postgres/usecase"
)

const (
	// 指定要連接的DB位置
	HOST     = "db"
	DATABASE = "postgres"
	USER     = "user"
	PASSWORD = "mysecretpassword"
)

func main() {
    // 連接DB
    db, err := sql.Open(
        "postgres",
        fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE),
    )

    dbConnected := false // 💡 設一個貼心小貼紙，記錄到底有沒有連上資料庫

    if err != nil {
        fmt.Println("⚠️ 資料庫連線失敗:", err)
    } else {
        // 檢查連接是否成功
        if err = db.Ping(); err != nil {
            fmt.Println("⚠️ 資料庫 Ping 不到:", err)
        } else {
            fmt.Println("Successfully created connection to database")
            dbConnected = true // 💡 成功連上了！
        }
    }

    // 💡 只有當資料庫真的有連上，才去建立資料表
    if dbConnected {
        _, err = db.Exec(`
            CREATE TABLE IF NOT EXISTS users (
                id SERIAL PRIMARY KEY,
                name VARCHAR(100) NOT NULL
            );
        `)
        if err != nil {
            log.Println("自動建表失敗:", err)
        } else {
            fmt.Println("資料表檢查/建立完成")
        }
    } else {
        fmt.Println("⚠️ 因為沒有資料庫，跳過自動建表步驟，直接啟動網頁伺服器！")
    }

    // ==========================================
    // 💡 這裡以下就是你原本處理 /user 路由和啟動 Server 的地方
    // ==========================================
    // (請保留你原本的 http.HandleFunc 和動態 Port 程式碼...)
    port := os.Getenv("PORT")
    if port == "" {
        port = "5000"
    }
    fmt.Println("Server is running on 0.0.0.0:" + port)
    log.Fatal(http.ListenAndServe("0.0.0.0:" + port, nil))
}