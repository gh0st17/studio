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
func PromptUint(prompt string) (_ uint, err error) {
	var s string
	if s, err = PromptString(prompt); err != nil {
		return 0, err
	}

	var i64 uint64
	if i64, err = strconv.ParseUint(s, 10, 0); err != nil {
		return 0, err
	}

	return uint(i64), nil
}

// Проверяет является ли стандратный ввод терминалом
func IsNonInteractive() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode()&os.ModeCharDevice) == 0 || !term.IsTerminal(int(os.Stdin.Fd()))
}
