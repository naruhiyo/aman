package main

/**
 * i*** : 実装モジュール
 * s*** : 構造体モジュール
 */
import (
	iio "aman/implement/io"
	imodel "aman/implement/model"
	ipagination "aman/implement/pagination"
	iutil "aman/implement/util"
	iwindow "aman/implement/window"
	"fmt"

	"github.com/nsf/termbox-go"
)

/**
 * 描画処理
 */
func render(input *iio.InputStruct, list *imodel.ManDataObjectStruct, pagination *ipagination.PaginationStruct, window *iwindow.WindowInfoStruct) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// ページネーション設定
	pagination.LocatePages(list.MapLineNumber(), window.Height)
	// 基本Viewの描画
	window.RenderQuery(input.Query)
	window.RenderCursor(input.CursorPosX)
	window.RenderOptionStack(input.Commands, input.Options)
	window.RenderPageNumber(pagination.Page, pagination.MaxPage, input.Query)

	var pageNum int = pagination.PageList[pagination.Page]
	var nextPageNum int = pageNum

	if len(list.Filtered) > 0 {
		nextPageNum = pagination.PageList[pagination.Page+1]
	}

	// 結果データの描画
	window.RenderResult(pageNum, nextPageNum, pagination.SelectedPos, list, input.Query)
	termbox.Flush()
}

/**
 * main実行後の後処理
 */
func postExecMain() {
	if r := recover(); r != nil {
		var recoverCommand *iutil.CommandStruct = iutil.NewCommand()
		recoverCommand.ExecWithStdin("stty", "sane")
		fmt.Printf("Terminated with error: %v\n", r)
	}
}

func main() {
	// panic時には、端末設定をデフォルトに戻す
	defer postExecMain()
	// 標準入力有効化
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetOutputMode(termbox.Output256)

	// 初期化
	var windowInfo *iwindow.WindowInfoStruct = iwindow.NewWindowInfo()
	var pagination *ipagination.PaginationStruct = ipagination.NewPagination()
	var input *iio.InputStruct = iio.NewInput()
	var list *imodel.ManDataObjectStruct = imodel.NewManDataObject()
	var command *iutil.CommandStruct = iutil.NewCommand()

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
			pagination.BackLine()
		case termbox.KeyArrowDown:
			// 表示しているページ件数までしかカーソルを動かせないようにする
			var pageList []int = pagination.PageList
			if len(list.Filtered) > 0 {
				var maxLength int = pageList[pagination.Page+1] - 1
				pagination.NextLine(maxLength)
			}
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
