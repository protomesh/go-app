package structures

func longestPrefix(k1, k2 string) int {
	max := len(k1)
	if l := len(k2); l < max {
		max = l
	}
	var i int
	for i = 0; i < max; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

type RadixTree[T comparable] interface {
	Insert(prefix string, value T) int
	Delete(prefix string) (int, bool)
	Match(prefix string) T
	MatchLongest(prefix string) T
}

type radixTree[T comparable] struct {
	root *radixNode[T]
}

func NewRadixTree[T comparable]() RadixTree[T] {
	var zeroValue T
	return &radixTree[T]{
		root: &radixNode[T]{
			prefix:    "",
			value:     zeroValue,
			children:  make([]*radixNode[T], 0),
			level:     0,
			zeroValue: zeroValue,
		},
	}
}

func (tree *radixTree[T]) Insert(prefix string, value T) int {
	return tree.root.Insert(prefix, value)
}

func (tree *radixTree[T]) Delete(prefix string) (int, bool) {
	return tree.root.Delete(prefix)
}

func (tree *radixTree[T]) Match(prefix string) T {
	return tree.root.Match(prefix)
}

func (tree *radixTree[T]) MatchLongest(prefix string) T {
	return tree.root.MatchLongest(prefix)
}

type radixNode[T comparable] struct {
	prefix    string
	value     T
	children  []*radixNode[T]
	level     int
	zeroValue T
}

func (node *radixNode[T]) Insert(prefix string, value T) int {
	if len(prefix) == 0 {
		node.value = value
		return node.level
	}

	for _, child := range node.children {

		if len(child.prefix) == 0 {
			continue
		}

		commonPrefix := longestPrefix(child.prefix, prefix)

		if commonPrefix == 0 {
			continue
		} else if commonPrefix == len(child.prefix) && commonPrefix == len(prefix) {
			child.value = value
			return child.level
		} else if commonPrefix < len(child.prefix) {
			child.split(commonPrefix)
		}

		return child.Insert(prefix[commonPrefix:], value)

	}

	node.children = append(node.children, &radixNode[T]{
		prefix:    prefix,
		value:     value,
		children:  make([]*radixNode[T], 0),
		level:     node.level + 1,
		zeroValue: node.zeroValue,
	})
	return node.level + 1
}

func (node *radixNode[T]) Delete(prefix string) (int, bool) {
	if len(prefix) == 0 {
		node.value = node.zeroValue
		return node.level, true
	}

	for _, child := range node.children {

		if len(child.prefix) == 0 {
			continue
		}

		commonPrefix := longestPrefix(child.prefix, prefix)

		if commonPrefix == 0 {
			continue
		} else if commonPrefix == len(child.prefix) && commonPrefix == len(prefix) {
			child.value = node.zeroValue
			return child.level, true
		} else if commonPrefix < len(child.prefix) {
			return child.level, false
		}

		return child.Delete(prefix[commonPrefix:])

	}

	return node.level, false
}

func (node *radixNode[T]) Match(prefix string) T {
	if len(prefix) == 0 {
		return node.value
	}

	for _, child := range node.children {

		if len(child.prefix) == 0 {
			continue
		}

		commonPrefix := longestPrefix(child.prefix, prefix)

		if commonPrefix == 0 {
			continue
		} else if commonPrefix == len(child.prefix) && commonPrefix == len(prefix) {
			return child.value
		} else if commonPrefix < len(child.prefix) {
			return child.zeroValue
		}

		return child.Match(prefix[commonPrefix:])

	}

	return node.zeroValue
}

func (node *radixNode[T]) MatchLongest(prefix string) T {
	if len(prefix) == 0 {
		return node.value
	}

	for _, child := range node.children {

		if len(child.prefix) == 0 {
			continue
		}

		commonPrefix := longestPrefix(child.prefix, prefix)

		if commonPrefix == 0 {
			continue
		} else if commonPrefix == len(child.prefix) && commonPrefix == len(prefix) {
			return child.value
		} else if commonPrefix < len(child.prefix) {
			return child.zeroValue
		}

		value := child.MatchLongest(prefix[commonPrefix:])

		if value == child.zeroValue {
			return child.value
		}

		return value

	}

	return node.value
}

func (node *radixNode[T]) split(index int) {
	if index >= len(node.prefix) {
		panic("index out of range")
	}

	prefix := node.prefix[index:]
	node.prefix = node.prefix[:index]

	newNode := &radixNode[T]{
		prefix:    prefix,
		value:     node.value,
		children:  node.children,
		level:     node.level + 1,
		zeroValue: node.zeroValue,
	}

	node.value = node.zeroValue
	node.children = []*radixNode[T]{newNode}
}
