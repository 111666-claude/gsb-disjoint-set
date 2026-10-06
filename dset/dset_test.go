package dset

import "testing"

func TestSingleMemberIsItsOwnGroup(t *testing.T) {
	s := New()
	if got := s.find("a1", &Counter{}); got != "a1" {
		t.Fatalf("单独成员应是自己的根，得到 %s", got)
	}
}

func TestUnionJoinsGroups(t *testing.T) {
	s := New()
	c := &Counter{}
	s.Union("a1", "a2", "r1", c)
	if s.find("a1", c) != s.find("a2", c) {
		t.Fatal("合并后两个成员应在同一组")
	}
}

func TestUnionCountsOnce(t *testing.T) {
	s := New()
	s.Union("a1", "a2", "r1", &Counter{})
	if s.Merges() != 1 {
		t.Fatalf("合并一次应记 1，得到 %d", s.Merges())
	}
}

func TestEachVisitsEveryGroup(t *testing.T) {
	s := New()
	c := &Counter{}
	s.Union("a1", "a2", "r1", c)
	s.Union("b1", "b2", "r2", c)
	groups := 0
	s.Each(func([]string) { groups++ }, c)
	if groups != 2 {
		t.Fatalf("应有两个组，得到 %d", groups)
	}
}

func TestRenderMentionsGroups(t *testing.T) {
	s := New()
	c := &Counter{}
	s.Union("a1", "a2", "r1", c)
	if got := s.Render(c); got == "" {
		t.Fatal("报告不该为空")
	}
}
