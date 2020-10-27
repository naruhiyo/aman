package filter

import (
	"strings"

	"github.com/aman/modules"
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
func IncrementalSearch(inputs *string, manLists []modules.ManData) []modules.ManData {
	separatedQuery := strings.Fields(*inputs)
	result := manLists

	for indexQuery := 0; indexQuery < len(separatedQuery); indexQuery++ {
		resultCandidate := []modules.ManData{}
		for indexResult := 0; indexResult < len(result); indexResult++ {
			if 0 <= strings.Index(result[indexResult].Contents, separatedQuery[indexQuery]) {
				resultCandidate = append(resultCandidate, modules.ManData{
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
