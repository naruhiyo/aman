package main

import (
	"github.com/aman/filter"
	"github.com/aman/iocontrol"
	"github.com/aman/modules"
	"github.com/nsf/termbox-go"
)

const (
	ESCAPE     = -1
	KEY        = 1
	ARROW_UP   = 90
	ARROW_DOWN = 91
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

	var manLists []string = modules.AnalyzeOutput(commandResult)
	var inputs string = ""
	var selectedPos int = -1
	var result []string = manLists

	iocontroller := iocontrol.NewIoController(result)
	iocontrol.RenderQuery(&inputs)
	pageList := iocontroller.LocatePages(result)
	iocontroller.RenderResult(selectedPos, result, pageList[:])
	for {
		var keyStatus int = iocontrol.ReceiveKeys(&inputs)
		iocontrol.RenderQuery(&inputs)

		switch keyStatus {
		case KEY:
			// 毎回 man 結果に対して検索を行う
			result = filter.IncrementalSearch(&inputs, manLists)
		case ARROW_UP:
			if selectedPos > 0 {
				selectedPos--
			}
		case ARROW_DOWN:
			if selectedPos < len(result)-1 {
				selectedPos++
			}
		case ESCAPE:
			return
		}
		pageList = iocontroller.LocatePages(result)
		iocontroller.RenderResult(selectedPos, result, pageList[:])
	}
}
