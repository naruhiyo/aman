package mio

import (
	"errors"
	"flag"
	"strings"
	"unicode/utf8"

	sio "github.com/aman/struct/io"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type InputStruct sio.InputStruct

/*
 * コンストラクタ
 */
func NewInput() *InputStruct {
	input := &InputStruct{
		Commands:   []string{},
		Options:    []string{},
		Query:      "",
		CursorPosX: 2,
	}
	input.Parse()
	return input
}

/**
* @description コマンドライン引数を取得
 */
func (myself *InputStruct) Parse() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		panic(errors.New("Error: No arguments"))
	}

	myself.Commands = args
}

/*
 * @description 入力を削除する
 */
func (myself *InputStruct) DeleteInput() {
	if 0 < len(myself.Query) {
		myself.CursorPosX -= runewidth.RuneWidth([]rune(myself.Query)[utf8.RuneCountInString(myself.Query)-1])
		myself.Query = string([]rune(myself.Query)[:utf8.RuneCountInString(myself.Query)-1])
	}
	if myself.CursorPosX < 2 {
		myself.CursorPosX = 2
	}
}

/*
 * @description 選択した行のオプションを抽出する
 * @param line オプション説明文
 */
func (myself *InputStruct) ExtractOption(line string) {
	// 文字列を空白区切で区切ったものの先頭がオプションのはずなのでそれを取得
	var selectedOption string = strings.Split(line, " ")[0]
	// 末端の改行を削除する
	selectedOption = strings.TrimRight(selectedOption, "\n")

	// 重複選択を制限する
	for _, option := range myself.Options {
		// 一致するオプションが見つかったら追加処理を行わず、returnする
		if option == selectedOption {
			return
		}
	}

	// ストック
	myself.Options = append(myself.Options, selectedOption)
}

/*
 * @description スペース入力
 */
func (myself *InputStruct) PutSpace() {
	myself.Query += " "
	myself.CursorPosX++
}

/*
 * @description キー入力
 */
func (myself *InputStruct) PutKey(ev termbox.Event) {
	myself.Query += string(ev.Ch)
	for _, r := range string(ev.Ch) {
		myself.CursorPosX += runewidth.RuneWidth(r)
	}
}
