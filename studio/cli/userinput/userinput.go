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
func PromptString(prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimRight(scanner.Text(), "\r\n"), nil
	}
	return "", ErrInput
}

// Обрабатывает пользовательский ввод
func PromptUint(prompt string) (nums []uint, err error) {
	var s string
	if s, err = PromptString(prompt); err != nil {
		return nil, err
	}

	var u64 uint64
	numsStr := strings.Split(s, " ")
	for _, numStr := range numsStr {
		if numStr == "" {
			continue
		}

		if u64, err = strconv.ParseUint(numStr, 10, 0); err != nil {
			return nil, err
		}
		nums = append(nums, uint(u64))
	}

	return nums, nil
}

// Проверяет является ли стандратный ввод терминалом
func IsNonInteractive() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode()&os.ModeCharDevice) == 0 || !term.IsTerminal(int(os.Stdin.Fd()))
}
