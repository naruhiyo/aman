package iocontrol

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/aman/modules"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

/*
 * height     ウィンドウの高さ
 * width      ウィンドウの幅
 * page       現在のページ番号 定義域は[0, maxPage]
 * maxPage    最大ページ番号
 * query      検索クエリ
 * cursorPosX カーソルのx座標
 */
type IoController struct {
	height     int
	width      int
	page       int
	maxPage    int
	query      string
	cursorPosX int
}

/*
 * @param manLists オプションとオプション説明が格納された文字列と、各オプション説明の行数の配列
 * @description IoControllerのコンストラクタ
 */
func NewIoController(manLists []modules.ManData) *IoController {
	width, height := termbox.Size()
	iocontroller := IoController{
		height:     height,
		width:      width,
		page:       0,
		maxPage:    0,
		query:      "",
		cursorPosX: 2,
	}
	return &iocontroller
}

func (iocontroller *IoController) GetQuery() string {
	return iocontroller.query
}

func (iocontroller *IoController) RenderTextLine(x, y int, texts string, fg, bg termbox.Attribute) {
	for _, r := range texts {
		termbox.SetCell(x, y, r, fg, bg)
		x += runewidth.RuneWidth(r)
	}
}

func (iocontroller *IoController) DeleteInput() {
	if 0 < len(iocontroller.query) {
		iocontroller.query = string([]rune(iocontroller.query)[:utf8.RuneCountInString(iocontroller.query)-1])
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
		iocontroller.cursorPosX++
		break
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if 0 < len(iocontroller.query) {
			iocontroller.cursorPosX -= runewidth.RuneWidth([]rune(iocontroller.query)[utf8.RuneCountInString(iocontroller.query)-1])
		}
		if iocontroller.cursorPosX < 2 {
			iocontroller.cursorPosX = 2
		}
		iocontroller.DeleteInput()
		break
	case termbox.KeyEnter:
		return 99
	default:
		iocontroller.page = 0
		*selectedPos = 0
		iocontroller.query += string(ev.Ch)
		for _, r := range string(ev.Ch) {
			iocontroller.cursorPosX += runewidth.RuneWidth(r)
		}
		break
	}
	return 1
}

func (iocontroller *IoController) RenderQuery() {
	iocontroller.RenderTextLine(0, 0, "> " + iocontroller.query, termbox.ColorDefault, termbox.ColorDefault)
}

func (iocontroller *IoController) RenderCursor() {
	termbox.SetCursor(iocontroller.cursorPosX, 0)
}

func (iocontroller *IoController) RenderPageNumber() {
	var pageNumberText string = strconv.Itoa(iocontroller.page+1) + "/" + strconv.Itoa(iocontroller.maxPage+1)
	var blankCounts int = iocontroller.width - len("> ") - len(iocontroller.query) - len(pageNumberText)
	var blanks string = strings.Repeat(" ", blankCounts)
	iocontroller.RenderTextLine(len("> ") + len(iocontroller.query), 0, blanks + pageNumberText, termbox.ColorDefault, termbox.ColorDefault)
}

func (iocontroller *IoController) RenderOptionStack(command []string, stackOptions []string) {
	var optionStack = ""
	for i := 0; i < len(command); i++ {
		optionStack += command[i] + " "
	}
	if 0 < len(stackOptions) {
		for i := 0; i < len(stackOptions); i++ {
			optionStack += stackOptions[i] + " "
		}
	}
	iocontroller.RenderTextLine(0, 1, optionStack, termbox.ColorDefault, termbox.ColorDefault)
}

func (iocontroller *IoController) RenderResult(selectedPos int, result []modules.ManData, pageList []int) {
	const SEPARATOR = "----------"
	var separatorFg, separatorBg termbox.Attribute = termbox.ColorDefault, termbox.ColorDefault

	iocontroller.RenderTextLine(0, 2, SEPARATOR, separatorFg, separatorBg)

	// rowは、表示行数を表す。
	// query行と先頭SEPARATORの2行分
	var row = 2

	if len(result) == 0 {
		return
	}

	for i := pageList[iocontroller.page]; i < pageList[iocontroller.page+1]; i++ {
		var contentsFg, contentsBg termbox.Attribute = termbox.ColorDefault, termbox.ColorDefault
		row += strings.Count(result[i].Contents, "\n") + 2
		if iocontroller.height <= row {
			return
		}
		if selectedPos == i {
			// 選択行だけ赤色に変更
			contentsFg = 167
			contentsBg = 160
		}
		var contentsLines []string = strings.Split(result[i].Contents, "\n")
		var startLine int = row - strings.Count(result[i].Contents, "\n") - 1
		for line := 0; line < len(contentsLines); line++ {
			iocontroller.RenderTextLine(0, startLine + line, contentsLines[line], contentsFg, contentsBg)
		}
		iocontroller.RenderTextLine(0, startLine + len(contentsLines), SEPARATOR, separatorFg, separatorBg)
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
