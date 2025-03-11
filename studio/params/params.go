// Пакет params предоставляет набор функции для
// обработки входных параметров программы
//
// Основные функции:
//   - ParseParams: Обрабатывает входные флаги и возвращает
//     структуру [Params] с результатом обработанных флагов
package params

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type InterfaceType uint

const (
	CLI InterfaceType = iota
	Web
)

type Params struct {
	It      InterfaceType // Тип интерфейса
	DBPath  string        // Путь к файлу базы данных
	Logging bool          // Печать логов
}

// Печатает справку
func printHelp() {
	program := filepath.Base(os.Args[0])

	fmt.Println("Использование: ", program, useExample)
	fmt.Printf("\nФлаги:\n")

	flag.PrintDefaults()
}

// Возвращает структуру Params с прочитанными
// входными аргументами программы
func ParseParams() (p *Params, err error) {
	p = &Params{}
	flag.Usage = printHelp
	var interfaceType string
	flag.StringVar(&interfaceType, "type", "", interfaceTypeDesc)
	flag.StringVar(&p.DBPath, "db", "", dbPathDesc)

	var IType string
	flag.StringVar(&IType, "c", "", IType)

	logging := flag.Bool("log", false, logDesc)
	version := flag.Bool("V", false, versionDesc)
	help := flag.Bool("help", false, helpDesc)

	flag.Parse()

	if !*logging {
		log.SetOutput(io.Discard)
	}
	if *version {
		fmt.Print(versionText)
		os.Exit(0)
	}
	if *help {
		printHelp()
		os.Exit(0)
	}

	return p, nil
}

// Проверяет параметр типа компрессора
func (p *Params) checkInterfaceType(it string) error {
	it = strings.ToLower(it)

	switch it {
	case "cli":
		p.It = CLI
	case "lzw":
		p.It = Web
	default:
		return ErrUnknownIType
	}

	return nil
}
