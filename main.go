package main

import (
	"github.com/aman/filter"
	"github.com/aman/iocontrol"
	"github.com/aman/modules"
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// 引数取得
	var args []string = modules.Parse()

	// コマンド実行
	var commandResult string = modules.ExecMan(args)

	manLists := modules.AnalyzeOutput(commandResult)
	var inputs = ""

	iocontroller := iocontrol.NewIoController(manLists)
	iocontrol.RenderQuery(&inputs)
	pageList := iocontroller.LocatePages(manLists)
	iocontroller.RenderResult(manLists, pageList[:])
	for {
		if iocontroller.ReceiveKeys(&inputs) == -1 {
			return
		}
		iocontrol.RenderQuery(&inputs)
		result := filter.IncrementalSearch(&inputs, manLists)
		pageList = iocontroller.LocatePages(result)
		iocontroller.RenderResult(result, pageList[:])
	}
}
