package iwindow

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"

	imodel "aman/implement/model"
	swindow "aman/struct/window"
)

const (
	SELECTED_TEXT_COLOR = 238
	SELECTED_BG_COLOR   = 160
	MATCHED_TEXT_COLOR  = 200
	SEPARATOR           = "----------"
)

type WindowInfoStruct swindow.WindowInfo

/*
 * コンストラクタ作成
 */
func NewWindowInfo() *WindowInfoStruct {
	width, height := termbox.Size()
	return &WindowInfoStruct{
		Width:     width,
		Height:    height,
		TextColor: termbox.ColorDefault,
		BgColor:   termbox.ColorDefault,
	}
}

/*
 * @description ページ番号を描画する
 * @param page ページ番号
 * @param maxPage 最大ページ件数
 * @param query クエリ
 */
func (myself *WindowInfoStruct) RenderPageNumber(page int, maxPage int, query string) {
	var pageNumberText string = strconv.Itoa(page+1) + "/" + strconv.Itoa(maxPage+1)
	var blankCounts int = myself.Width - len("> ") - len(query) - len(pageNumberText)
	var blanks string = strings.Repeat(" ", blankCounts)
	myself.renderTextLine(len("> ")+len(query), 0, blanks+pageNumberText)
}

/*
 * @description 選択しているオプション一覧を描画する
 * @param command 入力コマンド
 * @param stackOptions 選択しているオプション
 */
func (myself *WindowInfoStruct) RenderOptionStack(commands []string, options []string) {
	var optionStack = ""
	for i := 0; i < len(commands); i++ {
		optionStack += commands[i] + " "
	}
	if 0 < len(options) {
		for i := 0; i < len(options); i++ {
			optionStack += options[i] + " "
		}
	}
	myself.renderTextLine(0, 1, optionStack)
}

/*
 * @description 検索結果を表示する
 * @param pageNum ページの件数
 * @param nextPageNum 次ページの件数
 * @param selectedPos プションの選択位置
 * @param list 抽出結果
 * @param query クエリ
 */
func (myself *WindowInfoStruct) RenderResult(pageNum int, nextPageNum int, selectedPos int, list *imodel.ManDataObjectStruct, query string) {
	myself.TextColor = termbox.ColorDefault
	myself.BgColor = termbox.ColorDefault
	// startLineは、次に表示する行の行番号(0スタート)を表す。
	// iocontroller.RenderTextLine()が呼ばれた後にインクリメントする
	// query行、選択オプション行の2行分が既に表示されているので初期値は2
	var startLine = 2

	myself.renderTextLine(0, 2, SEPARATOR)
	startLine++

	if len(list.Filtered) == 0 {
		return
	}

	for i := pageNum; i < nextPageNum; i++ {
		// Contentsの最終行がターミナルの最終行まで表示可能かどうかを判定している
		// iocontroller.heightは1スタート、startLineは0スタート
		// どちらも単位はターミナル上での1行
		if myself.Height < startLine+strings.Count(list.Filtered[i].Contents, "\n") {
			return
		}
		if selectedPos == i {
			// 選択行だけハイライト
			myself.TextColor = SELECTED_TEXT_COLOR
			myself.BgColor = SELECTED_BG_COLOR
		}
		var contentsLines []string = strings.Split(list.Filtered[i].Contents, "\n")
		for line := 0; line < len(contentsLines); line++ {
			// texts内でqueryが部分文字列として一致する先頭index番号の配列
			myself.searchMatchedText(contentsLines[line], query, list)
			myself.renderColoredTextLine(0, startLine, contentsLines[line], list)
			startLine++
		}
		myself.TextColor = termbox.ColorDefault
		myself.BgColor = termbox.ColorDefault
		myself.renderTextLine(0, startLine, SEPARATOR)
		startLine++
	}
}

/*
 * @description 一致するテキストの色を変更しつつ、テキストの描画（標準出力）を行う
 * @param x テキストの出現位置(x座標)
 * @param y テキストの出現位置(y座標)
 * @param texts 描画されるテキスト文字列
 * @param list 抽出結果
 */
