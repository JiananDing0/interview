// 现在有107个用户，编号为1- 107，现在已知有m对关系，每一对关系给你两个数x和y，
// 代表编号为x的用户和编号为y的用户是在一个圈子中，例如：A和B在一个圈子中，B和C
// 在一个圈子中，那么A,B,C就在一个圈子中。现在想知道最多的一个圈子内有多少个用户。
//
// 输入例子1:
// 2
// 4
// 1 2
// 3 4
// 5 6
// 1 6
// 4
// 1 2
// 3 4
// 5 6
// 7 8
//
// 输出例子1:
// 4
// 2

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
