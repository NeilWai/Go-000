# 第二周：异常处理  学习笔记

## **Week02 作业题目：**

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？



## 个人理解

**当Dao层遇到sql.ErrNoRows时，我认为不应该Wrap这erro抛给上层。**

因为sql.ErrNoRows是数据库底层的错误，业务层不需要关心，而且sql.ErrNoRows是有局限性的，不方便dao层扩展。应该用自定义error返回。

业务层建议wrap Dao层的error直接返回到接入层。

一般error日志由最上层打印输出。

## 代码示例

```go
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

func main()  {
	service(0)
	service(1)
}

// 数据表执行
func dbQuery(id int) (*book,error)  {
	if id==0 {
		return nil,sql.ErrNoRows
	} else {
		return &book{content:"I am a book."},nil
	}
}

// Dao 层
func dao(id int) (*book,error) {
	b,err := dbQuery(id)
	if err!=nil {
		if errors.Is(err,sql.ErrNoRows) {
			return nil,errors.New("no data found")
		}
		return nil, gerrors.Wrap(err,"other error")
	}
	return b,nil
}

// business层
func business(id int)(*book,error)  {
	// 业务逻辑代码
	// 直接返回错误
	return dao(id)
}
//接入层代码
func service(id int) {
	b,err := business(id)
	if err!=nil {
		fmt.Printf("get book[%d] error: %v\n",id,err)
		return
	}
	fmt.Printf("get book[%d] content:%s\n",id,b.content)
}

```

