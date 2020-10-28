package iocontrol

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aman/modules"
	"github.com/nsf/termbox-go"
)

/*
 * height  ウィンドウの高さ
 * width   ウィンドウの幅
 * page    現在のページ番号 定義域は[0, maxPage]
 * maxPage 最大ページ番号
 */
type IoController struct {
	height  int
	width   int
	page    int
	maxPage int
	query   string
}

/*
 * @param manLists オプションとオプション説明が格納された文字列と、各オプション説明の行数の配列
 * @description IoControllerのコンストラクタ
 */
func NewIoController(manLists []modules.ManData) *IoController {
	width, height := termbox.Size()
	iocontroller := IoController{
		height:  height,
		width:   width,
		page:    0,
		maxPage: 0,
		query:   "",
	}
	return &iocontroller
}

func (iocontroller *IoController) GetQuery() string {
	return iocontroller.query
}

func (iocontroller *IoController) DeleteInput() {
	var space = ""
	for i := 0; i < len(iocontroller.query); i++ {
		space += " "
	}
	fmt.Printf("\r%s", space)
	if 0 < len(iocontroller.query) {
		iocontroller.query = (iocontroller.query)[:len(iocontroller.query)-1]
	}
}

func (iocontroller *IoController) ReceiveKeys(selectedPos *int) int {
	var ev termbox.Event = termbox.PollEvent()

	if ev.Type != termbox.EventKey {
		return 0
	}

	switch ev.Key {
	case termbox.KeyEsc:
		return -1
	case termbox.KeyArrowUp:
		return 90
	case termbox.KeyArrowDown:
		return 91
	case termbox.KeyArrowRight:
		iocontroller.page++
		if iocontroller.maxPage < iocontroller.page {
			iocontroller.page = iocontroller.maxPage
		}
	case termbox.KeyArrowLeft:
		iocontroller.page--
		if iocontroller.page < 0 {
			iocontroller.page = 0
		}
	case termbox.KeySpace:
		iocontroller.query += " "
		break
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		iocontroller.DeleteInput()
		break
	case termbox.KeyEnter:
		return 99
	default:
		iocontroller.page = 0
		*selectedPos = 0
		iocontroller.query += string(ev.Ch)
		break
	}
	return 1
}

func (iocontroller *IoController) RenderQuery() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
	fmt.Printf("\r> %s", iocontroller.query)
}

func (iocontroller *IoController) RenderPageNumber() {
	var pageNumberText string = strconv.Itoa(iocontroller.page+1) + "/" + strconv.Itoa(iocontroller.maxPage+1)
	var blankCounts int = iocontroller.width - len("> ") - len(iocontroller.query) - len(pageNumberText)
	var blanks string = strings.Repeat(" ", blankCounts)
	fmt.Printf("%s\n", blanks + pageNumberText)
}

func (iocontroller *IoController) RenderOptionStack(command []string, stackOptions []string) {
	for i := 0; i < len(command); i++ {
		fmt.Printf("%s ", command[i])
	}
	if 0 < len(stackOptions) {
		for i := 0; i < len(stackOptions); i++ {
			fmt.Printf("%s ", stackOptions[i])
		}
	}
	fmt.Println("")
}

func (iocontroller *IoController) RenderResult(selectedPos int, result []modules.ManData, pageList []int) {
	const SEPARATOR = "----------"
	// rowは、表示行数を表す。
	// query行と先頭SEPARATORの2行分
	var row = 2
	fmt.Println(SEPARATOR)

	if len(result) == 0 {
		return
	}

	for i := pageList[iocontroller.page]; i < pageList[iocontroller.page+1]; i++ {
		row += strings.Count(result[i].Contents, "\n") + 2
		if iocontroller.height <= row {
			return
		}
		var state string = "\r%s\n"
		if selectedPos == i {
			// 選択行だけ赤色に変更
			state = "\r\x1b[31m%s\x1b[0m\n"
		}
		fmt.Printf(state, result[i].Contents)
		fmt.Println(SEPARATOR)
	}
}

/*
 * @param manLists オプションとオプション説明が格納された文字列と、各オプション説明の行数の配列
 * @description 各ページの先頭となるオプション配列manListsのindex番号が格納された配列を生成する
 */
func (iocontroller *IoController) LocatePages(manLists []modules.ManData) []int {
	var maxLineNumber = -1
	pageList := []int{0}
	// query行、option stack行、SEPARATORの３行
	var lineCount = 3
	var page = 0
	iocontroller.maxPage = 0

	for i := 0; i < len(manLists); i++ {
		// for文を抜けた後に、ウィンドウの高さが低すぎて描画できないかを判定するために、
		// 一番行数の多いオプション説明文の行数を求める
		if maxLineNumber < manLists[i].LineNumber {
			maxLineNumber = manLists[i].LineNumber
		}

		// ウィンドウの高さをオーバーしてしまう場合、次のページにオプション説明を表示する
		if iocontroller.height < lineCount+manLists[i].LineNumber {
			lineCount = 2
			page++
			pageList = append(pageList, i)
			if iocontroller.maxPage < page {
				iocontroller.maxPage = page
			}
		}

		lineCount += manLists[i].LineNumber

		if i == len(manLists)-1 {
			pageList = append(pageList, i+1)
		}
	}

	// 2は、query行と先頭SEPARATORの2行分
	if iocontroller.height < maxLineNumber + 2 {
		panic(errors.New("Window height is too small"))
	}

	return pageList
}
