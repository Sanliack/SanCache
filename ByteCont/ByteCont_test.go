package ByteCont

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func hh(g Getter, key string) {
	data, err := g.Get(key)
	fmt.Println(data, err)
}

func Test_interFunc(t *testing.T) {
	hh(GetterFunc(func(key string) ([]byte, error) {
		return nil, errors.New("aa")
	}), "11")
}

func Test_Getter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		fmt.Println(key + "AAAAAAAA")
		return []byte(key), nil
	})
	ex := []byte("key1")
	if v, _ := f.Get("key1"); !reflect.DeepEqual(ex, v) {
		fmt.Println(ex, "-------", v)
		fmt.Println("errrrrr")
	}
}

func Test_Group(t *testing.T) {
	gg := NewGroup("test1.1", 20, GetterFunc(func(key string) ([]byte, error) {
		fmt.Println("use getter func ")
		if key == "key3" {
			return []byte("val3"), nil
		}
		return []byte{}, errors.New("find error")
	}))

	gg.cache.add("key1", ByteView{
		b: []byte("val1"),
	})
	gg.cache.add("key2", ByteView{
		b: []byte("val2"),
	})

	for i := 1; i <= 4; i++ {
		key := fmt.Sprintf("key%d", i)
		bv, err := gg.Get(key)
		fmt.Println("get ", key, "====", bv.String(), err)
	}

}

//type T1 interface {
//	A(a int) int
//}
//
//type AB struct {
//}
//
//func (ab *AB) A(a int) int {
//	fmt.Println("++++++++++++++++++++")
//	return 2
//}
//
//func Test_lxdy(t *testing.T) {
//	var gg T1 = &AB{}
//	ab := gg.(*AB)
//	fmt.Println(ab)
//}
