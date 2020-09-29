package main

import (
	"fmt"

	"github.com/aman/modules"
)

func main() {
	// 引数取得
	var args = modules.Parse()

	// コマンド実行
	var commandResult = modules.GetOptions(args)

	// オプションだけ取得
	var optionList = modules.AnalyzeOutput(commandResult)
	fmt.Println(optionList)
}
