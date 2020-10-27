package modules

import (
	"bytes"
	"errors"
	"flag"
	"os/exec"
	"strings"
)

type ManData struct {
	Contents string
	LineNumber int
}

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
* オプションの判定方法
*   - `-` (または `--`) を検索する
*   - 検索結果に対して`-`(または `--`) の出現位置(index)を計算する
*   - 出現位置(index)の値だけ空白文字を生成し、オプションと結合する
*   - 結合した値と元の値を比較し、一致すればオプションとみなす
* @params output manコマンド実行結果
* @return オプションテキストリスト
**/
func AnalyzeOutput(output string) []ManData {
	// === 条件 ===
	// ハイフンまたはダブルハイフンで始まる英単語
	var splitOutputs []string = strings.Split(output, "\n")

	// 検索フラグ
	// true: 検索中、false: 検索してない
	var isFinding bool = true

	// buffer の方が string結合より効率が良い
	var buffer bytes.Buffer // オプション説明のブロックを入れる変数
	var results []ManData

	// オプション条件を満たしているかをチェック
	var isOptionText = func(line string) bool {
		return strings.Contains(line, "-")
	}

	// isOptionHeaderText()内で生成する空白文字の文字数を求める
	var getOptionHeaderBlankCounts = func(line string) int {
		var count int = 0
		if strings.Contains(line, "-") {
			count = strings.Index(line, "-")
		}
		return count
	}

	// オプションのヘッダーであるかチェック（オプションの説明文にあるハイフンを弾く）
	var isOptionHeaderText = func(line string, count int) bool {
		var blanks string = strings.Repeat(" ", count)

		// 先頭から文字を見たときにオプション条件を満たしているか確認する
		return strings.HasPrefix(line, blanks+"-")
	}

	// オプションを配列に追加する
	var count int = 0
	// オプションに必要な空白の個数
	var definedOptionBlankCount int = -1
	for _, line := range splitOutputs {
		if isFinding {
			// オプションのヘッダーに来るまでスキップ
			if !isOptionText(line) {
				continue
			}

			// オプションの前に何文字空白があるか計算
			count = getOptionHeaderBlankCounts(line)

			// オプションの条件を満たしているか確認
			if !isOptionHeaderText(line, count) {
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

			// fmt.Println(line)
			isFinding = false
			buffer.WriteString(line[definedOptionBlankCount:])
		} else {
			// 改行だった場合次のオプションを探す
			if len(line) == 0 {
				results = append(results, ManData {
					Contents: buffer.String(),
					LineNumber: strings.Count(buffer.String(), "\n") + 2,
				})
				buffer.Reset()
				isFinding = true
				continue
			}

			// バッファに文字列追加
			buffer.WriteString("\n" + line[definedOptionBlankCount:])
		}
	}

	return results
}
