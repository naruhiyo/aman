package modules

import (
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
* コマンドを受け取り、オプションを取得する
* @params args 実行時引数
**/
func GetOptions(args []string) string {
	var options = [3]string{"-help", "--help", "-H"}
	var command string = strings.Join(args, " ")

	// それぞれのオプションでコマンドを実行する
	var results []string
	var cmdErrors []error
	for _, option := range options {
		var out, err = exec.Command(command, string(option)).Output()

		if err != nil {
			cmdErrors = append(cmdErrors, err)
		} else {
			results = append(results, string(out))
		}
	}

	if len(cmdErrors) == len(options) {
		panic(errors.New("Error: No results"))
	}

	return strings.Join(results, ",")
}

/**
* コマンド実行結果からオプションを抽出する
 */
func AnalyzeOutput(output string) string {
	// === 条件 ===
	// ハイフンまたはダブルハイフンで始まる英単語
	reg := regexp.MustCompile(`-?-[a-zA-Z0-9\-]+`)

	var optionList = reg.FindAllString(output, -1)

	return strings.Join(optionList, "\n")
}

// func Sample() {
// 	err := termbox.Init()
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer termbox.Close()

// 	for {
// 		switch ev := termbox.PollEvent(); ev.Type {
// 		case termbox.EventKey:
// 			switch ev.Key {
// 			case termbox.KeyEsc:
// 				return
// 			default:
// 				fmt.Println(ev.Key)
// 				continue
// 			}
// 		default:
// 			return
// 		}
// 	}
// }
