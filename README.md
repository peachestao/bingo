# bingo
一个轻量级的golang restful api框架
# 快速使用

## 1、下载 go get https://github.com/peachestao/bingo
## 2、例子
~~~
 package main

import (
	"github.com/peachestao/bingo"
)

func main(){
	engine := bingo.New()
	engine.GET("/api/v1/hello", func(c *bingo.Context) {

		c.JSON(200,bingo.Res{
			"status":0,
			"msg":"这是一个轻量级的golang restful api风格的后端框架",
			"data":"hello bingo",
		})
	})

	engine.Run("127.0.0.1:8008")
}
~~~
