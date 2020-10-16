package modules

import (
	"bytes"
	"errors"
	"flag"
	"os/exec"
	"regexp"
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

	// 検索フラグ
	// true: 検索中、false: 検索してない
	var isFinding bool = true

	// buffer の方が string結合より効率が良い
	var buffer bytes.Buffer // オプション説明のブロックを入れる変数
	var results []string
	// オプション条件
	//   - `-` または `--` から始まり、半角英字が続く文字列であること
	var reg *regexp.Regexp = regexp.MustCompile(`-{1,2}[a-zA-Z]`)

	// オプション条件を満たしているかをチェック
	var isOptionText = func(line string) bool {
		return reg.MatchString(line)
	}

	// オプションの前に何文字空白があるかカウントする
	var getOptionHeaderBlankCounts = func(line string) int {
		var count int = 0
		if strings.Contains(line, "--") {
			count = strings.Index(line, "--") | 0
		} else if strings.Contains(line, "-") {
			count = strings.Index(line, "-") | 0
		}
		return count
	}

	// オプションのヘッダーであるかチェック（オプションの説明文にあるハイフンを弾く）
	var isOptionHeaderText = func(line string, count int) bool {
		var blanks string = strings.Repeat(" ", count)

		// 先頭から文字を見たときにオプション条件を満たしているか確認する
		if strings.HasPrefix(line, blanks+"--") || strings.HasPrefix(line, blanks+"-") {
			return true
		}

		return false
	}

	// オプションを配列に追加する
	var count int = 0
	for _, line := range splitOutputs {
		if isFinding {
			// オプションのヘッダーに来るまでスキップ
			if !isOptionText(line) {
				continue
			}

			// オプションの前に何文字空白があるか計算
			count = getOptionHeaderBlankCounts(line)

			if !isOptionHeaderText(line, count) {
				continue
			}

			isFinding = false
			buffer.WriteString(line[count:])
		} else {
			// 改行だった場合次のオプションを探す
			if len(line) == 0 {
				results = append(results, buffer.String())
				buffer.Reset()
				isFinding = true
				continue
			}
			// バッファに文字列追加
			buffer.WriteString("\n" + line[count:])
		}
	}

	return results
}
