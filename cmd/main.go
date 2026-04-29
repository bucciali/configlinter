package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"configlinter/internal/config"
	"configlinter/internal/engine"
	"configlinter/internal/parser"
	"configlinter/internal/reporter"
	"configlinter/internal/rules"
	"configlinter/internal/scanner"
	"configlinter/internal/server"
)

const version = "1.0.0"

func buildParserRegistry() *parser.Registry {
	reg := parser.NewRegistry()
	reg.Register(&parser.YAMLParser{})
	reg.Register(&parser.JSONParser{})
	reg.Register(&parser.TOMLParser{})
	return reg
}

func buildEngine() *engine.Engine {
	return engine.New(
		rules.NewDebugLogRule(),
		rules.NewPlaintextPasswordRule(),
		rules.NewBindAllRule(),
		rules.NewTLSDisabledRule(),
		rules.NewWeakCryptoRule(),
		&rules.FilePermissionsRule{},
	)
}

func main() {
	if len(os.Args[1:]) == 0 {
		interactiveMode()
		return
	}
	silent := false
	stdin := false
	serve := false
	var dirPath string
	var filePath string

	for _, arg := range os.Args[1:] {
		switch arg {
		case "-s", "--silent":
			silent = true
		case "--stdin":
			stdin = true
		case "--serve":
			serve = true
		default:
			if strings.HasPrefix(arg, "--dir=") {
				dirPath = strings.TrimPrefix(arg, "--dir=")
			} else if arg == "--dir" {

				dirPath = "__next__"
			} else if dirPath == "__next__" {
				dirPath = arg
			} else if filePath == "" && !strings.HasPrefix(arg, "-") {
				filePath = arg
			}
		}
	}

	if serve {
		cfg := config.LoadServerConfig()
		reg := buildParserRegistry()
		eng := buildEngine()

		if err := server.Start(cfg, reg, eng); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка сервера: %v\n", err)
			os.Exit(1)
		}
		return
	}

	reg := buildParserRegistry()
	eng := buildEngine()
	if dirPath != "" {
		results, err := scanner.ScanDir(dirPath, reg, eng)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка сканирования: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(scanner.FormatResults(results))

		if !silent && scanner.HasFindings(results) {
			os.Exit(2)
		}
		return
	}

	var data []byte
	var p parser.Parser
	var err error

	if stdin {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка чтения stdin: %v\n", err)
			os.Exit(1)
		}
		format := detectFormat(data)
		p, err = reg.GetByFormat(format)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
	} else {
		if filePath == "" {
			fmt.Fprintf(os.Stderr, "Использование: configlinter <файл> [--flags]\n")
			fmt.Fprintf(os.Stderr, "  или: cat config.yaml | configlinter --stdin\n")
			fmt.Fprintf(os.Stderr, "  или: configlinter --dir <путь>\n")
			fmt.Fprintf(os.Stderr, "  или: configlinter --serve\n")
			os.Exit(1)
		}
		data, err = os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка чтения файла: %v\n", err)
			os.Exit(1)
		}
		p, err = reg.GetByFilename(filepath.Base(filePath))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
	}

	root, err := p.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсинга: %v\n", err)
		os.Exit(1)
	}

	if filePath != "" {
		info, statErr := os.Stat(filePath)
		if statErr == nil {
			root.FilePath = filePath
			root.FileMode = info.Mode().Perm()
		}
	}

	findings := eng.Analyze(root)

	rep := reporter.NewTextReporter()
	rep.Report(os.Stdout, findings)

	if len(findings) > 0 && !silent {
		os.Exit(2)
	}
}

func detectFormat(data []byte) string {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		return "json"
	}
	for _, line := range strings.Split(trimmed, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 2 && line[0] == '[' && strings.HasSuffix(line, "]") && !strings.Contains(line, "=") {
			return "toml"
		}
	}
	return "yaml"
}

func interactiveMode() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("  ConfigLinter v" + version)
	fmt.Println("  Проверка конфигов (YAML, JSON, TOML) на проблемы безопасности.")
	fmt.Println()
	fmt.Println("  1 — Проверить файл")
	fmt.Println("  2 — Проверить папку")
	fmt.Println("  3 — Запустить сервер (REST API)")
	fmt.Println("  4 — Справка")
	fmt.Println("  0 — Выход")
	fmt.Println()

	for {
		fmt.Print("  > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Print("  Путь к файлу: ")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)
			path = strings.Trim(path, "\"'")
			if path == "" {
				fmt.Println("  Путь не указан.")
				fmt.Println()
				continue
			}
			fmt.Println()
			runFileInteractive(path)
			waitExit(reader)
			return

		case "2":
			fmt.Print("  Путь к папке: ")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)
			path = strings.Trim(path, "\"'")
			if path == "" {
				fmt.Println("  Путь не указан.")
				fmt.Println()
				continue
			}
			fmt.Println()
			runDirInteractive(path)
			waitExit(reader)
			return

		case "3":
			fmt.Println()
			fmt.Println("  Запуск сервера...")
			cfg := config.LoadServerConfig()
			reg := buildParserRegistry()
			eng := buildEngine()
			if err := server.Start(cfg, reg, eng); err != nil {
				fmt.Fprintf(os.Stderr, "  Ошибка сервера: %v\n", err)
			}
			return

		case "4":
			fmt.Println()
			fmt.Println("  Использование:")
			fmt.Println("    configlinter <файл>           проверить один файл")
			fmt.Println("    configlinter --dir=<папка>    проверить все конфиги в папке")
			fmt.Println("    configlinter --stdin          читать из stdin")
			fmt.Println("    configlinter --serve          запустить сервер")
			fmt.Println()
			fmt.Println("  Флаги:")
			fmt.Println("    -s, --silent   не возвращать код ошибки")
			fmt.Println()

		case "0", "q", "exit":
			return

		default:
			fmt.Println("  Введите 1, 2, 3, 4 или 0.")
			fmt.Println()
		}
	}
}

func runFileInteractive(path string) {
	reg := buildParserRegistry()
	eng := buildEngine()

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Ошибка: %v\n", err)
		return
	}

	p, err := reg.GetByFilename(filepath.Base(path))
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Ошибка: %v\n", err)
		return
	}

	root, err := p.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Ошибка парсинга: %v\n", err)
		return
	}

	info, statErr := os.Stat(path)
	if statErr == nil {
		root.FilePath = path
		root.FileMode = info.Mode().Perm()
	}

	findings := eng.Analyze(root)
	rep := reporter.NewTextReporter()
	rep.Report(os.Stdout, findings)
}

func runDirInteractive(dir string) {
	reg := buildParserRegistry()
	eng := buildEngine()

	results, err := scanner.ScanDir(dir, reg, eng)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Ошибка: %v\n", err)
		return
	}
	fmt.Print(scanner.FormatResults(results))
}

func waitExit(reader *bufio.Reader) {
	fmt.Println()
	fmt.Println("  Нажмите Enter чтобы закрыть...")
	reader.ReadString('\n')
}
