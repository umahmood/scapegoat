package scapegoat

import (
	"errors"
	"math"

	"golang.org/x/exp/constraints"
)

var (
	ConstraintErr = errors.New("constraint error.")
	AlphaValueErr = errors.New("alpha value must be greater than zero.")
)

type node[T constraints.Ordered] struct {
	left   *node[T]
	right  *node[T]
	parent *node[T]
	key    T
}

type stats struct {
	TotalRebalances            uint64
	TotalRebalancesAfterInsert uint64
	TotalRebalancesAfterRemove uint64
	TotalInserts               uint64
	TotalRemovals              uint64
	TotalSearches              uint64
}

type Scapegoat[T constraints.Ordered] struct {
	root  *node[T]
	alpha float64
	n     int
	q     int
	Stats stats
}

const DefaultAlpha float64 = 1.5

// New creates a new Scapegoat instance. The parameter alpha allows
// flexibility in deciding how balanced the tree should be. A high
// alpha value results in fewer balances, making insertion quicker
// but lookups and deletions slower, and vice versa for a low alpha.
func New[T constraints.Ordered](alpha float64) (*Scapegoat[T], error) {
	if alpha <= 0 {
		return nil, AlphaValueErr
	}
	return &Scapegoat[T]{
		alpha: alpha,
	}, nil
}

// Insert a key into the tree. The key will not inserted if it already exists.
func (s *Scapegoat[T]) Insert(key T) error {
	newNode := &node[T]{key: key}
	depth := 0
	// base case - key already exists or nothing in the tree
	if s.searchHelper(key, false) {
		return nil
	} else if s.root == nil {
		s.root = newNode
	} else {
		// search to find the nodes correct position
		currentNode := s.root
		var potentialParent *node[T]
		for currentNode != nil {
			potentialParent = currentNode
			if newNode.key < currentNode.key {
				currentNode = currentNode.left
			} else {
				currentNode = currentNode.right
			}
			depth += 1
		}
		// assign parents and siblings to the new node
		newNode.parent = potentialParent
		if newNode.key < newNode.parent.key {
			newNode.parent.left = newNode
		} else {
			newNode.parent.right = newNode
		}
		newNode.left = nil
		newNode.right = nil
	}
	s.Stats.TotalInserts += 1
	s.n += 1
	s.q += 1
	if float64(depth) > logToBase(float64(s.q), s.alpha) {
		// there is a scapegoat in the tree
		scapegoat := s.findScapegoat(newNode)
		if scapegoat == nil {
			return nil
		}
		s.Stats.TotalRebalancesAfterInsert += 1
		tmp := s.rebalance(scapegoat)
		// assign the correct pointers to and from the scapegoat
		scapegoat.left = tmp.left
		scapegoat.right = tmp.right
		scapegoat.key = tmp.key
		if scapegoat.left != nil {
			scapegoat.left.parent = scapegoat
		}
		if scapegoat.right != nil {
			scapegoat.right.parent = scapegoat
		}
		// n and q must obey the following inequalities at all times
		if !(s.q/2 <= s.n && s.n <= s.q) {
			return ConstraintErr
		}
	}
	return nil
}

// Remove a key from the tree. Returns true if the key was removed
// otherwise false.
func (s *Scapegoat[T]) Remove(key T) bool {
	getMinKey := func(node *node[T]) *node[T] {
		for node.left != nil {
			node = node.left
		}
		return node
	}

	var remove func(T) bool

	remove = func(key T) bool {
		var (
			parent  *node[T]
			current *node[T] = s.root
		)
		for current != nil && current.key != key {
			parent = current
			if key < current.key {
				current = current.left
			} else {
				current = current.right
			}
		}
		if current == nil {
			return false
		}
		// case 1 - node to be deleted has no children
		if current.left == nil && current.right == nil {
			if current != s.root {
				if parent.left == current {
					parent.left = nil
				} else {
					parent.right = nil
				}
			} else {
				s.root = nil
			}
		} else if current.left != nil && current.right != nil {
			// case 2 - node to be deleted has two children
			successor := getMinKey(current.right)
			curKey := successor.key
			remove(successor.key)
			current.key = curKey
		} else {
			// case 3 - node to be deleted has one child
			var child *node[T]
			if current.left != nil {
				child = current.left
			} else {
				child = current.right
			}
			if current != s.root {
				if current == parent.left {
					parent.left = child
				} else {
					parent.right = child
				}
			} else {
				s.root = child
			}
		}
		return true
	} // close remove

	removed := remove(key)
	if removed {
		s.Stats.TotalRemovals += 1
		s.n -= 1
		if s.n <= s.q/2 {
			s.Stats.TotalRebalancesAfterRemove += 1
			s.root = s.rebalance(s.root)
			s.q = s.n
		}
	}
	return removed
}

// Search for a key in the tree.
func (s *Scapegoat[T]) Search(key T) bool {
	return s.searchHelper(key, true)
}

func (s *Scapegoat[T]) searchHelper(key T, logStat bool) bool {
	currentNode := s.root
	found := false
	for currentNode != nil {
		if currentNode.key == key {
			found = true
			break
		}
		if key < currentNode.key {
			currentNode = currentNode.left
		} else {
			currentNode = currentNode.right
		}
	}
	if logStat {
		s.Stats.TotalSearches += 1
	}
	return found
}

func (s *Scapegoat[T]) findScapegoat(node *node[T]) *node[T] {
	for 3*s.sizeOfSubtree(node) <= 2*s.sizeOfSubtree(node.parent) {
		node = node.parent
	}
	return node
}

func (s *Scapegoat[T]) sizeOfSubtree(node *node[T]) int {
	if node == nil {
		return 0
	}
	return 1 + s.sizeOfSubtree(node.left) + s.sizeOfSubtree(node.right)
}

func (s *Scapegoat[T]) rebalance(root *node[T]) *node[T] {
	s.Stats.TotalRebalances += 1
	nodes := &[]node[T]{}
	flatten(root, nodes)
	return buildTreeFromSortedList(nodes, 0, len((*nodes))-1)
}

func flatten[T constraints.Ordered](node *node[T], nodes *[]node[T]) {
	if node == nil {
		return
	}
	flatten(node.left, nodes)
	(*nodes) = append((*nodes), *node)
	flatten(node.right, nodes)
}

func buildTreeFromSortedList[T constraints.Ordered](nodes *[]node[T], start, end int) *node[T] {
	if start > end {
		return nil
	}
	// mid := int(math.Ceil(float64(start + (end-start)/2.0)))
	mid := start + (end-start)/2.0
	node := &node[T]{key: (*nodes)[mid].key}
	node.left = buildTreeFromSortedList(nodes, start, mid-1)
	node.right = buildTreeFromSortedList(nodes, mid+1, end)
	return node
}

// logToBase returns the logarithm of x to the given base b, calculated as
// log(x)/log(base).The argument passed must not be negative, otherwise it
// returns NaN.
func logToBase(x, b float64) float64 {
	return math.Log(x) / math.Log(b)
}
