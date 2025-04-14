// Пакет userinput предоставляет функции для внутренней
// обработки пользовательского ввода
package userinput

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// Обрабатывает пользовательский ввод
func PromptString(prompt string) string {
	fmt.Printf("%s: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			if scanner.Text() == "" {
				continue
			} else {
				return strings.TrimRight(scanner.Text(), "\r\n")
			}
		}
	}
}

// Обрабатывает пользовательский ввод
func PromptUint(prompt string) (nums []uint, err error) {
	var u64 uint64

	numsStr := strings.Split(PromptString(prompt), " ")
	for _, numStr := range numsStr {
		if numStr == "" {
			continue
		}

		if u64, err = strconv.ParseUint(numStr, 10, 0); err == nil {
			nums = append(nums, uint(u64))
		}
	}

	if len(nums) == 0 {
		return PromptUint(prompt)
	}

	return nums, nil
}

// Проверяет является ли стандратный ввод терминалом
func IsNonInteractive() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode()&os.ModeCharDevice) == 0 || !term.IsTerminal(int(os.Stdin.Fd()))
}
