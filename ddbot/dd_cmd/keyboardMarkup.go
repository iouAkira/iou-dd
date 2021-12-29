package dd_cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type InlineKeyboardMarkup tgbotapi.InlineKeyboardMarkup

// Markupable 按钮包装接口
type Markupable interface {
	IkType() string                           //当前指令的返回
	MakeKeyboardMarkup() InlineKeyboardMarkup //对返回数据的按钮做包装
}

//RdcMarkup 读取cookies类
type RdcMarkup struct {
	Cmd     string
	Cookies []string
	RowBtns int
	Suffix  string
	Prefix  string
}

func (markup *RdcMarkup) IkType() string {
	return markup.Cmd
}

// MakeKeyboardMarkup 业务层待实现
func (markup *RdcMarkup) MakeKeyboardMarkup() InlineKeyboardMarkup {
	var keyboardMarkup InlineKeyboardMarkup
	t := make([]interface{}, len(markup.Cookies))
	for i, v := range markup.Cookies {
		t[i] = v
	}
	list := Slice_chunk(t, markup.RowBtns)
	for k, v := range list {
		row := tgbotapi.NewInlineKeyboardRow()
		for i, n := range v {
			cmd := Command{Cmd: markup.Cmd, prefix: markup.Prefix}
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(n.(string), cmd.Run(strconv.Itoa(i+k+1))))
			//row = append(row, tgbotapi.NewInlineKeyboardButtonData("📄 WSKEY", cmd.Run(fmt.Sprintf("%s %s", markup.Suffix, n.(string)))))
		}
		keyboardMarkup.InlineKeyboard = append(keyboardMarkup.InlineKeyboard, row)

	}
	return keyboardMarkup
}

//WrapCancelBtn 添加取消按钮
func WrapCancelBtn(markup *InlineKeyboardMarkup) InlineKeyboardMarkup {
	var cancelBTN = Command{Cmd: "cancel", prefix: "/"}

	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("取消", cancelBTN.Run()),
	))
	return *markup

}

//WrapCancelWithExampleBtn 添加取消按钮并带有范例显示
func WrapCancelWithExampleBtn(markup *InlineKeyboardMarkup) InlineKeyboardMarkup {
	var cancelBTN = Command{Cmd: "cancel"}
	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("取消", cancelBTN.Run()),
		tgbotapi.NewInlineKeyboardButtonURL("查看配置示例🔗", "https://github.com/iouAkira/someDockerfile/blob/master/dd_scripts/genCodeConf.list"),
	))
	return *markup
}

//SpnodeMarkup 读取cookies类
type SpnodeMarkup struct {
	Cmd      string
	FilePath string
	RowBtns  int
	Suffix   string
}

func (markup *SpnodeMarkup) IkType() string {
	return markup.Cmd
}

// MakeKeyboardMarkup todo 业务层待实现
func (markup *SpnodeMarkup) MakeKeyboardMarkup() InlineKeyboardMarkup {
	var keyboardMarkup InlineKeyboardMarkup
	return keyboardMarkup
}

//LogsMarkup 读取cookies类
type LogsMarkup struct {
	Cmd      string
	FilePath string
	RowBtns  int
	Suffix   string
}

func (markup *LogsMarkup) IkType() string {
	return markup.Cmd
}

// MakeKeyboardMarkup todo 业务层待实现
func (markup *LogsMarkup) MakeKeyboardMarkup() InlineKeyboardMarkup {
	var keyboardMarkup InlineKeyboardMarkup
	return keyboardMarkup
}

// MakeKeyboard 实例化空按钮对象
func MakeKeyboard() InlineKeyboardMarkup {
	var keyboardMarkup InlineKeyboardMarkup
	return keyboardMarkup
}

func (markup InlineKeyboardMarkup) WithCancel() InlineKeyboardMarkup {
	return WrapCancelBtn(&markup)
}

func (markup InlineKeyboardMarkup) WithBack(cmd string) InlineKeyboardMarkup {
	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("返回", cmd),
	))
	return markup
}
func (markup InlineKeyboardMarkup) WithExampleBtn() InlineKeyboardMarkup {
	return WrapCancelWithExampleBtn(&markup)
}

func (markup InlineKeyboardMarkup) WithCommand(cmd Executable) InlineKeyboardMarkup {
	if cmd == nil {
		return markup
	}
	description := cmd.Description()
	if description == "" {
		description = cmd.GetCmd()
	}
	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(description, cmd.Run()),
	))
	return WrapCancelBtn(&markup)
}

func (markup InlineKeyboardMarkup) WithCommandStr(cmd, description string) InlineKeyboardMarkup {
	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(description, cmd),
	))
	return WrapCancelBtn(&markup)
}

func (markup InlineKeyboardMarkup) Get() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.InlineKeyboardMarkup(markup)
}

// Slice_chunk 切片数据分组
func Slice_chunk(slice []interface{}, size int) (chunkslice [][]interface{}) {
	if size >= len(slice) {
		chunkslice = append(chunkslice, slice)
		return
	}
	end := size
	for i := 0; i <= (len(slice) - size); i += size {
		chunkslice = append(chunkslice, slice[i:end])
		end += size
	}
	return
}
