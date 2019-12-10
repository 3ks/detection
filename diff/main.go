// 该小程序用于对比 content 下 old.md 和 source.md 文件的 HTTP 链接。
// 并将 old.md 文件中不存在的 http 链接输出至 result 目录下的 new.md。
package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// 加载解析 source.md
	// key:url  value:line
	m404 := Load404()
	// 加载解析 old.md
	old := LoadOld()

	count := 0
	data := make([]string, len(m404))
	for k, v := range m404 {
		if !old[k] {
			data[count] = v + "\n"
			count++
		}
	}
	data = data[:count]
	sort.Slice(data, func(i, j int) bool {
		return data[i] > data[j]
	})

	count = 0
	dataOrder := bytes.Buffer{}
	for _, v := range data {
		dataOrder.WriteString(strconv.Itoa(count+1) + "." + v)
		count++
	}

	_ = os.Mkdir("result", os.ModeDir)
	_ = ioutil.WriteFile("result/new.md", dataOrder.Bytes(), 0644)
}

func Load404() (m404 map[string]string) {
	m404 = make(map[string]string)
	f, err := os.Open("content/source.md")
	if err != nil {
		panic(err)
	}
	//reg:=regexp.MustCompile(`\[(.*?)\]`)
	reg := regexp.MustCompile(`\(http://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	regs := regexp.MustCompile(`\(https://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		ss := strings.Split(sc.Text(), ",")
		lk := reg.FindString(sc.Text())
		if lk != "" {
			m404[lk[1:]] = ss[0] + ", " + lk[1:]
			continue
		}
		lk = regs.FindString(sc.Text())
		if lk != "" {
			m404[lk[1:]] = ss[0] + ", " + lk[1:]
			continue
		}
	}
	return
}

func LoadOld() (old map[string]bool) {
	old = make(map[string]bool)
	f, err := os.Open("content/old.md")
	if err != nil {
		panic(err)
	}
	//reg:=regexp.MustCompile(`\[(.*?)\]`)
	reg := regexp.MustCompile(`http://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	regs := regexp.MustCompile(`https://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lk := reg.FindString(sc.Text())
		if lk != "" {
			old[lk] = true
			continue
		}
		lk = regs.FindString(sc.Text())
		if lk != "" {
			old[lk] = true
			continue
		}
	}
	return
}
