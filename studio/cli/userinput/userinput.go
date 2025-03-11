// Пакет userinput предоставляет функции для внутренней
// обработки пользовательского ввода
package userinput

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Обрабатывает пользовательский ввод
func Prompt(prompt string) (string, error) {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimRight(scanner.Text(), "\r\n"), nil
	}
	return "", ErrInput
}

// Проверяет является ли стандратный ввод терминалом
func IsNonInteractive() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode()&os.ModeCharDevice) == 0 || !term.IsTerminal(int(os.Stdin.Fd()))
}
