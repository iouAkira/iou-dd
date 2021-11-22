package dd_cmd

// IRoutes 定义所有路由 Handler 接口。
type IRoutes interface {
	RegCommand(string, string, ...HandlerFunc) IRoutes
	Handle(Executable, string, ...HandlerFunc) IRoutes
}

// IRouter 定义了所有的路由器Handler接口，包括单路由器和组路由器。
type IRouter interface {
	RegCommand(string, string, ...HandlerFunc) IRoutes
	Handle(Executable, string, ...HandlerFunc) IRoutes
	// Group(string, ...HandlerFunc) *RouterGroup
}

//Executable 对于一条指令来说需要用到以下两个方法，分别是Description和Run方法，
//Description 方法对当前指令的描述，返回值是一个字符串，
//Run 方法是执行当前指令的具体操作，
type Executable interface {
	Description(...string) string
	Run(...string) string
	Prefix() string
	SetCmd(string)
	GetCmd() string
}

// ParseExec 解析字符串参数
func ParseExec(args ...string) Executable {
	return nil
}
