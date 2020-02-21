package internal

import (
	"bufio"
	"fmt"
	"os"
)

func InArray(s string, ss []interface{}) bool {
	for i := 0; i < len(ss); i++ {
		if ss[i].(string) == s {
			return true
		}
	}
	return false
}

func Scanner(prompt string, f func(input string) bool) {
	for {
		fmt.Print(prompt)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if !f(scanner.Text()) {
			break
		}
	}
}
