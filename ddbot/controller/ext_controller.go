package controller

import (
	"ddbot/dd_cmd"
	"log"
)

// ExtansionController 设想对扩展功能实现动态执行的功能[所谓的插件化扩展]，由于golang本身对动态链接库的孱弱支持
// (或者说从业务层面就已经放弃了这种本地化的开发思路，转而采用分布式微服务架构来优化)，针对单机程序初步设想是外挂执行即对多语言的shell执行(利用本身的exec库实现)，
// 在ext {param} param检索到的js或者其他语言文件后执行
// Todo: 脚本路径待商榷,是需要在env.sh中给定还是通过某种方式动态链路?
func ExtansionController(ctx *dd_cmd.Context) {
	message := ctx.Message(ctx)
	path := ctx.Vars()
	log.Println(message)
	log.Println(path)
}

//可扩展文件信息接口
type IExtFile interface {
	GetPath() string     //获取绝对路径
	SetPath(string)      //设置绝对路径
	GetLanguage() string //获取执行语言
	SetLanguage(string)  //设置执行语言
}

//可扩展文件执行接口
type IExtCmd interface {
	GetFile(string) IExtFile                 //通过文件信息获取可执行扩展文件
	ExecFile(args ...string) (string, error) //可选参数对扩展文件执行并导出结果
	SetFile(IExtFile) error                  //注入可执行扩展文件
	CheckFile() error                        //检查文件是否合法
}
