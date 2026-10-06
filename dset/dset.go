// Package dset 维护集合合并并生成分组报告。
package dset

import (
	"fmt"
	"sort"
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
	seen   map[string]bool
}

// New 建一个空集合。
func New() *Set {
	return &Set{parent: map[string]string{}, size: map[string]int{}, seen: map[string]bool{}}
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
	root := x
	for s.parent[root] != root {
		c.Scanned++
		root = s.parent[root]
	}
	for s.parent[x] != root {
		next := s.parent[x]
		s.parent[x] = root
		x = next
	}
	return root
}

// Union 合并两个成员所在的集合。
func (s *Set) Union(a, b, reqID string, c *Counter) {
	if reqID != "" {
		if s.seen[reqID] {
			return
		}
		s.seen[reqID] = true
	}
	ra := s.find(a, c)
	rb := s.find(b, c)
	if ra == rb {
		return
	}
	if s.size[ra] < s.size[rb] || (s.size[ra] == s.size[rb] && rb < ra) {
		ra, rb = rb, ra
	}
	s.parent[rb] = ra
	s.size[ra] += s.size[rb]
	s.merges++
}

// Same 报告两个成员是否同组。
func (s *Set) Same(a, b string, c *Counter) bool {
	return s.find(a, c) == s.find(b, c)
}

// Merges 返回已经生效的合并次数。
func (s *Set) Merges() int { return s.merges }

// Each 逐个吐出分组：组之间按组内最小成员 id 升序，组内成员按 id 升序。
func (s *Set) Each(fn func([]string), c *Counter) {
	byRoot := map[string][]string{}
	for _, x := range s.order {
		r := s.find(x, c)
		byRoot[r] = append(byRoot[r], x)
	}
	groups := make([][]string, 0, len(byRoot))
	for _, members := range byRoot {
		sort.Strings(members)
		groups = append(groups, members)
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i][0] < groups[j][0] })
	for _, group := range groups {
		fn(group)
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
