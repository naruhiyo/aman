package sio

/*
 * Commands 	入力コマンド
 * Options	 	選択オプション
 * Query      検索クエリ
 * CursorPosX カーソルのx座標
 */
type InputStruct struct {
	Commands   []string
	Options    []string
	Query      string
	CursorPosX int
}
