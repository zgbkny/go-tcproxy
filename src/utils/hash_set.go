package utils


type HashSet struct {
	m map[interface{}]bool
}

func NewHashSet() *HashSet { // 返回*HashSet可以调用类型为HashSet或*HashSet的所有方法
	return &HashSet{m: make(map[interface{}]bool)}
}

func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e)
}

func (set *HashSet) Clear() {
	// 如果是直接遍历map并清除所有的k，v，那么在并发程序中，可能删除以后又被添加进去
	set.m = make(map[interface{}]bool)
}

func (set *HashSet) Len() int {
	return len(set.m)
}

func (set *HashSet) Contains(e interface{}) bool {
	if set.m[e] {
		return true
	}
	return false
}

func (set *HashSet) Same(other *HashSet) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}