// @@
// @ Author       : Eacher
// @ Date         : 2023-03-08 14:09:25
// @ LastEditTime : 2023-03-10 16:33:54
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  : 
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /array/array.go
// @@
package array

import(
    "sort"
    "sync"
    "reflect"
)

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

type Arrays[K Ordered, V any] struct {
	mutex 	sync.RWMutex
	len 	int
	isSort 	bool
	list 	[]Node[K,V]
}

type Node[K Ordered, V any] struct {
	Key 	K
	Value 	V
}

func New[K Ordered, V any]() *Arrays[K,V] { return &Arrays[K,V]{list: make([]Node[K,V], 0)} }

func (a *Arrays[K,V]) Len() int {
	return a.len
}

func (a *Arrays[K,V]) Swap(i, j int) {
	if a.isSort {
		a.list[i], a.list[j] = a.list[j], a.list[i]
	}
}

func (a *Arrays[K,V]) Less(i, j int) (ok bool) {
	defer func () {
		if err := recover(); err != nil {
			ok = false
		}
	}()
	if a.isSort {
		if reflect.ValueOf(a.list[i].Value).Comparable() && reflect.ValueOf(a.list[j].Value).Comparable() {
			ok = a.list[i].Key > a.list[j].Key
		}
	}
    return ok
}

func (a *Arrays[K,V]) Sort() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.isSort = true
	sort.Sort(a)
	a.isSort = false
}

/*-----------------------------------------------------------------数组的泛型方法 start ---------------------------------------------------*/

func (a *Arrays[K,V]) getIndex(key K) (i int, ok bool) {
	a.mutex.RLock()
	defer func () {
		if err := recover(); err != nil {
			i, ok = -1, false
		}
		a.mutex.RUnlock()
	}()
	for v, n := range a.list {
		if ok = reflect.ValueOf(key).Equal(reflect.ValueOf(n.Key)); ok {
			i = v
			break
		}
	}
	return
}

func (a *Arrays[K,V]) Exist(key K) (ok bool) {
	_, ok = a.getIndex(key)
	return
}

func (a *Arrays[K,V]) First() *Node[K,V] {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	if a.len > 0 {
		n := a.list[0]
		return &n
	}
	return nil
}

func (a *Arrays[K,V]) Last() *Node[K,V] {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	if a.len > 0 {
		n := a.list[a.len-1]
		return &n
	}
	return nil
}

func (a *Arrays[K,V]) PushFront(key K, value V) *Node[K,V] {
	if a.len > 65534 {
		return nil
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	n, list := Node[K,V]{Key: key, Value: value}, make([]Node[K,V], a.len+1)
	list[0] = n
	if a.len > 0 {
		copy(list[1:], a.list)
	}
	a.list = list
	a.len++
	return &n
}

func (a *Arrays[K,V]) PushBack(key K, value V) *Node[K,V] {
	if a.len > 65534 {
		return nil
	}
	n := Node[K,V]{Key: key, Value: value}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.list = append(a.list, n)
	a.len++
	return &n
}

func (a *Arrays[K,V]) Get(key K) (n Node[K,V], ok bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	var i int
	if i, ok = a.getIndex(key); ok {
		n = a.list[i]
	}
	return
}

func (a *Arrays[K,V]) Gets(keys ...K) ([]Node[K,V], bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	nodes, ok := make([]Node[K,V], 0, 1), false
	for _, k := range keys {
		if i, o := a.getIndex(k); o {
			ok = true
			nodes = append(nodes, a.list[i])
		}
	}
	return nodes, ok
}

func (a *Arrays[K,V]) Set(node Node[K,V]) (ok bool) {
	var i int
	if i, ok = a.getIndex(node.Key); ok {
		a.mutex.Lock()
		a.list[i] = node
		a.mutex.Unlock()
	}
	return
}

func (a *Arrays[K,V]) Sets(nodes ...Node[K,V]) ([]Node[K,V], bool) {
	values, ok := make([]Node[K,V], 0, 1), false
	for _, node := range nodes {
		if o := a.Set(node); o {
			ok = true
			values = append(values, node)
		}
	}
	return values, ok
}

func (a *Arrays[K,V]) Delete(key K) (n Node[K,V], ok bool) {
	var i int
	if i, ok = a.getIndex(key); ok {
		a.mutex.Lock()
		list := make([]Node[K,V], a.len-1)
		n = a.list[i]
		a.len--
		copy(list[0:i], a.list[0:i])
		copy(list[i:], a.list[i+1:])
		a.list = list
		a.mutex.Unlock()
	}
	return
}

func (a *Arrays[K,V]) Deletes(keys ...K) ([]Node[K,V], bool) {
	nodes, ok := make([]Node[K,V], 0, 1), false
	for _, k := range keys {
		if n, o := a.Delete(k); o {
			ok = true
			nodes = append(nodes, n)
		}
	}
	return nodes, ok
}

// 返回数组节点切片
func (a *Arrays[K,V]) GetNodes() []Node[K,V] {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	list := make([]Node[K,V], a.len)
	copy(list, a.list)
	return list
}

func (a *Arrays[K,V]) Clear() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.list = a.list[0:0]
	a.len = 0
	return true
}

func (a *Arrays[K,V]) Copy() *Arrays[K,V] {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	a1 := &Arrays[K,V]{list: make([]Node[K,V], a.len), len: a.len}
	copy(a1.list, a.list)
	return a1
}

// 合并数组将会把第一组之后的数组内容全部完整的复制到第一组内；相同键名的将以最后一个为准
func (a *Arrays[K,V]) Merges(arrs ...*Arrays[K,V]) {
	for _, v := range arrs {
		for _, n := range v.GetNodes() {
			i, ok := a.getIndex(n.Key)
			a.mutex.Lock()
			if !ok {
				a.list, i = append(a.list, n), a.len
				a.len++
			}
			a.list[i] = n
			a.mutex.Unlock()
		}
	}
}

// 切分数组 start 开始 size 长度
func (a *Arrays[K,V]) Slice(start, size int) *Arrays[K,V] {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var a1 *Arrays[K,V]
	if (a.len - start) > size {
		a1 = &Arrays[K,V]{list: make([]Node[K,V], size)}
		a1.len = size
		copy(a1.list, a.list[start:start+size])
		list := make([]Node[K,V], start)
		copy(list, a.list[:start])
		a.list, a.len = list, a.len - start
	}
	return a1
}

/*-----------------------------------------------------------------数组的泛型方法 end   ---------------------------------------------------*/