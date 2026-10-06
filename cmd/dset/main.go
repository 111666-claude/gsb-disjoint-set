// Command dset 是集合合并与分组报告的场景入口。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"example.com/disjoint-set/dset"
)

func emit(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func groupsOf(s *dset.Set, c *dset.Counter) [][]string {
	var out [][]string
	s.Each(func(group []string) {
		out = append(out, group)
	}, c)
	return out
}

func main() {
	sample := flag.String("sample", "", "order|members|same|replay|work")
	flag.Parse()

	switch *sample {
	case "order":
		s := dset.New()
		c := &dset.Counter{}
		s.Union("z1", "z2", "r1", c)
		s.Union("a1", "a2", "r2", c)
		emit(groupsOf(s, c))
	case "members":
		s := dset.New()
		c := &dset.Counter{}
		s.Union("m2", "m1", "r1", c)
		emit(groupsOf(s, c))
	case "same":
		s := dset.New()
		c := &dset.Counter{}
		s.Union("b1", "a1", "r1", c)
		s.Union("c1", "b1", "r2", c)
		emit(map[string]bool{"same": s.Same("a1", "c1", c)})
	case "replay":
		s := dset.New()
		c := &dset.Counter{}
		s.Union("a1", "a2", "r1", c)
		s.Union("a3", "a2", "r1", c)
		emit(map[string]int{"merges": s.Merges()})
	case "work":
		s := dset.New()
		c := &dset.Counter{}
		for i := 0; i < 2000; i++ {
			s.Union(fmt.Sprintf("m%d", i+1), fmt.Sprintf("m%d", i), fmt.Sprintf("r%d", i), c)
		}
		s.Render(c)
		emit(map[string]int{"merges": s.Merges(), "scanned": c.Scanned})
	default:
		fmt.Fprintln(os.Stderr, "未知场景："+*sample)
		os.Exit(2)
	}
}
