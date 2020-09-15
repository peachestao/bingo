package bingo

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int //增加状态码字段，方便记录请求日志
}

//func(w ResponseWriter)WriteHeader(code int){
//	//w.ResponseWriter.Write([]byte(data))
//}

