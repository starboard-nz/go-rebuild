## go-rebuild: a facepalm stopper

Ever happened to you that you made some changes to the source code of your
Go program but the result was exactly as before, only to realise that you
forgot to recompile it? Then this package is for you.

go-rebuild checks if the sources are newer than the executable when the
executable starts saving you having to slap your forehead. At least for this
particular reason...

To use it, simply import the package:

```
import (
	_ "github.com/starboard-nz/go-rebuild/force"
)

```

Importing `force` makes the executable exit (with Exit Status 127) if it needs
to be rebuilt.

If you prefer to be given a warning only, you can import `warning` instead:

```
import (
	_ "github.com/starboard-nz/go-rebuild/warning"
)

```

When building production code, you probably want to disable this, use the `prod`
build tag:

```
go build -tags prod
```

That's it.

## Assumptions

There should be a Go source file in the same directory as the executable. If there isn't one the check will fail with a warning.
The code checks all Go source files (`*.go`) in all subdirectories starting from the root of the source tree, i.e. where `go.mod` is found. If you organised your sources differently then this won't work, but you should probably consider following Go project conventions or you will likely to have to slap your forehead for many other reasons down the line.

The code does not parse `go.mod` (you're welcome to submit a PR), so if you refer to other source trees (e.g. `replace` directives pointing to `../other-project`, those source will not be checked.


## TODO

Automated tests.

Will add some along these lines:

```
laca@localhost:~/go/go-rebuild/test/hello$ go build
laca@localhost:~/go/go-rebuild/test/hello$ ./hello 
Hello, world!
laca@localhost:~/go/go-rebuild/test/hello$ touch ../../check.go
laca@localhost:~/go/go-rebuild/test/hello$ ./hello 
*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*
*~*                             *~*
*~*   ERROR: sources modified   *~*
*~*      Please rebuild me!     *~*
*~*                             *~*
*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*~*
Modified source: /home/laca/go/go-rebuild/check.go

laca@localhost:~/go/go-rebuild/test/hello$ go build -tags prod
laca@localhost:~/go/go-rebuild/test/hello$ ./hello 
Hello, world!
laca@localhost:~/go/go-rebuild/test/hello$ touch ../../check.go
laca@localhost:~/go/go-rebuild/test/hello$ ./hello 
Hello, world!
```
