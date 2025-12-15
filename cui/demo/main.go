package main

import (
	"fmt"
	"github.com/sw965/omw/cui"
	"time"
)

func main() {
	// 画面がまっさらになる → Goと表示される → 画面がまっさらになる → Pythonと表示される → 画面がまっさらになる → C++ と表示される
	for _, s := range []string{"Go", "Python", "C++"} {
		if err := cui.ClearConsole(); err != nil {
			fmt.Println("ClearConsole failed:", err)
			return
		}
		fmt.Println(s)
		time.Sleep(1 * time.Second)
	}

	err := cui.ClearConsole()
	if err != nil {
		fmt.Println("ClearConsole failed:", err)
	}
}
