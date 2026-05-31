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
	if err != nil {
		panic(err)
	}
	// 檢查連接是否成功
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Successfully created connection to database")
	// === 自動建立資料表 ===
_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL
    );
`)
if err != nil {
    log.Println("自動建表失敗:", err)
} else {
    log.Println("users 資料表檢查/建立成功！")
}
// =====================

	// ==========================================
	// 核心大絕招：Clean Architecture 零件組裝 (DI)
	// ==========================================

	// 1. 把真正的資料庫連線變數 db 傳入 Repository
	userRepo := repository.NewPostgresUserRepository(db)

	// 2. 把剛剛做好的 userRepo 傳入 Usecase
	userUsecase := usecase.NewUserUsecase(userRepo)

	// 3. 把剛剛做好的 userUsecase 傳入 Delivery/Handler
	userHandler := &delivery.UserHandler{
		UserUsecase: userUsecase,
	}

	// ==========================================
	// 路由設定
	// ==========================================

	// 原本的測試路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!!")
	})

	// 新增的 Clean Architecture 路由：取得使用者資料
	// 當瀏覽器輸入 http://localhost:5000/user?id=1 時，就會觸發我們整套架構
	// 修改後的 /user 路由：根據 HTTP 方法分流
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUserByID(w, r)
		case http.MethodPost:
			userHandler.PostToCreateUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// 💡 自動動態偵測雲端主機分配的 Port，找不到就預設用 5000
port := os.Getenv("PORT")
if port == "" {
    port = "5000"
}

// 啟動 http server
fmt.Println("Server is running on :" + port)
log.Fatal(http.ListenAndServe(":" + port, nil))
}