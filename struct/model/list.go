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
 * @brief 空白で区切られたどのqueryが探索文字列内のどのindex番号から
 *        開始する部分文字列なのかを表す
 * Original
 * Filtered
 * text  空白で区切られたqueryの一要素
 * index 探索文字列内でtextが部分文字列として一致するindex番号の先頭
 */
type List struct {
	Original          []ManData
	Filtered          []ManData
	Matched           []MatchedInfo
	MatchedInfoStruct MatchedInfo
}
