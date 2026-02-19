// Package main — утилита для генерации метода Reset() для структур с меткой // generate:reset.
//
// Использование:
//
//	go run ./cmd/reset
//
// Алгоритм:
// 1. Сканирует все пакеты начиная с корня
// 2. Для каждой структуры с комментарием // generate:reset
// 3. Генерирует метод Reset()
// 4. Сохраняет в reset.gen.go в том же пакете
package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	g := &Generator{
		fset:     token.NewFileSet(),
		packages: make(map[string]*Package),
	}

	// ❌ Убираем filepath.Walk + parser.ParseDir
	// ✅ Загружаем через packages
	if err := g.loadPackages(); err != nil {
		return fmt.Errorf("ошибка загрузки пакетов: %w", err)
	}

	// Генерируем файлы
	for _, pkg := range g.packages {
		if len(pkg.Structs) == 0 {
			continue
		}
		if err := g.generateResetFile(pkg); err != nil {
			return err
		}
		fmt.Printf("Сгенерировано: %s/reset.gen.go (%d структур)\n", pkg.Dir, len(pkg.Structs))
	}

	fmt.Println("Генерация завершена")
	return nil
}

func getPackageName(importPath string) string {
	// Убираем версии: github.com/pkg/v3 → github.com/pkg
	importPath = removeVersionSuffix(importPath)

	// Берём последнюю часть пути
	return filepath.Base(importPath)
}

// removeVersionSuffix убирает суффиксы вроде /v2, /v3 из пути
func removeVersionSuffix(path string) string {
	for i := len(path) - 1; i > 0 && path[i] >= '0' && path[i] <= '9'; i-- {
		if path[i-1] == '/' && path[i-2] == 'v' {
			return path[:i-2]
		}
	}
	return path
}
func isStdlibType(typeName, pkgPath string) bool {
	if pkgPath == "" {
		return false
	}
	pkgName := getPackageName(pkgPath)

	stdPkgs := map[string]bool{
		"sync": true, "os": true, "database/sql": true,
		"time": true, "context": true, "net": true,
	}

	if stdPkgs[pkgName] {
		return true
	}

	// Точные имена типов
	switch typeName {
	case "sync.Mutex", "sync.RWMutex",
		"*os.File", "*os.Stdout", "*os.Stderr", "*os.Stdin",
		"*sql.DB", "*sql.Tx", "sql.NullString", "sql.driver",
		"context.Context", "context.cancelCtx", "context.timerCtx",
		"time.Time", "time.Ticker", "time.Timer":
		return true
	}

	return false
}

// hasResetMethod — заглушка. В реальном анализе используй types.Info.
// Здесь мы просто предполагаем, что если имя содержит "Struct", то Reset() может быть.
// Для точности — нужно анализировать types.Named.Methods.
func hasResetMethod(typeName string) bool {
	return false // временно отключено — безопаснее
}

func formatExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return formatExpr(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + formatExpr(e.X)
	case *ast.ArrayType:
		if e.Len == nil {
			return "[]" + formatExpr(e.Elt)
		}
		return fmt.Sprintf("[%s]%s", e.Len, formatExpr(e.Elt))
	case *ast.MapType:
		return "map[" + formatExpr(e.Key) + "]" + formatExpr(e.Value)
	default:
		return "interface{}" // fallback
	}
}

func isPointer(expr ast.Expr) bool {
	_, ok := expr.(*ast.StarExpr)
	return ok
}

func isNamedStruct(expr ast.Expr, info *types.Info) bool {
	typ := info.TypeOf(expr)
	if typ == nil {
		return false
	}
	// Убираем указатели
	for {
		ptr, ok := typ.(*types.Pointer)
		if !ok {
			break
		}
		typ = ptr.Elem()
	}
	_, isStruct := typ.Underlying().(*types.Struct)
	return isStruct
}

func isSlice(expr string) bool { return strings.HasPrefix(expr, "[]") }
func isMap(expr string) bool   { return strings.HasPrefix(expr, "map[") }

func zeroValue(t string) string {
	switch {
	case isSlice(t), isMap(t):
		return t + "{}"
	case t == "string":
		return `""`
	case t == "bool":
		return "false"
	case strings.HasPrefix(t, "int") || strings.HasPrefix(t, "uint") || strings.HasPrefix(t, "float"):
		return "0"
	default:
		return t + "{}" // для структур
	}
}
