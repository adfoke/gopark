package main

import (
	"encoding/json"
	"gopark/pkg/hello"
	"github.com/sirupsen/logrus"
)

// 定义一个 User 结构体
type User struct {
	ID   uint `json:"id"`
	Name string `json:"name"`
	Mail  string `json:"mail"`
}

var log = logrus.New()


func main() {
	// 初始化一个 User 结构体
	user := User{
		ID:   1,
		Name: "test",
		Mail:  "test@gmail.com",
	}
	// 使用 json.Marshal() 方法将 User 结构体转换为 JSON 格式的数据
	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	// 打印 JSON 格式的数据
	log.Info(string(data))
	//引用hello包
	log.Info(hello.SayHello())


}