package cui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

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
