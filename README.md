<div align="center">
	<a>
		<img src="https://socialify.git.ci/codeiiest/test-runner/image?forks=1&issues=1&language=1&pattern=Circuit%20Board&pulls=1&stargazers=1&theme=Light" />
	</a>
	<br/>
	<b>Test-Runner</b>
	<br/>
	<a href="https://goreportcard.com/report/github.com/codeiiest/test-runner">
		<img src="https://goreportcard.com/badge/github.com/codeiiest/test-runner" alt="Go Report Card" />
	</a>
</div>	

Test-Runner is a Go module which can be used to execute programs against particular
test cases and get corresponding outputs. The goal is to provide an easy module
for testing programs for their correctness. This is one of the many modules
required to create a completely functional backend for an "Online Judge".

- Runs programs in a completely isolated environment using containers
- Set and record memory and time limits using cgroups
- Easy to use interface

## Table of Contents

1. [ Installation ](#install)
2. [ Example ](#example)
3. [ To-do ](#todo)
4. [ Contributing ](#contrib)
5. [ Working ](#working)

<a name="install"></a>

## 1. Installation

`go get github.com/codeiiest/test-runner`

<a name="example"></a>

## 2. Example

```
import (
	"io/ioutil"
	"log"
	"github.com/codeiiest/test-runner/runner/run"
)

func test(code string, lang string, filename string){
	in := []string{"2", "4", "5"}
	out := []string{"4", "16", "25"}
	timeLimit := 2              // Time in seconds
	memLimit := 500*1024*1024   // Memory in bytes

	res := run.Evaluate(code, lang, filename, in, out,
        len(in), timeLimit, memLimit)
	log.Print(res)
}
```

<a name="todo"></a>

## 3. To-do

- Limit Cpu count through cgroups (perhaps choose CPU?)
- Support for more languages (compilation and run scripts)
- Unit Tests
- Makefile

<a name="contrib"></a>

## 4. Contributing

Please use the issue tracker.
All contributions are more than welcome :)

<a name="working"></a>

## 5. Working

The gist of this module is that it runs the code passed to it in a docker container, evaluates it by comparing the output of the sent program against the expected output.

---

- Code is first compiled using their respective compilers [here](./internal/container/run.go) in the function Compile(), or if it is interpreted, directly Run()

- The compilation takes place inside the container via a [bind mounts](https://docs.docker.com/storage/bind-mounts/) and the executable is produced

- That executable is then run with the [evaluate script](./internal/container/docker/evaluate) with parameters of time limit, memory limit, number of test cases, and finally the runner script, which depends on whether the language is compiled or interpreted (for example, directly run python scripts, but C/C++ needs compialtion)

- For a thorough understanding of how to use this module, have a look at the tests created for each language.
