package abac

import (
	"fmt"
	"strings"
)

// 节点（字典树）
type MyNode struct {
	// Data
	children map[interface{}]*MyNode //子节点
	isEnd    bool
}

// 二叉树
type MyTree struct {
	Root *MyNode
}

// 构造访问策略树
func NewTree() *MyTree {
	return &MyTree{Root: NewTreeNode()}
}

// 构造树节点
func NewTreeNode() *MyNode {
	return &MyNode{
		// Data:  data,
		children: make(map[interface{}]*MyNode),
		isEnd:    false,
	}
}

// 向访问策略树插入一个访问策略
func (mt *MyTree) Append(policy []interface{}) {
	node := mt.Root
	for i := 0; i < len(policy); i++ {
		_, ok := node.children[policy[i]]
		// 如果该节点不存在，则构造一个节点
		if !ok {
			node.children[policy[i]] = NewTreeNode()
		}
		node = node.children[policy[i]]
	}
	node.isEnd = true
}

// 搜索树中是否存在指定单词
func (mt *MyTree) Search(policy []interface{}) bool {
	node := mt.Root
	for i := 0; i < len(policy); i++ {
		_, ok := node.children[policy[i]]
		// fmt.Println("ok", policy[i])

		if !ok {
			return false
		}
		node = node.children[policy[i]]
	}
	return node.isEnd
}

// 判断树中是否有指定前缀的单词
func (mt *MyTree) StartsWith(prefix []interface{}) bool {
	node := mt.Root
	for i := 0; i < len(prefix); i++ {
		_, ok := node.children[prefix[i]]
		if !ok {
			return false
		}
		node = node.children[prefix[i]]
	}
	return true
}

// 把policy变成tree存储
func PolicyToTree(policy *Policy) *MyTree {
	// var alltree []*MyTree
	tree := NewTree()
	f := 0
	for i := 0; i < len(policy.SubRules); i++ {
		// if f {
		// 	i = i - 1
		// }
		rule := policy.SubRules[i]
		var value []interface{}
		// fmt.Println(reflect.TypeOf(rule))

		// values := dataProcess(rule)
		values := strings.Split(rule, ",")
		// fmt.Println(values)
		for _, v := range values {
			// fmt.Println(v)
			if v == "action:rw" {
				if f == 0 {
					value = append(value, "action:r")
					// f = true
					f++
					i--
				} else if f == 1 {
					value = append(value, "action:w")
					// f = false
					f = 0
				}
			} else if v == "action:all" {
				if f == 0 {
					value = append(value, "action:r")
					f++
					i--
				} else if f == 1 {
					value = append(value, "action:w")
					f++
					i--
				} else if f == 2 {
					value = append(value, "action:d")
					f = 0
				}

			} else {
				value = append(value, v)
			}
			// fmt.Println(value)
		}
		tree.Append(value)
	}
	// for _, rule := range policy.SubRules {
	// 	var value []interface{}
	// 	// fmt.Println(reflect.TypeOf(rule))

	// 	// values := dataProcess(rule)
	// 	values := strings.Split(rule, ",")
	// 	// fmt.Println(values)
	// 	for _, v := range values {
	// 		fmt.Println(v)
	// 		if v == "action:rw" {
	// 			if !f {
	// 				value = append(value, "action:r")
	// 				f = true
	// 			} else {
	// 				value = append(value, "action:w")
	// 				f = false
	// 			}
	// 		} else {
	// 			value = append(value, v)
	// 		}
	// 		fmt.Println(value)
	// 	}
	// 	tree.Append(value)
	// }
	return tree
}

func PreorderPrint(mn *MyNode) {
	if mn == nil {
		return
	}
	// fmt.Print(len(mn.children), " ")
	fmt.Print(mn.children, " ")
	for _, child := range mn.children {
		// fmt.Println("child is ", child, " ")
		// fmt.Println(mn.children[child], " ")

		PreorderPrint(child)

	}
	// for i := 0; i < len(mn.children); i++ {
	// 	PreorderPrint(mn.children[i])

	// }
}
