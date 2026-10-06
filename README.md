# disjoint-set

集合合并与分组报告：把合并请求吃进去，再按组生成文本报告。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/dset --sample=order
go run ./cmd/dset --sample=members
go run ./cmd/dset --sample=same
go run ./cmd/dset --sample=replay
go run ./cmd/dset --sample=work
```

## 口径

- **同组判定**：两个成员是否同组，看的是它们所在集合的根，不是它们的直接父亲。
- **合并策略**：把成员少的集合挂到成员多的集合下（小挂大），同规模时按成员 id 升序里靠前的当根。
- **路径**：查找时把沿途节点压到根上（路径压缩）。
- **报告顺序**：组之间按组内最小成员 id 升序，组内成员按 id 升序。
- **幂等**：同一个 `req_id` 只生效一次，重复请求不改变结构也不增加合并次数。

## 不变量

- 合并次数不多于生效的不同请求数。
- 报告里每个成员只出现一次。
- `scanned` 不随成员数平方级放大：两千次串联合并的 `scanned` 不超过 9000。

## 输出契约

报告每个组一行：`group N: 成员1,成员2`
末尾一行 `groups=N merges=N`。
`--sample=work` 只打印 `{"merges": N, "scanned": N}`。
