// 该小程序用于扫描 content 下所有 md 文件的 http 链接的可用性。
// 并将扫描结果 md 文件输出至 result 目录下。
package main

import (
	"detection/links"
)

func main() {
	links.Start("")
}
