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
	"studio/filesystem"
)

type InterfaceType uint

const (
	CLI InterfaceType = iota
	Web
)

type Params struct {
	IType   InterfaceType // Тип интерфейса
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

	logging := flag.Bool("log", false, logDesc)
	version := flag.Bool("V", false, versionDesc)
	help := flag.Bool("help", false, helpDesc)

	flag.Parse()

	if *version {
		fmt.Print(versionText)
		os.Exit(0)
	}
	if *help {
		printHelp()
		os.Exit(0)
	}

	if interfaceType == "" {
		return nil, ErrMissingIType
	}
	if p.DBPath == "" {
		return nil, ErrMissingDBPath
	}
	if !filesystem.Exsists(p.DBPath) {
		return nil, ErrDBNotExists
	}
	if err = p.checkInterfaceType(interfaceType); err != nil {
		return nil, err
	}

	if !*logging {
		log.SetOutput(io.Discard)
	}

	return p, nil
}

// Проверяет параметр для типа интерфейса
func (p *Params) checkInterfaceType(it string) error {
	it = strings.ToLower(it)

	switch it {
	case "cli":
		p.IType = CLI
	case "web":
		p.IType = Web
	default:
		return ErrUnknownIType
	}

	return nil
}
