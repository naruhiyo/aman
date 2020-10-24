package modules

import (
	"bytes"
	"errors"
	"flag"
	"os/exec"
	"strings"
)

/**
* コマンドライン引数を取得
 */
func Parse() []string {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		panic(errors.New("Error: No arguments"))
	}

	return args
}

/**
* man コマンドを実行する
* @params args 実行時引数
**/
func ExecMan(args []string) string {
	// man コマンドは空白区切のコマンドをハイフンで管理しているため、ハイフンつなぎに変更
	const MAN string = "man"
	var command string = strings.Join(args, "-")

	// manコマンドを実行する
	var out, err = exec.Command(MAN, command).Output()

	if err != nil {
		panic(errors.New("Error: No results"))
	}

	return string(out)
}

/**
* コマンド実行結果からオプションを抽出する
* @params output manコマンド実行結果
* @return オプションテキストリスト
**/
func AnalyzeOutput(output string) []string {
	// === 条件 ===
	// ハイフンまたはダブルハイフンで始まる英単語
	var splitOutputs []string = strings.Split(output, "\n")

	// オプションを見つけたかどうか
	var isFinding bool = true

	// buffer の方が string結合より効率が良い
	var buffer bytes.Buffer // オプション説明のブロックを入れる変数
	var results []string

	for _, line := range splitOutputs {
		// オプション判定条件
		// 1. 半角3文字以上の空白後に`--` または `-` で始まるオプションであること
		// 2. 前のオプションの説明ではないこと
		if isFinding {
			// オプションが来るまでスキップ
			if !(strings.Contains(line, "   --") || strings.Contains(line, "   -")) {
				continue
			}

			isFinding = false
			buffer.WriteString(line)
		} else {
			// 改行だった場合次のオプションを探す
			if len(line) == 0 {
				results = append(results, buffer.String())
				buffer.Reset()
				isFinding = true
				continue
			}
			// バッファに文字列追加
			buffer.WriteString("\n" + line)
		}
	}

	return results
}
