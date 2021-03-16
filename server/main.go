package main
import(
"geecache"
"fmt"
"log"
"net/http"
)

var localDB = map[string]string{
	"tom" : "100",
	"bob" : "200",
	"ami" : "500",
}

func getFromLocal(key string)([]byte, error){
	if v, ok := localDB[key]; ok {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("get from local error")
}

func main(){

   fmt.Println("hhaa")

   geecache.NewGroup("score", 1000, geecache.GetterFunc(getFromLocal) )

/*
   v, err := scoreGroup.Get("tom")
   if err != nil {
   	 fmt.Println("1:get tom from group err:", err)
   }else{
   	 fmt.Println("1:get tom from group success:", v.String())
   }


   v, err = scoreGroup.Get("tom")
   if err != nil {
   	 fmt.Println("2:get tom from group err:", err)
   }else{
   	 fmt.Println("2: get tom from group success:", v.String())
   }
*/

   ipaddr := ":6665" //  ipaddr := ":6663" //  ipaddr := ":6661"
   h := geecache.NewCacheHttp(ipaddr)
   log.Fatal(http.ListenAndServe(ipaddr, h))

   //curl  http://127.0.0.1:6665/score/tom

}