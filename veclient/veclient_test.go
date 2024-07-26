package veclient

import "testing"

func TestSomething(t *testing.T) {
	m := []int{1, 2, 3, 4, 5}
	p := len(m)
	n1 := m[:p]
	n2 := m[p:]
	for _, v := range n1 {
		t.Log("n1 item", v)
	}
	for _, v := range n2 {
		t.Log("n2 item", v)
	}
}
