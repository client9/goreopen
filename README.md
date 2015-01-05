

```go
writer := goreopen.File("afile")
logger.SetOutput(writer)

c := make(chan os.Signal, 1)
signal.Notify(c, syscall.SIGHUP)
go func() {
    for {
        <-c
        l.Reopen()
    }
}()
```