package agent

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTransIDString(t *testing.T) {
	bu := NewByteUtil("")
    bu.PutByte(TVERSION)
	bu.PutPrefixedString("i am cnogo")
	bu.PutVar64(time.Now().UnixNano()/1000)
	bu.PutVar64(-9857834)
	bys := bu.GetBytes()

	fmt.Println(time.Now().UnixNano()/1000)

	pa := NewByteUtil(string(bys))
	fmt.Println(pa.ReadByte())
	fmt.Println(pa.ReadPrefixedString())
	fmt.Println(pa.ReadVar64())
	fmt.Println(pa.ReadVar64())
}



func Test_Str(t *testing.T) {
	m(nil)
}

func m(s interface{}) {
	a, ok := s.(int32)
	if !ok {
		fmt.Println("it is nil ", a)
	}

	fmt.Println(a)
}


