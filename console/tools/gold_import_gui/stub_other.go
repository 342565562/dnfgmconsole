//go:build !windows

// GUI 工具仅支持 Windows(依赖 lxn/walk)。非 Windows 平台提供空实现，
// 使 `go build ./...` 在 Linux 编译机上不报错。
package main

func main() {}
