// Package dset 维护集合合并并生成分组报告。
package dset

import (
	"fmt"
	"strings"
)

// Counter 统计本轮沿父链走过的步数。
type Counter struct {
	Scanned int
}

// Set 是一组不交集合。
type Set struct {
	parent map[string]string
	size   map[string]int
	order  []string
	merges int
	seen   []string
}

// New 建一个空集合。
func New() *Set {
	return &Set{parent: map[string]string{}, size: map[string]int{}}
}

func (s *Set) add(x string) {
	if _, ok := s.parent[x]; ok {
		return
	}
	s.parent[x] = x
	s.size[x] = 1
	s.order = append(s.order, x)
}

func (s *Set) find(x string, c *Counter) string {
	s.add(x)
	for s.parent[x] != x {
		c.Scanned++
		x = s.parent[x]
	}
	return x
}

// Union 合并两个成员所在的集合。
func (s *Set) Union(a, b, reqID string, c *Counter) {
	ra := s.find(a, c)
	rb := s.find(b, c)
	if ra == rb {
		return
	}
	s.parent[rb] = ra
	s.merges++
	s.seen = append(s.seen, reqID)
}

// Same 报告两个成员是否同组。
func (s *Set) Same(a, b string, c *Counter) bool {
	return s.parent[a] == s.parent[b]
}

// Merges 返回已经生效的合并次数。
func (s *Set) Merges() int { return s.merges }

// Each 逐个吐出分组，组内成员按加入顺序。
func (s *Set) Each(fn func([]string), c *Counter) {
	byRoot := map[string][]string{}
	var roots []string
	for _, x := range s.order {
		r := s.find(x, c)
		if _, ok := byRoot[r]; !ok {
			roots = append(roots, r)
		}
		byRoot[r] = append(byRoot[r], x)
	}
	for _, r := range roots {
		fn(byRoot[r])
	}
}

// Render 生成文本报告。
func (s *Set) Render(c *Counter) string {
	var b strings.Builder
	n := 0
	s.Each(func(group []string) {
		n++
		b.WriteString(fmt.Sprintf("group %d: %s\n", n, strings.Join(group, ",")))
	}, c)
	b.WriteString(fmt.Sprintf("groups=%d merges=%d\n", n, s.merges))
	return b.String()
}
