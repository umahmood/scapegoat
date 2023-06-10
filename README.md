# Scapegoat Tree

A Go library which implements the Scapegoat Tree data structure.

A Scapegoat Tree is a self-balancing binary search tree, which is based on the
concept of a scapegoat. A scapegoat is typically a person who is blamed when
there is a problem, and Scapegoat Trees use this term because they identify a
node to “blame” when the tree becomes unbalanced. The tree could become
unbalanced after an insertion or deletion of a node, and then an element would
be identified as “the scapegoat” to accept the problem. This problem would be
corrected by rebalancing the subtree at the scapegoat

A Scapegoat Tree provides worst-case O(log n) lookup time and O(log n) amortised
insertion and deletion time.

## Installation

> go get github.com/umahmood/scapegoat

## Usage

```
packgage main

import (
	"fmt"

	"github.com/umahmood/scapegoat"
)

func main() {
    tree, err := scapegoat.New[int](scapegoat.DefaultAlpha)
    // handle err
    keys := []int{42, 99, 3}
    for _, key := range keys {
    	err := sg.Insert(key)
        // handle err
    }
    found := tree.Search(42)
    fmt.Println("found 42?", found) // true
    removed := tree.Remove(88)
    fmt.Println("removed 88?", removed) // false
    // stats
    fmt.Println(tree.Stats.TotalInserts)  // 3
    fmt.Println(tree.Stats.TotalRemovals) // 0
    fmt.Println(tree.Stats.TotalSearches) // 1
}
```

## Documentation

- [https://pkg.go.dev/github.com/umahmood/scapegoat](https://pkg.go.dev/github.com/umahmood/scapegoat)

## References

- [Open Data Structures - Scapegoat Trees (Chapter 8)](https://opendatastructures.org/ods-java.pdf)

- [Wikipedia - Scapegoat Tree](https://en.wikipedia.org/wiki/Scapegoat_tree)

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
