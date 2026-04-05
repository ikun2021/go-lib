package unitTest

import (
	"encoding/json"
	"net/http"
)

// Compute 一个简单的普通函数
func Compute(a, b int) int {
	return a + b
}

// User 一个结构体
type User struct {
	Name string
}

// FetchData User 的方法，包含网络请求（需要 Mock）
func (u *User) FetchData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return "real data", nil
}

// ValidateJSON 依赖标准库 json.Unmarshal（演示 Mock 标准库）
func ValidateJSON(data []byte) bool {
	var tmp interface{}
	err := json.Unmarshal(data, &tmp)
	return err == nil
}
