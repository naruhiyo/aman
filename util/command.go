package util

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

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
	execWithStdin("stty", "-echo") // エコーバックを OFF
	// コマンドをターミナル上に出力
	var command string = strings.Join(args, " ") + " " + strings.Join(stackOptions, " ")
	// ターミナルをクリアする
	robotgo.TypeStr(command)
	time.Sleep(time.Millisecond * 15) // システムが記憶している入力をクリア
	execWithStdin("stty", "echo")     // エコーバックを ON
}

/**
* @description コマンドを標準入力から実行する
* @params name コマンド
* @params option コマンドオプション
**/
func execWithStdin(name string, option ...string) {
	c := exec.Command(name, option...)
	c.Stdin = os.Stdin
	c.Run()
}
