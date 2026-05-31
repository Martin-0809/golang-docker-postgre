package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"golang-docker-postgres/domain"
)

// UserHandler 負責接收 HTTP 請求，它需要用到 Usecase 來處理業務
type UserHandler struct {
	UserUsecase domain.UserUsecase
}

// GetUserByID 處理 /user?id=123 這樣的 HTTP 請求
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// 1. 從 URL 參數取得 id 並且轉成數字
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "無效的 ID", http.StatusBadRequest)
		return
	}

	// 2. 呼叫 Usecase 層處理邏輯
	user, err := h.UserUsecase.GetProfile(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. 將結果轉成 JSON 回傳給 Client 端
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
// ... 原本的 GetUserByID 保持不動 ...

// PostToCreateUser 處理 POST /user 請求
func (h *UserHandler) PostToCreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. 限制必須是 POST 請求
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允許", http.StatusMethodNotAllowed)
		return
	}

	// 2. 解析前端傳來的 JSON 內容
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "無效的 JSON 格式", http.StatusBadRequest)
		return
	}

	// 3. 建立 domain 實體並呼叫 Usecase
	newUser := &domain.User{
		Name: req.Name,
	}
	if err := h.UserUsecase.Store(r.Context(), newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. 成功，回傳建立好的完整資料（含資料庫生成的 ID）
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}