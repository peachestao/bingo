package bingo

type Error struct{
	title string
	msg string
}
func(err Error) New(title string,msg string){
	err.title=title
	err.msg=msg
}
func(err Error)Error()string {
	return err.title+" "+err.msg
}