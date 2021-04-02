# omegleapi
An Omegle API for GoLang, built from the ground up and designed to be concurrency safe.

## Usage
Getting started is as simple as a few lines. For a full example, check out the `echo` example in `examples`.
```go
log := logrus.New()
log.Formatter = &logrus.TextFormatter{ForceColors: true}
log.Level = logrus.DebugLevel
	
err := omegleapi.NewSlave(&CustomHandler{log: log}, []string{"women"}, omegleapi.SlaveOptions{
	Logger: log,
}).Start()
if err != nil {
	panic(err)
}
```
