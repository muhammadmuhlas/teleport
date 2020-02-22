package internal

import (
	"bufio"
	"fmt"
	"os"
)

func InArray(s string, ss []string) bool {
	for i := 0; i < len(ss); i++ {
		if ss[i] == s {
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

func BitbucketHTTPURLBuilder(username, namespace, path string) string {
	return "https://" + username + "@bitbucket.org/" + namespace + "/" + path + ".git"
}
