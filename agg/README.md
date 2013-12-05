# Agg

A small library for aggregating statistics about arbitrary actions

# Usage

There's two public functions, `Agg` and `Print`. Read the comments on them for
their usage.

For printing I've been using the following code:

```go
go func() {
    log.Println("Waiting for signal")
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
    <-c
    log.Println("Got SIG")
    agg.Print(1)
    time.Sleep(500 * time.Millisecond)
    os.Exit(0)
}()
```
