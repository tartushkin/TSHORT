package main

import "go/token"

// Generator управляет процессом генерации.
type Generator struct {
	fset     *token.FileSet
	packages map[string]*Package // path → Package
}

// Package содержит информацию о Go-пакете.
type Package struct {
	Name    string
	Dir     string
	Structs []*Struct
}

// Struct — структура, помеченная для генерации Reset().
type Struct struct {
	Name   string
	Fields []*Field
}

// Field — поле структуры.
type Field struct {
	Name      string
	Type      string
	IsPointer bool
	IsStruct  bool
	PkgPath   string // полный import (если внешний)
}
