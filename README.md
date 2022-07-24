# hack virtual machine

Hack virtual machine built as an exercise of project 7 from [nand2tetris course](https://www.nand2tetris.org/course).

# How to use

Given a single vm file:

    $ go run main.go -f FunctionCalls/SimpleFunction/SimpleFunction.vm

Given a directory with multiple vm files:

    $ go run main.go -d FunctionCalls/FibonacciElement/
