package lru
import(
"testing"
_ "reflect"
"fmt"
)


type String string

func(d String) Len()int{
	return len(d)
}


func TestAdd(t *testing.T) {
	testLru := New(int64(100), nil)
	testLru.Add("key1", String("abc") )
	testLru.Add("key2", String("haha") )

	if testLru.nbytes != int64( len("key1") + len("abc") + len("key2") + len("haha") ){
		fmt.Println(testLru.nbytes)
		t.Fatalf("cache add key error")
	}
}