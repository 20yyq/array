package main

import(
	"fmt"
	"gitee.com/maxrang/array"
)

func main() {
    a := array.New[int,string]()
    a.PushFront(123, "123")
    a.PushFront(1111, "dsgfdg")
    ok2 := a.Set(array.Node[int,string]{Key: 1111, Value: "123421416666"})
    v, ok := a.Get(1111)
    // k1, v1, ok1 := a.Sets(array.Node[any,string]{Key: 1111, Value: "12342141"})
    // fmt.Println("ArraysTest11111111", ok2, v, ok, k1, v1, ok1)
    nodes, ok1 := a.Gets(1111, 123)
    fmt.Println("ArraysTest22222222", ok2, v, ok, nodes, ok1, a.First().Value)
    a.Asc = true
    a.Sort()
    nodes, ok1 = a.Gets(1111, 123)
    fmt.Println("ArraysTest22222222", ok2, v, ok, nodes, ok1, a.First().Value)
}