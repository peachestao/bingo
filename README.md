## bingo

一个轻量级的golang restful api框架

## 快速使用

### 1、下载 go get github.com/peachestao/bingo
### 2、例子
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
## 功能介绍
### 1、中间件
用法：
engine.use(m MiddleWare),支持链式调用engine.use(m1).use(m2)...
例如，实现一个记录所有请求的处理时间的功能，我们可以定义如下方法
~~~
func logger()bingo.HandlerFunc{
	return bingo.HandlerFunc(func(c *bingo.Context){
		startTime:= time.Now()
		c.Next()
		endTime:=time.Now()
		path:=c.Req.URL.Path
		takeTime:=strconv.FormatInt(endTime.Sub(startTime).Milliseconds(),10)// strconv.FormatInt((endTime-startTime),10)
		method:=c.Req.Method
		fmt.Printf("%s %s %sms\n",path,method,takeTime )
	})
}
~~~
然后engine.use(logger)即可
### 2、路由
支持GET、POST、PUT、DELETE等http谓词,用法如下：
~~~
engine.GET("/api/v1/hello", func(c *bingo.Context) {

}))
engine.POST("/api/v1/hello", func(c *bingo.Context) {

}))
engine.DELETE("/api/v1/hello", func(c *bingo.Context) {

}))
engine.PUT("/api/v1/hello", func(c *bingo.Context) {

}))
~~~
### 3、请求上下文
将req,res数据封装在bingo.Context对象中:
1、Query方法用于获取get传值
2、PostForm方法用于后去post传值
未完，待续~
