package models

// User 定义了用户的数据结构
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
}
