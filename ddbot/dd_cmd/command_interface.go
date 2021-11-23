package dd_cmd

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
