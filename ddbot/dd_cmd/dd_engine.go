package dd_cmd

import (
	"fmt"
	"log"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Engine 中注入了Bot所需要的一些，相当与一个大框架，使用时候需要New()来对Engine初始化操作
type Engine struct {
	RouterGroup
	trees  commandTrees
	pool   sync.Pool
	Token  string
	Userid int64
	bot    *tgbotapi.BotAPI
}

// New 返回一个 Engine 实体,初始化操作并不包含任何路由和中间件
func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
	}
	engine.RouterGroup.engine = engine
	return engine
}

// addRoute 对指令集的添加 method 实则是指令的 prefix，函数创建了一个trees挂载到 Engine 中，
// 将所有method相同的指令合并到一颗树组合中。
func (engine *Engine) addRoute(method Executable, path string, handlers HandlersChain) {
	log.Printf("> route on [ %s %s ]", method.Prefix(), path)
	//查询 method 节点是否存在  例如 "/" "@"等这些是root节点，在root节点下存在一个trees切片集合包含该方法下的所有指令，
	//对为空对象是初始化添加节点组
	root := engine.trees.get(method)
	if root == nil {
		tree := new(commandTree)
		tree.method = method
		tree.root = &commandNodes{}
		engine.trees = append(engine.trees, tree)
		root = engine.trees.get(method)
	}
	//todo 解决获取节点到底是哪里的问题
	root.root.addNode(path, handlers)
}

// 获取当前engine的指令前缀集合
func (engine *Engine) GetCommandPrefixs() []string {
	var prefixs []string
	for _, v := range engine.trees {
		prefixs = append(prefixs, v.method.Prefix())
	}
	return prefixs
}

// Run 主入口函数包含两个两个参数，token tg的bot token,userid 用户的id
// 函数在启动时候会直接将 tgbotapi 进行初始化操作,并在结束进程时候停止信息的接受
func (engine *Engine) Run(token string, userid int64) {
	engine.Token = token
	engine.Userid = userid
	bot, err := tgbotapi.NewBotAPI(token)
	bot.Debug = false
	if err != nil {
		log.Panicf("start bot failed with some error %v", err)
	}
	log.Printf("Telegram bot stared，Bot info ==> %s %s[%s]", bot.Self.FirstName, bot.Self.LastName, bot.Self.UserName)
	c := &Context{Request: bot}
	c.reset()
	c.Request = bot
	c.engine = engine
	engine.pool.Put(c)
	engine.bot = bot
	engine.handleHTTPRequest(c)
}

// handleHTTPRequest
func (engine *Engine) handleHTTPRequest(c *Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.Request.GetUpdatesChan(u)
	for update := range updates {
		//log.Println(update.Message)
		//log.Println(update.CallbackQuery)
		if update.Message != nil && update.Message.Document != nil {
			fmt.Printf("%+v\n", "sss")
			continue
		}
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		c.mu.Lock()
		c.Update = &update
		c.mu.Unlock()
		engine.handleRequest(c)
	}
}

// handleRequest 对请求数据处理进不同的路由中,由解析器将用户返回的消息解析加工到 Executable 中,
func (engine *Engine) handleRequest(c *Context) {
	var msg string
	if c.Update.Message != nil {
		log.Printf("来源Message: %s", c.Update.Message.Text)
		msg = c.Update.Message.Text
	} else {
		log.Printf("来源CallbackQuery: %s", c.Update.CallbackQuery.Data)
		msg = c.Update.CallbackQuery.Data
	}

	msg = strings.Trim(msg, " ")
	//对路由集合遍历的查询开头与请求一致的指令
	//对入口路由进行查询到各个路由树中
	var routes *commandTree
	var msgPrefix string
	for _, tree := range engine.trees {
		log.Printf("查询路由: %v", tree.method.Prefix())
		if strings.HasPrefix(msg, tree.method.Prefix()) {
			routes = tree
			msgPrefix = tree.method.Prefix()
		}
	}

	if routes != nil && len(*routes.root) > 0 {
		if c.Request.Debug {
			log.Printf("入口命令：%s", msgPrefix)
		}
		for _, route := range *routes.root {
			log.Printf("%+v", fmt.Sprintf("%s%s", routes.method.Prefix(), route.path))
			if fmt.Sprintf("%s%s", routes.method.Prefix(), route.path) == cleanPath(msg, 0)[0] {
				if c.Request.Debug {
					log.Printf("存在命令 %s", route.path)
					log.Printf("等待命令 %s", msg)
					log.Printf("接受命令 %s", msg)
				}
				for _, handler := range route.handlers {
					handler(c)
				}
			}
		}

	}
}

func cleanPath(cmd string, offset int) []string {
	cmdMsgSplit := strings.Split(cmd[offset:], " ")
	var arr []string
	for _, v := range cmdMsgSplit {
		if v == "" {
			continue
		}
		arr = append(arr, v)
	}
	return arr
}
