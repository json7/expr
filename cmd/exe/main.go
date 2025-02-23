package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
	"github.com/sanity-io/litter"
	"io/ioutil"
	"os"
)

var (
	bytecode bool
	debug    bool
	run      bool
	ast      bool
	dot      bool
	repl     bool
	opt      bool
)

func init() {
	flag.BoolVar(&bytecode, "bytecode", false, "disassemble bytecode")
	flag.BoolVar(&debug, "debug", false, "debug program")
	flag.BoolVar(&run, "run", false, "run program")
	flag.BoolVar(&ast, "ast", false, "print ast")
	flag.BoolVar(&dot, "dot", false, "dot format")
	flag.BoolVar(&repl, "repl", false, "start repl")
	flag.BoolVar(&opt, "opt", true, "do optimization")
}

func main() {
	flag.Parse()

	if ast {
		printAst()
		os.Exit(0)
	}
	if bytecode {
		printDisassemble()
		os.Exit(0)
	}
	if run {
		runProgram()
		os.Exit(0)
	}
	if debug {
		debugger()
		os.Exit(0)
	}
	if repl {
		startRepl()
		os.Exit(0)
	}

	flag.Usage()
	os.Exit(2)
}

func input() string {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func check(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func printAst() {
	tree, err := parser.Parse(input())
	check(err)
	if !dot {
		litter.Dump(tree.Node)
		return
	}

	if opt {
		optimizer.Optimize(&tree.Node)
	}

	dotAst(tree.Node)
}

func printDisassemble() {
	tree, err := parser.Parse(input())
	check(err)

	_, err = checker.Check(tree, nil)
	check(err)

	if opt {
		optimizer.Optimize(&tree.Node)
	}

	program, err := compiler.Compile(tree, nil)
	check(err)

	_, _ = fmt.Fprintf(os.Stdout, program.Disassemble())
}

func runProgram() {
	out, err := expr.Eval(input(), nil)
	check(err)

	litter.Dump(out)
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		line := scanner.Text()
		out, err := expr.Eval(line, nil)
		if err != nil {
			fmt.Printf("%v\n", err)
			goto prompt
		}

		fmt.Printf("%v\n", litter.Sdump(out))

	prompt:
		fmt.Print("> ")
	}
}
