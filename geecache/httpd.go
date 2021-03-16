package geecache
import(
	_ "fmt"
	"net/http"
	"strings"
)

type CacheHttp struct{
	ipaddr   string
}

func NewCacheHttp(ip string) *CacheHttp {

	return &CacheHttp{
		ipaddr : ip,
	}
}


func (*CacheHttp)ServeHTTP(w http.ResponseWriter, req *http.Request) {

	path := strings.SplitN(req.URL.Path, "/", 3)
	groupName := path[1]
	key := path[2]

	group := GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group: " + groupName, http.StatusNotFound)
		return
	}

	v, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(v.ByteSlice())
}