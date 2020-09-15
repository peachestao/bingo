package bingo

import (
	"image"
	"image/jpeg"

	//"loveChatApp/pkg/errors"
	"mime/multipart"
	"net/http"
	//"fmt"
	//"errors"
	"encoding/json"
	"fmt"
)

type Context struct{
	middleWareIndex int
	handlerIndex int
	handlers []HandlerFunc
	W ResponseWriter
	Req *http.Request
	bingo *Bingo
	Params Params
	DiyParam DiyParam
	tsr bool
}

func(c *Context)Next() (err error) {
	err =nil
	r := c.bingo.router
	mLen := len(r.middleWares)

	if c.middleWareIndex==mLen|| mLen == 0 { //if the last middleWare or no middleware,then invoke handlers
		path := c.Req.URL.Path

		handlers:=c.handlers
		if root := r.trees[c.Req.Method]; root != nil {
			if handlers!=nil{
				c.handlerIndex++ //关键代码 必须在handler执行前进行，否则会一直调handler[0],导致无限循环，无法执行下面的代码
				if c.handlerIndex > len(handlers) {
					panic("no enough handlerfunc to invoke,please check Next() function invoke time\n")
					//err=Error{"f","fas"}
					//mErr.New("f","fas")
					//err=mErr
					return
				}
				c.W.statusCode=200
				handlers[c.handlerIndex-1].ServeHTTP(c)
				return
			} else if c.Req.Method != http.MethodConnect && path != "/" {
				// Moved Permanently, request with GET method
				code := http.StatusMovedPermanently
				if c.Req.Method != http.MethodGet {
					// Permanent Redirect, request with same method
					code = http.StatusPermanentRedirect
				}

				if c.tsr && r.RedirectTrailingSlash {
					if len(path) > 1 && path[len(path)-1] == '/' {
						c.Req.URL.Path = path[:len(path)-1]
					} else {
						c.Req.URL.Path = path + "/"
					}
					http.Redirect(c.W, c.Req, c.Req.URL.String(), code)
					return
				}

				// Try to fix the request path
				if r.RedirectFixedPath {
					fixedPath, found := root.findCaseInsensitivePath(
						CleanPath(path),
						r.RedirectTrailingSlash,
					)
					if found {
						c.Req.URL.Path = fixedPath
						http.Redirect(c.W, c.Req, c.Req.URL.String(), code)
						return
					}
				}
			}


		}



		if c.Req.Method == http.MethodOptions && r.HandleOPTIONS {
			// Handle OPTIONS requests
			if allow := r.allowed(path, http.MethodOptions); allow != "" {
				c.W.Header().Set("Allow", allow)
				if r.GlobalOPTIONS != nil {
					r.GlobalOPTIONS.ServeHTTP(c.W, c.Req)
				}
				return
			}
		} else if r.HandleMethodNotAllowed { // Handle 405
		c.W.statusCode=405;
			if allow := r.allowed(path, c.Req.Method); allow != "" {
				c.W.Header().Set("Allow", allow)

				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed.ServeHTTP(c.W, c.Req)
				} else {
					http.Error(c.W,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}

		// Handle 404
		c.W.statusCode=404;
		if r.notFound != nil {
			r.notFound.ServeHTTP(c)
		} else {
			http.NotFound(c.W, c.Req)
		}

		//
		//
		//handlers := c.handlers
		//if handlers == nil {
		//	if c.bingo.mode == DEBUG {
		//		fmt.Printf("your request resource not found\n")
		//	}
		//	c.W.WriteHeader(http.StatusNotFound)
		//	c.W.Write([]byte("your request resource not found\n"))
		//
		//	return
		//}
		//c.handlerIndex++ //关键代码 必须在handler执行前进行，否则会一直调handler[0],导致无限循环，无法执行下面的代码
		//if c.handlerIndex > len(handlers) {
		//	panic("no enough handlerfunc to invoke,please check Next() function invoke time\n")
		//	//err=Error{"f","fas"}
		//	//mErr.New("f","fas")
		//	//err=mErr
		//	return
		//}
		//handlers[c.handlerIndex-1].ServeHTTP(c)
		//c.handlerIndex++ 关键代码 自增量不能放在handler执行后执行，因为会一直调handler[0],导致无限循环，该行代码永远无法执行
	} else { //if not last middleWare,then exec the next middleWare with inverted order
		c.middleWareIndex++ //同handlerIndex 必须在middleWare执行前进行
		r.middleWares[mLen-c.middleWareIndex]().ServeHTTP(c)
	}
	return
}

/*
返回json格式流
*/
 func(c *Context)JSON(httpStatus int , res Res){
 	//if value,ok:=res["msg"]; ok{
 	//	if value!=""{
	//		res["msg"]= errors.GetMsg((res["code"].(int)))+" "+value.(string)
	//	}else{
	//		res["msg"]= errors.GetMsg((res["code"].(int)))
	//	}
	//}else{
	//	res["msg"]= errors.GetMsg((res["code"].(int)))
	//}

	//将因token过期重新生成的token返回给客户端
	 newToken:=c.DiyParam.Get("newToken")
	 if newToken!=nil&&newToken!=""{
		 res["newToken"]=newToken
	 }

 	dataBytes,err:= json.Marshal(res)
 	if err!=nil{
 		panic("convert to json string error")
	}
	jsonStr:=string(dataBytes)
	fmt.Printf("json:"+jsonStr)
	c.W.statusCode=httpStatus
	c.W.Header().Set("Content-Type","application/json") //设置header必须放在WriteHeader前面，否则无效
	c.W.WriteHeader(httpStatus)

 	c.W.Write(dataBytes)
 }

 /*
 输出jpeg格式流
 */
func(c *Context)JPEG(disImg image.Image, quality... int){
	header := c.W.Header()
	header.Add("Content-Type", "image/jpeg")
	qua:=100
	if len(quality)>0{
		qua=quality[0]
	}
	jpeg.Encode(c.W,disImg,&jpeg.Options{Quality:qua}) //设置图片质量 范围：1-100
}

/*
自定义参数 用于在回话各函数间传值
*/
type DiyParam map[string]interface{}

func(dp DiyParam)Set(key string,value interface{}){
	dp[key]=value
}

func(dp DiyParam)Get(key string)(value interface{}){
	value,ok:=dp[key]
	if !ok{
		value=nil
	}
	return
}

/*
获取get方式传值
*/
func(c *Context)Query(key string)(value string){
	//value=c.Req.Form.Get(key)  这种方法有问题？
	value=c.Req.URL.Query().Get(key)
	return
}

/*
获取post方式传值
*/
func(c *Context)PostForm(key string)(value string){
	value=c.Req.PostForm.Get(key)
	return
}

/*
获取路由参数值   如:api/v1/users/:id  id的值
*/
func (c *Context)Param(name string)(value string){
	params:=c.Params
	for i := range params {
		if params[i].Key == name {
			return params[i].Value
		}
	}
	return ""
}

/*
获取上传文件
*/
func (c *Context)FormFile(name string)(file multipart.File,header *multipart.FileHeader,err error) {
	file, header, err = c.Req.FormFile(name)
	return
}

