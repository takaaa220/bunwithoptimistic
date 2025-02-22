package main

import (
	"github.com/takaaa220/bunwithoptimistic"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(bunwithoptimistic.Analyzer) }
