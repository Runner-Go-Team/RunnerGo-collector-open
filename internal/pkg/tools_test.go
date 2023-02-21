package pkg

import (
	"fmt"
	"runtime"
	"testing"
)

func TestPost(t *testing.T) {
	//Post("https://dev-tools-env-test.apipost.cn/api/demo/news_list?mobile=18289454846&theme_news=国际新闻&page=1&pageSize=20", "")
	fmt.Println(runtime.NumCPU())
}
