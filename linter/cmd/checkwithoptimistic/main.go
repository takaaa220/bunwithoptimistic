package main

import (
	"github.com/takaaa220/bunwithoptimistic/linter"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(linter.Analyzer) }
