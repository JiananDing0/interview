package main

import "fmt"

type UnionFind struct {
	parent   []int
	setCount []int
	max      int
}

func unionCreate() UnionFind {
	u := UnionFind{}
	u.parent = make([]int, 10)
	u.setCount = make([]int, 10)
	u.max = 0
	for i := 0; i < 10; i++ {
		u.parent[i] = i
		u.setCount[i] = 1
	}
	return u
}

func (u UnionFind) root(i int) int {
	for i != u.parent[i] {
		u.parent[i] = u.parent[u.parent[i]]
		i = u.parent[i]
	}
	return i
}

func (u *UnionFind) union(i, j int) {
	rootI := u.root(i)
	rootJ := u.root(j)
	if rootI == rootJ {
		return
	}
	u.parent[rootI] = rootJ
	u.setCount[rootJ] += u.setCount[rootI]
	if u.setCount[rootJ] > u.max {
		u.max = u.setCount[rootJ]
	}
}

func main() {
	var count, arrLen, numA, numB int
	fmt.Scanln(&count)
	for i := 0; i < count; i++ {
		fmt.Scanln(&arrLen)
		u := unionCreate()
		for j := 0; j < arrLen; j++ {
			fmt.Scanln(&numA, &numB)
			u.union(numA, numB)
		}
		fmt.Println(u.max)
	}
}
