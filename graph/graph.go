package graph

import (
	"container/list"
	"fmt"
	"go-parser/tool"
)

type Graph struct{
    // 這次的 list 改用 map 來存，其中以 Vertex 來做索引
    // 突破只能用 int 來當 key 的限制
    adj_list map[*Vertex](*list.List)
}

type Vertex struct{
    Val interface{}    // Val 改為任意型態
}

// 創件 graph 時不在需要告訴有幾個節點
// 新增節點由 NewVertex 來添加
func NewGraph() *Graph{
    g := new(Graph)
    g.adj_list = make(map[*Vertex](*list.List))
    return g
}

// 新增節點的步驟改成綁定 Graph
// 利用 map 的特性來當成 adjacency list 最左側的陣列，可以更方便地處理新節點
func (g *Graph) NewVertex(number int) *Vertex{
    v := new(Vertex)
    v.Val = number
    g.adj_list[v] = list.New()
    return v
}

func (g *Graph) AddEdge(start, end *Vertex){
    g.adj_list[start].PushBack(end)
    g.adj_list[end].PushBack(start)
}

func (g *Graph) EasyTraversal(){
    // k's type is *Vertex
    // v's type is *list.List
    for k, v := range g.adj_list{
        fmt.Printf("[%v] -> ", k.Val)
        for current := v.Front(); current != nil; current = current.Next(){
            fmt.Printf("%v -> ", current.Value.(*Vertex).Val)
        }
        fmt.Println("nil")
    }
}

func (g *Graph) DFS(start *Vertex){
    // 記錄誰有走訪誰沒有走訪
    visited := make(map[*Vertex]bool)
    // 宣告一個堆疊以方便回去前一步
    path_stack := StackNew()

    // 處理頭點
    path_stack.Push(start)
    visited[start] = true
    fmt.Printf("%v -> ", start.Val)

    // 疊出堆疊 第 12 行
    for v := path_stack.Pop(); v != nil; v = path_stack.Pop(){ // 第 13 行
        // 透過節點來走訪串列
        // element 是 list.Element
        for element := g.adj_list[v.(*Vertex)].Front(); element != nil; element = element.Next(){
            adj_v := element.Value.(*Vertex)
            // 若相鄰節點尚未被走訪
            if visited[adj_v] == false{
                // 原本第 13 行疊出的節點要記得疊入回去
                path_stack.Push(v)
                path_stack.Push(adj_v)
                visited[adj_v] = true
                fmt.Printf("%v -> ", adj_v.Val)
                // 離開走訪串列(因為是深度優先啊，所以找到新點後就改以新點為優先)
                // 而不是沉溺在這裡面
                // 回到走訪 g.adj_list
                break
            }
        }
    }
}

func Dfs(graph map[string][]string, start string, visited []string, vidList *[][]string) {
    visited = append(visited, start)
    for _, ele := range graph[start] {
        if tool.StrContains(visited, ele) == -1 {
            viss := make([]string, len(visited))
            copy(viss, visited)
            Dfs(graph, ele, viss, vidList)
        }
    }
    *vidList = append(*vidList, visited)
}