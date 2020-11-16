package iocontrol

import (
	"strconv"
	"strings"
	"unicode/utf8"

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
 * @brief 空白で区切られたどのqueryが探索文字列内のどのindex番号から
 *        開始する部分文字列なのかを表す
 * text  空白で区切られたqueryの一要素
 * index 探索文字列内でtextが部分文字列として一致するindex番号の先頭
 */
type MatchedInfo struct {
	text  string
	index int
}

/*
 * @param manLists オプションとオプション説明が格納された文字列と、各オプション説明の行数の配列
 * @description IoControllerのコンストラクタ
 */
func NewIoController(manLists []ManData) *IoController {
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

/*
 * @description クエリ文字列の取得
 */
func (iocontroller *IoController) GetQuery() string {
	return iocontroller.query
}

/*
 * @description テキストの描画（標準出力）を行う
 * @param x テキストの出現位置(x座標)
 * @param y テキストの出現位置(y座標)
 * @param texts 描画されるテキスト文字列
 * @param fg, bg 描画時の色（テキストと背景色）
 */
func (iocontroller *IoController) RenderTextLine(x, y int, texts string, fg, bg termbox.Attribute) {
	for _, r := range texts {
		termbox.SetCell(x, y, r, fg, bg)
		x += runewidth.RuneWidth(r)
	}
}

/*
 * @brief originalText内に空白で区切られたqueryが、部分文字列として一致する
 *        先頭のindex番号及び一致したqueryのMatchedInfo配列を求める
 * @param originalText オプション説明文
 * @example originalText: "hoge hogera", query: "og a"の場合、
 *          matchedInfos: { MatchedInfo{ text: "og", index: 1 },
 *                          MatchedInfo{ text: "og", index: 6 },
 *                          MatchedInfo{ text: "a", index: 10 },
 *                        }
 * @description
 *  1. iocontroller.queryを空白区切りに分割し、separetedQueryに格納する
 *  2. separetedQueryの各要素(query)に対し、targetText内にqueryが部分文字列として存在するかチェックする
 *  2. 部分文字列なら、先頭のindexをmatchedInfo.index, MatchedInfo.textをqueryとして、appendする
 *  3. 2.で一致したindexの次の文字以降をtargetTextとして更新し、1.に戻る。
 *     targetText内に全queryが存在しなくなるまで繰り返す。
 */
func (iocontroller *IoController) GetMatchedInfos(originalText string) []MatchedInfo {
	// 探索文字列
	var matchedInfos []MatchedInfo
	var separatedQuery = strings.Fields(iocontroller.query)
	if 0 < len(iocontroller.query) {
		for _, query := range separatedQuery {
			var startIndex = 0
			var targetText = originalText
			for {
				var matchedIndex = strings.Index(targetText, query)
				if matchedIndex == -1 {
					break
				}

				matchedInfos = append(matchedInfos, MatchedInfo{
					text:  query,
					index: startIndex + utf8.RuneCountInString(targetText[:matchedIndex]),
				})
				startIndex += utf8.RuneCountInString(targetText[:matchedIndex]) + 1
				targetText = string([]rune(originalText)[startIndex:])
			}
		}
	}

	return matchedInfos
}

/*
 * @description 一致するテキストの色を変更しつつ、テキストの描画（標準出力）を行う
 * @param x テキストの出現位置(x座標)
 * @param y テキストの出現位置(y座標)
 * @param texts 描画されるテキスト文字列
 * @param fg, bg 描画時の色（テキストと背景色）
 */
func (iocontroller *IoController) RenderColoredTextLine(x, y int, texts string, fg, bg termbox.Attribute) {
	// texts内でqueryが部分文字列として一致する先頭index番号の配列
	var matchedInfos []MatchedInfo = iocontroller.GetMatchedInfos(texts)
	// 注目したいmatchedIndexesのindex番号
	var textsRune = []rune(texts)
	var matchedFg termbox.Attribute = 200
	var getTargetIndex = func(index int) int {
		for targetIndex, matchedInfo := range matchedInfos {
			if index == matchedInfo.index {
				return targetIndex
			}
		}
		return -1
	}

	for index := 0; index < len(textsRune); index++ {
		if 0 < len(matchedInfos) {
			var targetIndex = getTargetIndex(index)
			if targetIndex != -1 {
				for _, qr := range matchedInfos[targetIndex].text {
					termbox.SetCell(x, y, qr, matchedFg, bg)
					x += runewidth.RuneWidth(qr)
				}
				index += len(matchedInfos[targetIndex].text) - 1
				continue
			}
		}
		termbox.SetCell(x, y, textsRune[index], fg, bg)
		x += runewidth.RuneWidth(textsRune[index])
	}
}

/*
 * @description 入力を削除する
 */
func (iocontroller *IoController) DeleteInput() {
	if 0 < len(iocontroller.query) {
		iocontroller.query = string([]rune(iocontroller.query)[:utf8.RuneCountInString(iocontroller.query)-1])
	}
}

/*
 * @description キー入力を受け付ける
 * @param オプションの選択位置
 */
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

/*
 * @description 入力しているクエリを描画する
 */
func (iocontroller *IoController) RenderQuery() {
	iocontroller.RenderTextLine(0, 0, "> "+iocontroller.query, termbox.ColorDefault, termbox.ColorDefault)
}

/*
 * @description カーソルを描画する
 */
func (iocontroller *IoController) RenderCursor() {
	termbox.SetCursor(iocontroller.cursorPosX, 0)
}

/*
 * @description ページ番号を描画する
 */
func (iocontroller *IoController) RenderPageNumber() {
	var pageNumberText string = strconv.Itoa(iocontroller.page+1) + "/" + strconv.Itoa(iocontroller.maxPage+1)
	var blankCounts int = iocontroller.width - len("> ") - len(iocontroller.query) - len(pageNumberText)
	var blanks string = strings.Repeat(" ", blankCounts)
	iocontroller.RenderTextLine(len("> ")+len(iocontroller.query), 0, blanks+pageNumberText, termbox.ColorDefault, termbox.ColorDefault)
}

/*
 * @description 選択しているオプション一覧を描画する
 * @param command 入力コマンド
 * @param stackOptions 選択しているオプション
 */
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

/*
 * @description 検索結果を表示する
 * @param selectedPos プションの選択位置
 * @param result 抽出結果
 * @param pageList 各ページの先頭となるmanListのindex番号
 */
func (iocontroller *IoController) RenderResult(selectedPos int, result []ManData, pageList []int) {
	const SEPARATOR = "----------"
	var separatorFg, separatorBg termbox.Attribute = termbox.ColorDefault, termbox.ColorDefault
	// startLineは、次に表示する行の行番号(0スタート)を表す。
	// iocontroller.RenderTextLine()が呼ばれた後にインクリメントする
	// query行、選択オプション行の2行分が既に表示されているので初期値は2
	var startLine = 2

	iocontroller.RenderTextLine(0, 2, SEPARATOR, separatorFg, separatorBg)
	startLine++

	if len(result) == 0 {
		return
	}

	for i := pageList[iocontroller.page]; i < pageList[iocontroller.page+1]; i++ {
		var contentsFg, contentsBg termbox.Attribute = termbox.ColorDefault, termbox.ColorDefault
		// Contentsの最終行がターミナルの最終行まで表示可能かどうかを判定している
		// iocontroller.heightは1スタート、startLineは0スタート
		// どちらも単位はターミナル上での1行
		if iocontroller.height < startLine+strings.Count(result[i].Contents, "\n") {
			return
		}
		if selectedPos == i {
			// 選択行だけハイライト
			contentsFg = 238
			contentsBg = 160
		}
		var contentsLines []string = strings.Split(result[i].Contents, "\n")
		for line := 0; line < len(contentsLines); line++ {
			iocontroller.RenderColoredTextLine(0, startLine, contentsLines[line], contentsFg, contentsBg)
			startLine++
		}
		iocontroller.RenderTextLine(0, startLine, SEPARATOR, separatorFg, separatorBg)
		startLine++
	}
}
