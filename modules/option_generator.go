package modules

import "strings"

/*
 * 選択した行のオプションを抽出する
 */
func ExtractOption(line string) string {
	// 文字列を空白区切で区切ったものの先頭がオプションのはずなのでそれを取得
	return strings.Split(line, " ")[0]
}
