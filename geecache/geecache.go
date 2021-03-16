package geecache
import(
	"fmt"
	"log"
	"sync"
	"geecache/singleflight"
	"time"
)

type Group struct{
	name       string
	getter     Getter
	mainCache  cache
	loader	   *singleflight.Group
}

type Getter interface{
	Get(key string)( []byte, error)
}

type GetterFunc func(key string)( []byte, error)

func (f GetterFunc) Get(key string)( []byte, error){
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {

	if getter == nil {
		panic("getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:	   &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func(g *Group) Get(key string) (ByteView, error){
	if key == "" {
		return ByteView{}, fmt.Errorf("key is empty")
	}

	if value, ok := g.mainCache.get(key); ok{
		log.Println("geecache hit:", key)
		return value, nil
	}

	log.Println("geecache not hit！！:", key)
	return g.load(key)
}

func(g *Group) load(key string) (ByteView, error){
	
	//todo 是否应该从远程节点获取
	
	//远程节点也没有，则通过个回调函数获取
	v, err := g.loader.Do(key, func()(interface{}, error){
								return g.getLocally(key)
					   })

	return v.(ByteView), err
	
}


func(g *Group) getLocally(key string) (ByteView, error){
	
	log.Println("getLocally key:", key)
	log.Println("start sleep  10s")
	time.Sleep(10*time.Second)
	log.Println("stop sleep  10s")

	vByte, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(vByte)}
	g.popucate(key, value)

	return value, nil
}

func(g *Group) popucate(key string, value ByteView) {
	g.mainCache.add(key, value)
}
