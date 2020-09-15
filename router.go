package bingo

import (
	"net/http"
	"fmt"
	"strings"
	"sync"
)

type HandlerFunc func(*Context)

func(h HandlerFunc)ServeHTTP(c *Context){
	/*
		关键代码 relation A context不能定义为全局变量，否则多个请求时单个请求中的上下文数据会混乱 如：middleWareIndex当前请求为2，
	下次请求就从2开始执行handler,显然是有问题的，context要放在每个请求入口处理函数中初始化，然后作为参数一路传递到中间件最后执行handler

	context.Req=req
	context.W=w
	h(context)
		 */

	h(c)
}
type MiddleWare func()HandlerFunc
type Group struct{
	path string

}
type Router struct{
	middleWares []MiddleWare
	mux map[string][]HandlerFunc
	trees map[string]*node
	group Group
	globalAllowed string
	paramsPool sync.Pool
	maxParams  uint16
	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 308 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 308 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// An optional http.Handler that is called on automatic OPTIONS requests.
	// The handler is only called if HandleOPTIONS is true and no OPTIONS
	// handler for the specific path was set.
	// The "Allowed" header is set before calling the handler.
	GlobalOPTIONS http.Handler

	// Cached value of global (*) allowed methods

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	notFound HandlerFunc

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}
type Param struct {
	Key   string
	Value string
}
type Params []Param

func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}



//var context=&Context{bingo:bingo} //关键代码 relation A context不能定义为全局变量，否则多个请求时单个请求中的上下文数据会混乱 如：middleWareIndex当前请求为2，
// 下次请求就从2开始执行handler,显然是有问题的，context要放在每个请求入口处理函数中初始化，然后作为参数一路传递到中间件最后执行handler

/*
add a middleWare,support chain invoke
*/
func(r *Router)Use(m MiddleWare)*Router{
	r.middleWares=append(r.middleWares,m)
	return r
}
///*
//add a handler func,support chain invoke
//*/
//func(r *Router)Add(route string,handlers ...func( *Context)) *Router {
//	for _,h:=range handlers{
//		r.mux[route]=append(r.mux[route],HandlerFunc(h))
//	}
//	return r
//}

func(r *Router)GET(path string, handlers ...func(c *Context)) *Router{
	r.Handle(http.MethodGet, path,handlers)
	return r
}

func(r *Router)HEAD(path string, handlers ...func(c *Context)) *Router{
	r.Handle(http.MethodHead, path, handlers)
	return r
}
func(r *Router)OPTIONS(path string, handlers ...func(c *Context)) *Router{
	r.Handle(http.MethodOptions, path, handlers)
	return r
}
func(r *Router)PATH(path string, handlers ...func(c *Context)) *Router{
	r.Handle(http.MethodPatch, path, handlers)
	return r
}

func(r *Router)POST(path string, handlers ...func(c *Context)) *Router{
	r.Handle(http.MethodPost, path, handlers)
	return r
}
func (r *Router)PUT(path string ,handlers ...func(c *Context)) *Router{
	r.Handle(http.MethodPut,path,handlers)
	return r
}
func(r *Router)DELETE(path string,handlers ...func(c *Context)) *Router{
	r.Handle(http.MethodDelete,path,handlers)
	return r
}

func (r *Router) Handle(method, path string, handlers []func(c *Context)) {
	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if handlers == nil {
		panic("handle must not be nil")
	}


	var handlersNew []HandlerFunc
	for _,handle:=range handlers{
		handlersNew=append(handlersNew,HandlerFunc(handle))
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root

		r.globalAllowed = r.allowed("*", "")
	}

	root.addRoute(path, handlersNew)

	// Update maxParams
	if pc := countParams(path); pc > r.maxParams {
		r.maxParams = pc
	}

	// Lazy-init paramsPool alloc func
	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)
			return &ps
		}
	}
}

func(r *Router)NotFound(handler func(c *Context))*Router{
	r.notFound=handler
	return r
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" { // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range r.trees {
				if method == http.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else { // specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == http.MethodOptions {
				continue
			}

			handle, _, _ := r.trees[method].getValue(path, nil)
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Add request method to list of allowed methods
		allowed = append(allowed, http.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}
	return
}

/*
start a http server
*/
func (r *Router)Run(addr string )(err error){
	err=http.ListenAndServe(addr,r)
    if(err!=nil) {
		fmt.Printf("error hapend:%s", err.Error())
	}
	return
}

func (r *Router) getParams() *Params {
	ps := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0] // reset slice
	return ps
}

func (r *Router) putParams(ps *Params) {
	if ps != nil {
		r.paramsPool.Put(ps)
	}
}
/*
every request's Entrance func
每个请求的入口函数
*/
func(r *Router)ServeHTTP(w http.ResponseWriter,req *http.Request){
	/*
	关键代码 relation A context不能定义为全局变量，否则多个请求时单个请求中的上下文数据会混乱 如：middleWareIndex当前请求为2，
	下次请求就从2开始执行handler,显然是有问题的，context要放在每个请求入口处理函数中初始化，然后作为参数一路传递到中间件最后执行handler
	*/
	context:=&Context{W:ResponseWriter{w,0},Req:req,bingo:bingo,handlers:make([]HandlerFunc,1),DiyParam:make(map[string]interface{},0)}

	//if no bind middleWares then invoke the handler
	path:=req.URL.Path
	//handlers:=r.mux[path]
	//if handlers==nil{
	//	w.WriteHeader(http.StatusNotFound)
	//	w.Write([]byte("your request resource not found\n"))
	//	return
	//}
	//handlers[0].ServeHTTP(context)

	if root := r.trees[req.Method]; root != nil {
		if handles, ps, tsr := root.getValue(path, r.getParams); handles != nil {
			context.handlers=handles
			context.tsr=tsr
			if ps != nil {
				//context.handlers=handles
				context.Params=*ps
				r.putParams(ps)
			} else {
				//handles(w, req, nil)
				context.Params=nil
			}
		}else{
			context.handlers=nil
		}
	}
	context.Next()


}