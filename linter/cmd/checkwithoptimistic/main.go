package main

import (
	"github.com/takaaa220/bunwithoptimistic/linter"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(linter.Analyzer) }
