package bingo

import (
	"fmt"
	"time"
	"strconv"
)

func logger()HandlerFunc{
	return HandlerFunc(func(c *Context){
		startTime:= time.Now()
		c.Next()
		//time.Sleep(time.Second)
		endTime:=time.Now()
		path:=c.Req.URL.Path
		takeTime:=strconv.FormatInt(endTime.Sub(startTime).Milliseconds(),10)// strconv.FormatInt((endTime-startTime),10)
		method:=c.Req.Method
		statusCode:=strconv.Itoa(c.W.statusCode)
		//statusCode:=c.W.Header().
		fmt.Printf("%s %s %s %sms\n",path,method,statusCode,takeTime )
	})
}