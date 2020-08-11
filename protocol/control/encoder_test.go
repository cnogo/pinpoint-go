package control

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_Encoder(testing *testing.T) {
	encoder := NewEncoder()
	testMap := make(map[string]interface{})
	subMap := make(map[string]interface{})
	testMap["name"] = "alias"
	testMap["age"] = 10
	testMap["money"] = 12.15
	testMap["isTop"] = true

	subMap["address"] = "fuzhou"
	subMap["mail-code"] = "365000"
	subMap["pay"] = false
	subMap["time"] = 10.25

	testMap["info"] = subMap

	decoder := NewDecoder(encoder.Encode(testMap))
	m := decoder.Decode()
	fmt.Println(m)
}

func Test_switch(testing *testing.T) {
	a, err := strconv.ParseInt("10", 10, 64)
	fmt.Println(a, " error: ", err)
}
