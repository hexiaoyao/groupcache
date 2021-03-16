package main
import(
"fmt"
"hash/crc32"
"sync"
"strconv"
"sort"
"os/exec"
)


const(
	spotPerNodeDefault = 100
)

type node struct{
	nodeKey string
	nodeHashValue uint32
}


//follow the test of sort package
type nodesArray []node

func (a nodesArray) Len() int           { return len(a) }
func (a nodesArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nodesArray) Less(i, j int) bool { return a[i].nodeHashValue < a[j].nodeHashValue }


type hashCons struct{
	spotPerNode uint32
	nodes nodesArray
	mu	sync.RWMutex
}

func (h *hashCons ) addNodes( nodesKey []string){
	h.mu.Lock()
	defer h.mu.Unlock()
	
	for _, nodeKey := range nodesKey {
		var nodeTmp node
		nodeTmp.nodeKey = nodeKey
		h.nodes = append(h.nodes, nodeTmp)
	}

	h.generateHashcons()
}


func (h *hashCons ) generateHashcons() {

	var nodesArraynNew nodesArray
	for _, nodeInfo := range h.nodes {
		for i := 1; i <= int(h.spotPerNode); i++ {
			var nodeNew node
			nodeHashValueKey := nodeInfo.nodeKey + strconv.Itoa(i)
			nodeHashValue := h.genValue(nodeHashValueKey)
			nodeNew.nodeKey = nodeInfo.nodeKey
			nodeNew.nodeHashValue = nodeHashValue

			nodesArraynNew = append(nodesArraynNew, nodeNew)
		}
	}

	sort.Sort(nodesArraynNew)
	h.nodes = nodesArraynNew

}



func (h *hashCons ) genValue( k string ) uint32 {
	return crc32.ChecksumIEEE([]byte(k))
}


func (h *hashCons ) getNode( value string) string {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	valueHash := h.genValue(value)
	i := sort.Search(len(h.nodes), func(i int) bool { return h.nodes[i].nodeHashValue >= valueHash })
	if len(h.nodes) == i {
		i = 0
	}

	return h.nodes[i].nodeKey
	
}


func NewHashCons() *hashCons {

	h := &hashCons{
		spotPerNode: spotPerNodeDefault,
	}
	return h
}



func main(){
	testNodes := []string{"6665", "6663", "6661"}

	testHash := NewHashCons()

	testHash.addNodes(testNodes)

	v1 := "tom"
	get(testHash.getNode(v1), v1)

	v2 := "jock"
	get(testHash.getNode(v2), v2)

	v3:= "bob"
	get(testHash.getNode(v3), v3)

	v5:= "666666666666"
	get(testHash.getNode(v5), v5)

	v6:= "abcd"
	get(testHash.getNode(v6), v6)
}


func get(port string, key string){
	fmt.Println(port, key)

	command := `/usr/bin/curl http://127.0.0.1:` + port + `/score/` + key
	fmt.Println(command)
	cmd := exec.Command("/bin/bash", "-c", command)
	bytes,err := cmd.Output()
	if err != nil {
     	fmt.Println(err)
	}
	resp := string(bytes)
	fmt.Println(resp)

	fmt.Println(">>>>>>>>>>>>>>>>>>>>")

}