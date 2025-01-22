Simple demo that can be cross compiled from Linux to Windows and executed via Wine with a printer named `PDF`

```go
GOOS=windows GOARCH=amd64 go build cmd
```
