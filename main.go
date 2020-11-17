package main

import (
	"github.com/aman/iocontrol"
	"github.com/aman/util"
	"github.com/nsf/termbox-go"
)

const (
	ESCAPE     = -1
	ANYKEY     = 1
	ARROW_UP   = 90
	ARROW_DOWN = 91
	ENTER      = 99
)

func main() {
	// 標準入力有効化
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetOutputMode(termbox.Output256)

	// 引数取得
	var args []string = iocontrol.Parse()

	// コマンド実行
	var commandResult string = util.ExecMan(args)

	var manLists []iocontrol.ManData = iocontrol.AnalyzeMan(commandResult)
	// 選択位置
	var selectedPos int = 0
	// 検索結果
	var result []iocontrol.ManData = manLists
	// 選択したオプション格納
	var stackOptions []string

	// 初期化
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	iocontroller := iocontrol.NewIoController(result)
	iocontroller.RenderQuery()
	iocontroller.RenderCursor()
	pageList := iocontroller.LocatePages(result)
	iocontroller.RenderPageNumber()
	iocontroller.RenderOptionStack(args, stackOptions)
	iocontroller.RenderResult(selectedPos, result, pageList[:])
	termbox.Flush()

	// ユーザーからの入力を受け付ける
	//   - キーボード入力一回ごとにループ実行させることでインタラクティブな処理を実現
	//   - `ESC`キーで処理を終了
loop:
	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		var keyStatus int = iocontroller.ReceiveKeys(&selectedPos)

		switch keyStatus {
		// 毎回 man 結果に対して検索を行う
		case ANYKEY:
			result = iocontrol.IncrementalSearch(iocontroller.GetQuery(), manLists)
		case ARROW_UP:
			if selectedPos > 0 {
				selectedPos--
			}
		case ARROW_DOWN:
			if selectedPos < len(result)-1 {
				selectedPos++
			}
		case ENTER:
			var option string = iocontrol.ExtractOption(result[selectedPos].Contents)
			stackOptions = iocontrol.DistinctOption(option, stackOptions)
		case ESCAPE:
			break loop
		}
		// 描画処理 & ページネーション
		iocontroller.RenderQuery()
		iocontroller.RenderCursor()
		pageList = iocontroller.LocatePages(result)
		iocontroller.RenderPageNumber()
		iocontroller.RenderOptionStack(args, stackOptions)
		iocontroller.RenderResult(selectedPos, result, pageList[:])
		termbox.Flush()
	}

	// deferを利用すると 全ての処理が終わった後に呼ばれる
	// termbox を先に終了しておかないとコマンドプロンプト上に標準出力されない
	termbox.Close()

	// コマンド実行
	util.CmdOutput(args, stackOptions)
}
