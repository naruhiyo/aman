package ipagination

import (
	"errors"

	spagination "github.com/naruhiyo/aman/struct/pagination"
)

type PaginationStruct spagination.Pagination

/*
 * コンストラクタ
 */
func NewPagination() *PaginationStruct {
	return &PaginationStruct{
		Page:        0,
		MaxPage:     0,
		SelectedPos: 0,
		PageList:    []int{},
	}
}

/*
 * @description 各ページの先頭となるオプション配列manListsのindex番号が格納された配列を生成する
 * @param lineNumnbers 各オプション説明の行数の配列
 * @param windowHeight ウィンドウの高さ
 */
func (myself *PaginationStruct) LocatePages(lineNumnbers []int, windowHeight int) {
	var maxLineNumber = -1
	// query行、option stack行、SEPARATORの３行
	var lineCount = 3
	var page = 0
	myself.MaxPage = 0
	myself.PageList = []int{0}

	for i := 0; i < len(lineNumnbers); i++ {
		// for文を抜けた後に、ウィンドウの高さが低すぎて描画できないかを判定するために、
		// 一番行数の多いオプション説明文の行数を求める
		if maxLineNumber < lineNumnbers[i] {
			maxLineNumber = lineNumnbers[i]
		}

		// ウィンドウの高さをオーバーしてしまう場合、次のページにオプション説明を表示する
		if windowHeight < lineCount+lineNumnbers[i] {
			lineCount = 2
			page++
			myself.PageList = append(myself.PageList, i)
			if myself.MaxPage < page {
				myself.MaxPage = page
			}
		}

		lineCount += lineNumnbers[i]

		if i == len(lineNumnbers)-1 {
			myself.PageList = append(myself.PageList, i+1)
		}
	}

	// 2は、query行と先頭SEPARATORの2行分
	if windowHeight < maxLineNumber+2 {
		panic(errors.New("Window height is too small"))
	}
}

/*
 * @description 次の行へ遷移
 */
func (myself *PaginationStruct) NextLine(maxLength int) {
	if myself.SelectedPos < maxLength {
		myself.SelectedPos++
	}
}

/*
 * @description 前の行へ遷移
 * @param maxLength 最大行数
 */
func (myself *PaginationStruct) BackLine() {
	if myself.SelectedPos > 0 {
		myself.SelectedPos--
	}
}

/*
 * @description 次のページへ遷移
 */
func (myself *PaginationStruct) NextPage() {
	myself.Page++
	if myself.MaxPage < myself.Page {
		myself.Page = myself.MaxPage
	}
	myself.SelectedPos = myself.PageList[myself.Page]
}

/*
 * @description 前のページへ遷移
 */
func (myself *PaginationStruct) BackPage() {
	myself.Page--
	if myself.Page < 0 {
		myself.Page = 0
	}
	myself.SelectedPos = myself.PageList[myself.Page]
}

/*
 * @description リセット
 */
func (myself *PaginationStruct) Reset() {
	myself.Page = 0
	myself.SelectedPos = 0
}
