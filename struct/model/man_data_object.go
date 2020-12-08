package model

/*
 * Contents       オプション内容
 * LineNumber     行番号
 */
type ManData struct {
	Contents   string
	LineNumber int
}

/*
 * @brief 空白で区切られたどのqueryが探索文字列内のどのindex番号から
 *        開始する部分文字列なのかを表す
 * text  空白で区切られたqueryの一要素
 * index 探索文字列内でtextが部分文字列として一致するindex番号の先頭
 */
type MatchedInfo struct {
	Text  string
	Index int
}

/*
 * Original 	Man実行結果配列
 * Filtered		検索後のMan実行結果配列
 * Matched  	検索してマッチした結果を格納する配列
 * MatchedInfoStruct 構造体参照用
 */
type ManDataObject struct {
	Original          []ManData
	Filtered          []ManData
	Matched           []MatchedInfo
	MatchedInfoStruct MatchedInfo
}
