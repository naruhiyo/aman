package main

/**
 * m*** : 実装モジュール
 * s*** : 構造体モジュール
 */
import (
	mio "github.com/aman/modules/io"
	mmodel "github.com/aman/modules/model"
	mpagination "github.com/aman/modules/pagination"
	mutil "github.com/aman/modules/util"
	mwindow "github.com/aman/modules/window"
	"github.com/nsf/termbox-go"
)

/**
 * 描画処理
 */
func render(input *mio.InputStruct, list *mmodel.ListStruct, pagination *mpagination.PaginationStruct, window *mwindow.WindowInfoStruct) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// ページネーション設定
	pagination.LocatePages(list.MapLineNumber(), window.Height)
	// 基本Viewの描画
	window.RenderQuery(input.Query)
	window.RenderCursor(input.CursorPosX)
	window.RenderOptionStack(input.Commands, input.Options)
	window.RenderPageNumber(pagination.Page, pagination.MaxPage, input.Query)

	var pageNum int = pagination.PageList[pagination.Page]
	var nextPageNum int = pagination.PageList[pagination.Page+1]

	// 結果データの描画
	window.RenderResult(pageNum, nextPageNum, pagination.SelectedPos, list, input.Query)
	termbox.Flush()
}

func main() {
	// 標準入力有効化
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetOutputMode(termbox.Output256)

	// 初期化
	var windowInfo *mwindow.WindowInfoStruct = mwindow.NewWindowInfo()
	var pagination *mpagination.PaginationStruct = mpagination.NewPagination()
	var input *mio.InputStruct = mio.NewInput()
	var list *mmodel.ListStruct = mmodel.NewList()
	var command *mutil.CommandStruct = mutil.NewCommand()

	// コマンド実行
	command.ExecMan(input.Commands)
	list.AnalyzeMan(command.ManResult)

	// 描画処理 & ページネーション
	render(input, list, pagination, windowInfo)

	// ユーザーからの入力を受け付ける
	//   - キーボード入力一回ごとにループ実行させることでインタラクティブな処理を実現
	//   - `ESC`キーで処理を終了
loop:
	for {
		var ev termbox.Event = termbox.PollEvent()

		if ev.Type != termbox.EventKey {
			break loop
		}

		// 毎回 man 結果に対して検索を行う
		switch ev.Key {
		case termbox.KeyEsc:
			break loop
		case termbox.KeyArrowUp:
			pagination.NextLine()
		case termbox.KeyArrowDown:
			var maxLength int = len(list.Filtered) - 1
			pagination.BackLine(maxLength)
		case termbox.KeyArrowRight:
			pagination.NextPage()
		case termbox.KeyArrowLeft:
			pagination.BackPage()
		case termbox.KeyEnter:
			var pos int = pagination.SelectedPos
			var contents string = list.Filtered[pos].Contents
			input.ExtractOption(contents)
		case termbox.KeySpace:
			input.PutSpace()
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			input.DeleteInput()
			list.IncrementalSearch(input.Query)
		default:
			pagination.Reset()
			input.PutKey(ev)
			list.IncrementalSearch(input.Query)
		}
		// 描画処理 & ページネーション
		render(input, list, pagination, windowInfo)
	}

	// deferを利用すると 全ての処理が終わった後に呼ばれる
	// termbox を先に終了しておかないとコマンドプロンプト上に標準出力されない
	termbox.Close()

	// コマンド実行
	command.CmdOutput(input.Commands, input.Options)
}
