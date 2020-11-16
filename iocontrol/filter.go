package iocontrol

import (
	"strings"
)

/*
@param inputs   クエリ
@param manLists オプションとオプション説明が格納された文字列と、各オプション説明の行数の配列
@description
1. クエリを空白類で区切って配列化する
2. 区切ったクエリを1要素ごとに取り出す
3. 取り出したクエリが、オプション説明文字列の部分文字列なら、次回に取り出すクエリに対する検索対象として、オプション説明文字列を配列に格納する
4. 区切ったクエリをすべて取り出し終えるか、次回の検索対象のオプション説明文字列が無くなるまで2.と3.を繰り返す
*/
func IncrementalSearch(inputs string, manLists []ManData) []ManData {
	// クエリを空白区切りで取得
	separatedQuery := strings.Fields(inputs)
	result := manLists

	for indexQuery := 0; indexQuery < len(separatedQuery); indexQuery++ {
		resultCandidate := []ManData{}
		for indexResult := 0; indexResult < len(result); indexResult++ {
			// クエリの取り出し
			if 0 <= strings.Index(result[indexResult].Contents, separatedQuery[indexQuery]) {
				// クエリと一致する場合は結果に追加
				resultCandidate = append(resultCandidate, ManData{
					Contents:   result[indexResult].Contents,
					LineNumber: result[indexResult].LineNumber,
				})
			}
		}
		result = resultCandidate
		if len(result) == 0 {
			break
		}
	}
	return result
}
