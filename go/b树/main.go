package main

// b树
// 1. 每个节点最多有m个子节点
// 2. 除根节点和叶子节点外，其他每个节点至少有m/2个子节点
// 3. 若根节点不是叶子节点，则至少有两个子节点
// 4. 所有叶子节点都在同一层
// 5. 每个节点存放n个关键字，n+1个指向子节点的指针
// 6. 关键字按递增顺序排列，位于子节点指针的左边
const (
	order = 3
)

type Node struct {
	keys     []int
	children []*Node
	isLeaf   bool
}

func NewNode() *Node {
	return &Node{
		keys:     make([]int, 0),
		children: make([]*Node, 0),
		isLeaf:   true,
	}
}

type BTree struct {
	root *Node
}

func NewBTree() *BTree {
	return &BTree{
		root: NewNode(),
	}
}

func (b *BTree) Search(key int) *Node {
	return b.searchRec(b.root, key)
}

func (b *BTree) searchRec(node *Node, key int) *Node {
	i := 0
	for ; i < len(node.keys); i++ {
		if node.keys[i] >= key {
			break
		}
	}
	if i < len(node.keys) && node.keys[i] == key {
		return node
	}
	if node.isLeaf {
		return nil
	}
	return b.searchRec(node.children[i], key)
}
