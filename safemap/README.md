# safemap 包

`safemap` 包提供了线程安全的集合数据结构。该包包含以下主要类型和函数：

## 类型

- `HashSet[T]` 是一个集合数据结构，使用 Go 的内置 `map` 实现。
- `SafeHashSet[T]` 是一个线程安全的 `HashSet`，使用读写锁实现。
- `SafeMap[T]` 是一个线程安全的 `map[string]T`，使用读写锁实现。
- `SharedSafeMap[T]` 是一个线程安全的 `map[string]T`，使用分片思路优化多些性能。
- `SharedChannel[T]` 是一个线程安全的 `chan T`，使用消息分片思路优化多些性能。


