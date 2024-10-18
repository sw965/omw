package cui

import (
	"os"
	"os/exec"
	"runtime"
)

func ClearConsole() {
	var clearCmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		clearCmd = exec.Command("cmd", "/c", "cls")
	default:
		clearCmd = exec.Command("clear")
	}
	clearCmd.Stdout = os.Stdout
	clearCmd.Run()
}