func (myself *WindowInfoStruct) renderColoredTextLine(x, y int, texts string, list *imodel.ManDataObjectStruct) {
	// 注目したいmatchedIndexesのindex番号
	var textsRune = []rune(texts)
	var matchedFg termbox.Attribute = MATCHED_TEXT_COLOR

	for index := 0; index < len(textsRune); index++ {
		if 0 < len(list.Matched) {
			var targetIndex = myself.getTargetIndex(index, list)
			if targetIndex != -1 {
				for _, qr := range list.Matched[targetIndex].Text {
					termbox.SetCell(x, y, qr, matchedFg, myself.BgColor)
					x += runewidth.RuneWidth(qr)
				}
				index += len(list.Matched[targetIndex].Text) - 1
				continue
			}
		}
		termbox.SetCell(x, y, textsRune[index], myself.TextColor, myself.BgColor)
		x += runewidth.RuneWidth(textsRune[index])
	}
}

/*
 * @description 入力しているクエリを描画する
 * @param query クエリ
 */
func (myself *WindowInfoStruct) RenderQuery(query string) {
	myself.renderTextLine(0, 0, "> "+query)
}

/*
 * @description カーソルを描画する
 * @param cursorPosX カーソルの位置
 */
func (myself *WindowInfoStruct) RenderCursor(cursorPosX int) {
	termbox.SetCursor(cursorPosX, 0)
}

/*
 * @description originalText内に空白で区切られたqueryが、部分文字列として一致する
 *        			先頭のindex番号及び一致したqueryのMatchedInfo配列を求める
 *  1. iocontroller.queryを空白区切りに分割し、separetedQueryに格納する
 *  2. separetedQueryの各要素(query)に対し、targetText内にqueryが部分文字列として存在するかチェックする
 *  2. 部分文字列なら、先頭のindexをmatchedInfo.index, MatchedInfo.textをqueryとして、appendする
 *  3. 2.で一致したindexの次の文字以降をtargetTextとして更新し、1.に戻る。
 *     targetText内に全queryが存在しなくなるまで繰り返す。
 * @param originalText オプション説明文
 * @param query クエリ
 * @param list 抽出結果
 * @example originalText: "hoge hogera", query: "og a"の場合、
 *          matchedInfos: { MatchedInfo{ text: "og", index: 1 },
 *                          MatchedInfo{ text: "og", index: 6 },
 *                          MatchedInfo{ text: "a", index: 10 },
 *                        }
 */
func (myself *WindowInfoStruct) searchMatchedText(originalText string, query string, list *imodel.ManDataObjectStruct) {
	// 初期化
	list.Matched = nil
	// 探索文字列
	var separatedQuery = strings.Fields(query)
	if 0 < len(query) {
		for _, q := range separatedQuery {
			var startIndex = 0
			var targetText = originalText
			for {
				var targetTextLower = strings.ToLower(targetText)
				var qLower = strings.ToLower(q)
				var matchedIndex = strings.Index(targetTextLower, qLower)
				if matchedIndex == -1 {
					break
				}
				startIndex += utf8.RuneCountInString(targetText[:matchedIndex])
				list.Matched = append(
					list.Matched,
					list.GetMatchedInfo(
						targetText[matchedIndex:matchedIndex+len(q)],
						startIndex,
					),
				)
				startIndex += 1
				targetText = string([]rune(originalText)[startIndex:])
			}
		}
	}
}

/*
 * @description テキストの描画（標準出力）を行う
 * @param x テキストの出現位置(x座標)
 * @param y テキストの出現位置(y座標)
 * @param texts 描画されるテキスト文字列
 */
func (myself *WindowInfoStruct) renderTextLine(x, y int, texts string) {
	for _, text := range texts {
		termbox.SetCell(x, y, text, myself.TextColor, myself.BgColor)
		x += runewidth.RuneWidth(text)
	}
}

/*
 * @description インデックスを取得
 * @param index 選択されたインデックス
 * @param list 抽出結果
 * @return インデックス
 */
func (myself *WindowInfoStruct) getTargetIndex(index int, list *imodel.ManDataObjectStruct) int {
	for targetIndex, matchedInfo := range list.Matched {
		if index == matchedInfo.Index {
			return targetIndex
		}
	}
	return -1
}
