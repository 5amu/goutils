# shuttles

[shuttles](shuttles/) is a process launcher that is easy enough to parallelize thanks to Go easy concurrency. See my [implementation](cmd/shuttles/main.go)

```go
factory := shuttles.NewShuttleFactory(shuttleFlagSet.Args(), lanes)
for index, fname := range fuel {
    go func(n string, i int) {
        if err := loadFuel(factory, n, i); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    }(fname, index)
}

if err := factory.Start(context.Background(), lanes); err != nil {
    fmt.Println(err)
    os.Exit(1)
}

outputs := factory.GetShuttleOutputs()
```