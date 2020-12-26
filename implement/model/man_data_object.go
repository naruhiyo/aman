package imodel

import (
	"bytes"
	"strings"

	smodel "github.com/naruhiyo/aman/struct/model"
	"github.com/nsf/termbox-go"
)

type ManDataObjectStruct smodel.ManDataObject

/*
 * コンストラクタ
 */
func NewManDataObject() *ManDataObjectStruct {
	return &ManDataObjectStruct{
		Original: []smodel.ManData{},
		Filtered: []smodel.ManData{},
		Matched:  []smodel.MatchedInfo{},
	}
}

/*
 * @description 	マッチ情報を構造体にセット
 * @params text 	テキスト
 * @params index インデックス
 * @return 			構造体
 */
func (myself *ManDataObjectStruct) GetMatchedInfo(text string, index int) smodel.MatchedInfo {
	return smodel.MatchedInfo{
		Text:  text,
		Index: index,
	}
}

/*
 * @description コマンド実行結果からオプションを抽出する
 * オプションの判定方法
 *   - `-` (または `--`) を検索する
 *   - 検索結果に対して`-`(または `--`) の出現位置(index)を計算する
 *   - 出現位置(index)の値だけ空白文字を生成し、オプションと結合する
 *   - 結合した値と元の値を比較し、一致すればオプションとみなす
 * @params manResult manコマンド実行結果
 */
func (myself *ManDataObjectStruct) AnalyzeMan(manResult string) {
	// === 条件 ===
	// ハイフンまたはダブルハイフンで始まる英単語
	var splitOutputs []string = strings.Split(manResult, "\n")

	// 検索フラグ
	// true: 検索中、false: 検索してない
	var isFinding bool = true

	// buffer の方が string結合より効率が良い
	var buffer bytes.Buffer // オプション説明のブロックを入れる変数
	var results []smodel.ManData

	// オプションを配列に追加する
	var count int = 0
	// オプションに必要な空白の個数
	var definedOptionBlankCount int = -1
	// 右端詰めのための空白の個数
	var paddingCounts int
	// 右端詰めのための空白
	var padding string
	width, _ := termbox.Size()
	for _, line := range splitOutputs {
		if isFinding {
			// オプションのヘッダーに来るまでスキップ
			if !myself.isOptionText(line) {
				continue
			}

			// オプションの前に何文字空白があるか計算
			count = myself.getOptionHeaderBlankCounts(line)

			// オプションの条件を満たしているか確認
			if !myself.isOptionHeaderText(line, count) {
				continue
			}

			// オプションが初めて見つかった時、空白の個数を記憶しておく
			if definedOptionBlankCount == -1 {
				definedOptionBlankCount = count
			}

			// 説明文の中にオプションの条件を満たす文があった場合、ヘッダーと判断してはダメなのでスキップ
			if definedOptionBlankCount != count {
				continue
			}

			isFinding = false
			paddingCounts = width - len(line[definedOptionBlankCount:])
			padding = strings.Repeat(" ", paddingCounts)
			buffer.WriteString(line[definedOptionBlankCount:] + padding)
		} else {
			// 改行だった場合次のオプションを探す
			if len(line) == 0 {
				results = append(results, smodel.ManData{
					Contents:   buffer.String(),
					LineNumber: strings.Count(buffer.String(), "\n") + 2,
				})
				buffer.Reset()
				isFinding = true
				continue
			}

			paddingCounts = width - len(line[definedOptionBlankCount:])
			padding = strings.Repeat(" ", paddingCounts)
			// バッファに文字列追加
			buffer.WriteString("\n" + line[definedOptionBlankCount:] + padding)
		}
	}

	// 最初は検索結果とオリジナルデータに同じものを挿入
	myself.Original = results
	myself.Filtered = results
}

/*
 * @param query   クエリ
 * @description
    1. クエリを空白類で区切って配列化する
    2. 区切ったクエリを1要素ごとに取り出す
    3. 取り出したクエリが、オプション説明文字列の部分文字列なら、次回に取り出すクエリに対する検索対象として、オプション説明文字列を配列に格納する
    4. 区切ったクエリをすべて取り出し終えるか、次回の検索対象のオプション説明文字列が無くなるまで2.と3.を繰り返す
*/
func (myself *ManDataObjectStruct) IncrementalSearch(query string) {
	// クエリを空白区切りで取得
	separatedQuery := strings.Fields(query)
	myself.Filtered = myself.Original

	for indexQuery := 0; indexQuery < len(separatedQuery); indexQuery++ {
		resultCandidate := []smodel.ManData{}
		for indexResult := 0; indexResult < len(myself.Filtered); indexResult++ {
			// クエリの取り出し
			if 0 <= strings.Index(myself.Filtered[indexResult].Contents, separatedQuery[indexQuery]) {
				// クエリと一致する場合は結果に追加
				resultCandidate = append(resultCandidate, smodel.ManData{
					Contents:   myself.Filtered[indexResult].Contents,
					LineNumber: myself.Filtered[indexResult].LineNumber,
				})
			}
		}
		myself.Filtered = resultCandidate
		if len(myself.Filtered) == 0 {
			break
		}
	}
}

/*
 * @description LineNumber の配列を返す
 * @return インデックスの配列
 */
func (myself *ManDataObjectStruct) MapLineNumber() []int {
	var result []int = []int{}
	for i := 0; i < len(myself.Filtered); i++ {
		result = append(result, myself.Filtered[i].LineNumber)
	}
	return result
}

/*
 * @description MatchedText の配列を返す
 * @return テキストの配列
 */
func (myself *ManDataObjectStruct) MapMatchedText() []string {
	var result []string = []string{}
	for i := 0; i < len(myself.Filtered); i++ {
		result = append(result, myself.Matched[i].Text)
	}
	return result
}

/*
 * @description オプション条件を満たしているかをチェック
 * @param line 文字列
 * @return 真偽値
 */
func (myself *ManDataObjectStruct) isOptionText(line string) bool {
	return strings.Contains(line, "-")
}

/*
 * @description isOptionHeaderText()内で生成する空白文字の文字数を求める
 * @param line 文字列
 * @return 文字数
 */
func (myself *ManDataObjectStruct) getOptionHeaderBlankCounts(line string) int {
	var count int = 0
	if strings.Contains(line, "-") {
		count = strings.Index(line, "-")
	}
	return count
}

/*
 * @description オプションのヘッダーであるかチェック（オプションの説明文にあるハイフンを弾く）
 * @param line 文字列
 * @param count 空白文字数
 */
func (myself *ManDataObjectStruct) isOptionHeaderText(line string, count int) bool {
	var blanks string = strings.Repeat(" ", count)

	// 先頭から文字を見たときにオプション条件を満たしているか確認する
	return strings.HasPrefix(line, blanks+"-")
}
