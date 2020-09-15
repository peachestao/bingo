package bingo

type Bingo struct{
	router *Router
	mode int8
}

type ResEntity struct{
	Status int
	Code int
	Msg string
	Data interface{}
}

type Res map[string]interface{}

/*run mode list */
const DEBUG int8=0
const PRODUCE int8=1
var bingo *Bingo
/*
initialize bingo variable
*/
func init(){
	 bingo=&Bingo{router:&Router{RedirectTrailingSlash:  true,
		 RedirectFixedPath:      true,
		 HandleMethodNotAllowed: true,
		 HandleOPTIONS:          true, mux:make(map[string][]HandlerFunc)}}
}

/*
initialize a bingo engine,in fact is a router
*/
func New()*Router{
	if bingo.mode==DEBUG{
		bingo.router.Use(logger)
	}

	return bingo.router
}
/*
set run mode,has two mode:DEBUG mode:would print log in console,PRODUCE mode:would not print log
 */
func SetMode(mode int8){
	bingo.mode=mode
}



