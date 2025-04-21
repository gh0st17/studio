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
	IType    InterfaceType // Тип интерфейса
	Logging  bool          // Печать логов
	Reg      bool          // Регистрация нового клиента
	HttpSoc  string        // Сокет для веб-сервера
	RedisSoc string        // Сокет для подключения к Redis
	PgSqlSoc string        // Сокет для подключения к pgSQL
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
	flag.StringVar(&interfaceType, "type", "web", interfaceTypeDesc)
	flag.StringVar(&p.HttpSoc, "http", ":8080", httpSocDesc)
	flag.StringVar(&p.RedisSoc, "redis", "localhost:6379", redisSocDesc)
	flag.StringVar(&p.PgSqlSoc, "pgsql", "localhost:5432", pgSqlSocDesc)
	flag.BoolVar(&p.Reg, "registration", false, registrationDesc)

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

	if err = p.checkInterfaceType(interfaceType); err != nil {
		return nil, err
	}
	if p.IType != CLI && p.Reg {
		return nil, ErrRegCLI
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
