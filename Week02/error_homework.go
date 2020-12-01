package main

import (
	"database/sql"
	"errors"
	"fmt"
	gerrors "github.com/pkg/errors"
)

// 定义对象书
type book struct {
	//.... 其它字段
	content string
}

func main() {
	service(0)
	service(1)
}

// 数据表执行
func dbQuery(id int) (*book, error) {
	if id == 0 {
		return nil, sql.ErrNoRows
	} else {
		return &book{content: "I am a book."}, nil
	}
}

// Dao 层
func dao(id int) (*book, error) {
	b, err := dbQuery(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no data found")
		}
		return nil, gerrors.Wrap(err, "other error")
	}
	return b, nil
}

// business层
func business(id int) (*book, error) {
	// 业务逻辑代码
	// 直接返回错误
	return dao(id)
}

//接入层代码
func service(id int) {
	b, err := business(id)
	if err != nil {
		fmt.Printf("get book[%d] error: %v\n", id, err)
		return
	}
	fmt.Printf("get book[%d] content:%s\n", id, b.content)
}
