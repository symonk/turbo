# Turbo 🧬

⚡ A blazing-fast worker pool library for Go, built with concurrency and simplicity in mind.

## Features

- 🌀 Minimal API
- ⚙️ Efficient goroutine scheduling
- 🧵 Safe concurrent task execution
- 📦 Designed as a pure library – no executables

## Example

```go
p := turbo.NewPool(4)
defer p.Close()

p.Submit(func() {
    fmt.Println("Hello from a worker!")
})

p.Wait()