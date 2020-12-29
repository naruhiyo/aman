package iutil

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

	sutil "aman/struct/util"

	"github.com/go-vgo/robotgo"
	"github.com/mattn/go-pipeline"
)

type CommandStruct sutil.CommandStruct

/*
 * コンストラクタ作成
 */
func NewCommand() *CommandStruct {
	return &CommandStruct{
		ManResult: "",
	}
}

/*
 * @description man コマンドを実行する
 * @params args 実行時引数
 */
func (myself *CommandStruct) ExecMan(commands []string) {
	// man コマンドは空白区切のコマンドをハイフンで管理しているため、ハイフンつなぎに変更
	var command string = strings.Join(commands, "-")

	// manコマンドを実行する
	//   - manの結果には\bや\tが入っているためcolで
	//   - \bを除外し、\tを半角スペースに変換する
	out, err := pipeline.Output(
		[]string{"man", command},
		[]string{"col", "-bx"},
	)

	if err != nil {
		panic(errors.New("No man command results"))
	}

	myself.ManResult = string(out)
}

/*
 * @description オプション付きコマンドをターミナルに出力する
 * @params commands 実行時引数
 * @params options 選択したオプション
 */
func (myself *CommandStruct) CmdOutput(commands []string, options []string) {
	// エコーバックを OFF
	myself.ExecWithStdin("stty", "-echo")
	// コマンドをターミナル上に出力
	var result string = strings.Join(commands, " ") + " " + strings.Join(options, " ")
	// ターミナルをクリアする
	robotgo.TypeStr(result)
	// システムが記憶している入力をクリア
	time.Sleep(time.Millisecond * 15)
	// エコーバックを ON
	myself.ExecWithStdin("stty", "echo")
}

/*
 * @description コマンドを標準入力から実行する
 * @params name コマンド
 * @params option コマンドオプション
 */
func (myself *CommandStruct) ExecWithStdin(name string, option ...string) {
	c := exec.Command(name, option...)
	c.Stdin = os.Stdin
	c.Run()
}
