# pkg

## 简介

`pkg` 是一个 Go 语言的工具包，提供了一些常用的工具函数、数据结构和常用组件封装。

## 功能特性

* [safemap](https://github.com/go-inspire/pkg/tree/main/safemap): 线程安全、高性能的 map 结构
* [ringbuffer](https://github.com/go-inspire/pkg/tree/main/ringbuffer): 环形数据结构
* [casbin](https://github.com/go-inspire/pkg/tree/main/casbin): Casbin 的 upper 适配器
* [encoding](https://github.com/go-inspire/pkg/tree/main/encoding): 编解码 json，针对不同的平台使用不同的库，默认使用 [gojson](https://github.com/goccy/go-json)
* [log](https://github.com/go-inspire/pkg/tree/main/log): 日志封装, 允许通过编辑配置文件动态修改日志级别

## 如何使用

```shell
go get -u github.com/go-inspire/pkg@main

```

## 许可证

`pkg` 采用 MIT 许可证。 详见 [LICENSE](./LICENSE) 文件。


