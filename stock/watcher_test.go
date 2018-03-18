package stock

import (
	"testing"
	"fmt"
)

func Test_fetch(t *testing.T) {
	data, err := fetch("http://hq.sinajs.cn/list=sz000001,sh600000")

	if err != nil {
		t.Error("fetch 测试未通过");
	}

	str := format(data)
	fmt.Println(str)
}