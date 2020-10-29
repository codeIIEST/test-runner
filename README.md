Test-Runner
===========

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

<a name="install"></a>
## 1. Installation

``` go get github.com/codeiiest/test-runner ```

<a name="example"></a>
## 2. Example

```
import (
	"io/ioutil"
	"log"
	"github.com/codeiiest/test-runner/runner/run"
)

func test(code string, lang string, filename string)
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
