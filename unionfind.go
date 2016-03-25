package compiler

// Author: Ersi Ni
// This is mainly copy paste of my own unionfind implementation, see details
// at https://github.com/nilbot/algo.go/blob/master/connectivity/unionfind.go
// and https://blog.nilbot.net/2015/02/go-percolation-threshold/

type unionfind struct {
	parents []int
	size    []int
	cnt     int
}

func (u *unionfind) Find(p int) int {
	return u.root(p)
}

func (u *unionfind) Count() int {
	return u.cnt
}

func (u *unionfind) Connected(a, b int) bool {
	return u.root(a) == u.root(b)
}

func (u *unionfind) root(position int) int {
	for position != u.parents[position] {
		u.parents[position] = u.parents[u.parents[position]]
		position = u.parents[position]
	}
	return position
}

func (u *unionfind) Union(a, b int) {
	i := u.root(a)
	j := u.root(b)
	if i == j {
		return
	}
	if u.size[i] < u.size[j] {
		u.parents[i] = j
		u.size[j] += u.size[i]
	} else {
		u.parents[j] = i
		u.size[i] += u.size[j]
	}
	u.cnt--
}

// UF interface for other type to embed unionfind
type UF interface {
	Union(a, b int)
	Find(a int) int
	Connected(a, b int) bool
	Count() int
}
