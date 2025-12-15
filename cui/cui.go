// Package cui provides utilities for Console User Interface interactions.
//
// Package cui は、コンソールユーザーインターフェースの操作に関するユーティリティを提供します。
package cui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ClearConsole clears the current terminal screen.
//
// On Windows, it executes the "cls" command to ensure compatibility with older command prompts.
// On other operating systems (Linux, macOS, etc.), it uses ANSI escape sequences
// ("\033[H\033[2J") for faster performance without spawning a new process.
// It returns an error if the Windows command fails to execute.
//
// ClearConsole は現在のターミナル画面を消去します。
//
// Windowsでは、古いコマンドプロンプトとの互換性を保つために "cls" コマンドを実行します。
// その他のOS（Linux, macOSなど）では、新しいプロセスを生成せず高速に動作させるため、
// ANSIエスケープシーケンス ("\033[H\033[2J") を使用します。
// Windowsコマンドの実行に失敗した場合、エラーを返します。
func ClearConsole() error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}

	// Unix系システムではANSIエスケープシーケンスを使用
	// \033[H  -> カーソルをホーム位置（左上）へ移動
	// \033[2J -> 画面全体を消去
	fmt.Print("\033[H\033[2J")
	return nil
}
