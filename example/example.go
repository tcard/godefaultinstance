// Package example is an example package for godefaultinstance.
package example

//go:generate godefaultinstance -type=MyType
//go:generate godefaultinstance -type=*MyOtherType -name=unexportedMyOtherType

var DefaultMyType = MyType{1, 2}

type MyType struct {
	a, b int
}

func (m MyType) Avg() int {
	return m.sum() / 2
}

func (m MyType) sum() int {
	return m.a + m.b
}

func (m *MyType) Swap() {
	m.a, m.b = m.b, m.a
}

type MyOtherType struct {
	c, d int
}

func (m MyOtherType) Max() int {
	if m.c > m.d {
		return m.c
	}
	return m.d
}

func (m *MyOtherType) Set(c, d int) {
	m.c, m.d = c, d
}

func (m MyOtherType) BothIn(xs ...int) bool {
	cin, din := false, false
	for _, x := range xs {
		if x == m.c {
			cin = true
		}
		if x == m.d {
			din = true
		}
		if cin && din {
			return true
		}
	}
	return false
}
