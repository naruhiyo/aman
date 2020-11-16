package util

import (
	"errors"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/mattn/go-pipeline"
)

/**
* @description man コマンドを実行する
* @params args 実行時引数
**/
func ExecMan(args []string) string {
	// man コマンドは空白区切のコマンドをハイフンで管理しているため、ハイフンつなぎに変更
	const MAN string = "man"
	var command string = strings.Join(args, "-")

	// manコマンドを実行する
	//   - manの結果には\bや\tが入っているためcolで
	//   - \bを除外し、\tを半角スペースに変換する
	out, err := pipeline.Output(
		[]string{MAN, command},
		[]string{"col", "-bx"},
	)

	if err != nil {
		panic(errors.New("Error: No results"))
	}

	return string(out)
}

/**
* @description オプション付きコマンドをターミナルに出力する
* @params args 実行時引数
* @params stackOptions 選択したオプション
**/
func CmdOutput(args []string, stackOptions []string) {
	// コマンドをターミナル上に出力
	var command string = strings.Join(args, " ") + " " + strings.Join(stackOptions, " ")
	robotgo.TypeStr(command)
}
