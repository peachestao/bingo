## bingo

一个轻量级的golang restful api web框架,基于httprouter组件

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
### 1、handler
一个请求支持任意多的handler处理函数,通过Context的Next方法串连执行，如果想终止后续调用不调用Next方法返回响应即可，这种形式经常用于用户身份验证，如果不合法提前退出  
~~~
engine.GET("/api/v1/hello", func(c *bingo.Context) {

                token:=c.Query('token')
		if !checkValid(token){
		
		 c.JSON(200,bingo.Res{
		        "status":-1,
			"msg":"认证不通过",
			"data":null,		
		})
		
		return
		
		}
		
		c.Next()
		
	}, func(c *bingo.Context) {

		c.JSON(200,bingo.Res{
			"status":0,
			"msg":"这是一个轻量级的golang restful api风格的后端框架",
			"data":"hello bingo",
		})
	})
~~~
一般我们将handler分成单独函数的形式，代码组织清爽  

~~~

func auth(c *bingo.Context){

                token:=c.Query("token")
		if !checkValid(token){
		
		 c.JSON(200,bingo.Res{
		        "status":-1,
			"msg":"认证不通过",
			"data":null,		
		})
		
		return
		
		}		
		c.Next()
}

func hello(c *bingo.Context) {

		c.JSON(200,bingo.Res{
			"status":0,
			"msg":"这是一个轻量级的golang restful api风格的后端框架",
			"data":"hello bingo",
		})
	})

engine.GET("/api/v1/hello", auth, hello)
~~~

### 2、中间件
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
### 3、路由
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
### 4、请求上下文
将req,res数据封装在bingo.Context对象中:  
1）、Query方法用于获取get传值  
~~~
engine.GET("/api/v1/hello", func(c *bingo.Context) {
               name:=c.Query("name") 

		c.JSON(200,bingo.Res{
			"status":0,
			"msg":"这是一个轻量级的golang restful api风格的后端框架",
			"data":name,
		})
	})
~~~
2）、PostForm方法用于获取post 表单传值  
~~~
engine.GET("/api/v1/hello", func(c *bingo.Context) {
               name:=c.PostForm("name") 

		c.JSON(200,bingo.Res{
			"status":0,
			"msg":"这是一个轻量级的golang restful api风格的后端框架",
			"data":name,
		})
	})
~~~
3）、Param方法获取存储url中":id"形式的可变参数值，通过Param("id")形式获取 

如：http://127.0.0.1:8008/api/v1/hello/peachesTao
~~~
engine.GET("/api/v1/hello/:uid", func(c *bingo.Context) {

               uid:=c.Param("uid")) 

		c.JSON(200,bingo.Res{
			"status":0,
			"msg":"这是一个轻量级的golang restful api风格的后端框架",
			"data":"hello "+uid,
		})
	})
~~~

4）、JSON方法返回json格式响应数据,第一个参数为http 状态码，第二个参数为map类型的数据，可以自定义键值对  
~~~
engine.GET("/api/v1/hello", func(c *bingo.Context) {    
		c.JSON(200,bingo.Res{
			"status":0,
			"msg":"这是一个轻量级的golang restful api风格的后端框架",
			"data":"hello bingo",
		        "diyKey":"diyValue",	
		})
	})
~~~
5）、DiyParam属性用于各个handler之间传值，set方法用于设置某个key的值，get方法用于获取某个key值。它非常有用，比如请求过来后第一个handler负责认证,如果通过则解析出uid，然后通过DiyParam Set方法传递给其他handler,其他handler直接通过Get方法获取，不需要每个handler要获取uid时都去解析  
~~~

func auth(c *bingo.Context){

                token:=c.Query("token")
		if valided,uid:=checkValid(token);!valided{
		
		 c.JSON(200,bingo.Res{
		        "status":-1,
			"msg":"认证不通过",
			"data":null,		
		})
		
		return
		
		}
		c.DiyParam.Set("uid",uid) //认证通过后将解析出来的uid存储到DiyParam中，供后续handler调用
		c.Next()
}

func hello(c *bingo.Context) {

               uid:=c.DiyParam.Get("uid") //获取存储在DiyParam中的uid
		c.JSON(200,bingo.Res{
			"status":0,
			"msg":"这是一个轻量级的golang restful api风格的后端框架",
			"data":"hello "+uid,
		})
	})

engine.GET("/api/v1/hello", auth, hello)
~~~
### 如果觉得对你有帮助，欢迎Star支持一下^_^

