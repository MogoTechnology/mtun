# README

## 单元测试

单元测试需要设置：`go env -w CGO_ENABLED=1`, 不然报错：
```
ping\ping_test.go:76:16: undefined: getScor
```

因为 ping.go 有 `import "C"`，CGO_ENABLED=0 时不会编译。
