package dd_cmd

import (
	"log"
	"reflect"
	"strings"
	"testing"
)

type eee struct {
	IPrefixHandler
}
func (eee2 *eee) GetCommandPrefixs() []string{
	return []string{"/",">"}
}
func (eee2 *eee) GetPrefix(s string) string{
	prefixs := eee2.GetCommandPrefixs()
	for _, v := range prefixs {
		if strings.HasPrefix(s, v) {
			return v
		}
	}
	log.Println("无匹配")
	return ""
}
func TestParseCmd(t *testing.T) {
	type args struct {
		cmd    string
		engine IPrefixHandler
	}


	tests := []struct {
		name    string
		args    args
		want    Command
		wantErr bool
	}{
		{name: "/ command",args: args{cmd: CommandHelp("/","cmd 1 2"),engine: &eee{}},want: Command{Cmd: "cmd",prefix: "/",Params: []string{"1","2"}},wantErr: false},
		{name: "no command prefix has error",args: args{cmd: CommandHelp("+","cmd"),engine: &eee{}},want: Command{},wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCmd(tt.args.cmd, tt.args.engine)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCmd() got = %v, want %v", got, tt.want)
			}
		})
	}
}
