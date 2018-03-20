package stock

import (
	"testing"
	"fmt"
)

func Test_fetch(t *testing.T) {
	data, err := fetch("sz000016")

	if err != nil {
		t.Error("fetch 测试未通过");
	}

	str, err := match(data)

	if err != nil {
		t.Error(err.Error())
	}

	for _, val := range str {
		fmt.Println(string(val))
	}
}