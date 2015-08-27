package main

import "fmt"
import "time"

import "../../safemap"

func main() {

	sm := safemap.New()

	go sm.Insert("test", 111)
	go sm.Insert("abc", 123)
	go sm.Insert("def", 345)
	sm.Insert("test", 789)

	time.Sleep(time.Second * 1)

	fmt.Println("all data in map1: ", sm.Close())

	sm2 := safemap.New()

	sm2.Insert("test", 111)
	sm2.Insert("abc", 123)
	sm2.Insert("def", 345)

	fmt.Println("length of map:", sm2.Len())

	sm2.Delete("test")

	fmt.Println("length of map:", sm2.Len())

	value, found := sm2.Find("abc")
	fmt.Printf("find value of key abc in map  value:%d found:%v\n", value, found)

	sm2.Update("def", func(price interface{}, found bool) interface{} {
		if !found {
			return 0.0
		}

		return float64(price.(int)) * 1.05
	})

	fmt.Println("all data in map2: ", sm2.Close())
}
