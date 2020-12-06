package spagination

/*
 * Page       	現在のページ番号 定義域は[0, maxPage]
 * MaxPage    	最大ページ番号
 * SelectedPos 	選択行
 * PageList			ページ毎に表示する件数
 */
type Pagination struct {
	Page        int
	MaxPage     int
	SelectedPos int
	PageList    []int
}
