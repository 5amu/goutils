# goutils

Repository containing project-independent utilities to be imported in go projects.

## filequeue

[filequeue](filequeue/filequeue.go) is a file reader that sends the line into a channel.

```go
import (
    "fmt"
    "log"

    "github.com/5amu/goutils/filequeue"
)

func main() {

    filename := "test.txt"

    fq := filequeue.NewFileQueue()
    go func() {
        if err := fq.ScanFile(filename); err != nil {
            log.Fatal(err)
        }
    }()

    for {
        select{
        case line := <-fq.Pop():
            fmt.Print(line)
        case <-fq.IsEmpty():
            log.Fatal("queue is empty")
        }
    }
}
```