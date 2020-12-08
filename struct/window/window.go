package swindow

import "github.com/nsf/termbox-go"

/*
 * Height     ウィンドウの高さ
 * Width      ウィンドウの幅
 * TextColor  テキストカラー
 * BgColor    テキストの背景カラー
 */
type WindowInfo struct {
	Height    int
	Width     int
	TextColor termbox.Attribute
	BgColor   termbox.Attribute
}
