package links

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	files     []string
	fileCount int
	link      map[string]string // [path:line]=http:url
	links     map[string]string // [path:line]=https:url

)

func Start(base string) {
	if base == "" {
		base = "content/en"
	}

	fmt.Println("加载文件列表中……")
	files = make([]string, 1000)
	err := GetAllFile(base)
	if err != nil {
		panic(err)
	}
	fmt.Printf("成功加载 %d 个 .md 文件\n", len(files))

	fmt.Println("加载 url 列表中……")
	link = make(map[string]string)
	links = make(map[string]string)
	for k, v := range files {
		if v == "" {
			continue
		}
		GetAllURL(files[k])
	}
	fmt.Printf("成功加载 %d 个 HTTP 链接\n", len(link))
	fmt.Printf("成功加载 %d 个 HTTPS 链接\n", len(links))

	fmt.Println("检测 url 列表中……")
	wg := &sync.WaitGroup{}
	for k, v := range link {
		wg.Add(1)
		time.Sleep(time.Millisecond * 10)
		go handle(k, v, wg, false)
	}
	wg.Wait()

	for k, v := range links {
		wg.Add(1)
		time.Sleep(time.Millisecond * 10)
		go handle(k, v, wg, true)
	}
	wg.Wait()
	out.OutPut()
	fmt.Println("任务完成！")
}

func GetAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			err = GetAllFile(path.Join(pathname, fi.Name()))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			if path.Ext(fi.Name()) == ".md" {
				files[fileCount] = path.Join(pathname, fi.Name())
				fileCount++
				if fileCount+100 > len(files) {
					files = append(files, make([]string, len(files))...)
				}
			}
		}
	}
	return err
}

func GetAllURL(pathname string) {
	f, err := os.Open(pathname)
	if err != nil {
		panic(err)
	}
	line := 0
	//reg:=regexp.MustCompile(`\[(.*?)\]`)
	reg := regexp.MustCompile(`\(http://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	regs := regexp.MustCompile(`\(https://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line++
		lk := reg.FindString(sc.Text())
		if lk != "" {
			link[pathname+":"+strconv.Itoa(line)] = lk[1:]
			continue
		}
		lk = regs.FindString(sc.Text())
		if lk != "" {
			links[pathname+":"+strconv.Itoa(line)] = lk[1:]
			continue
		}
	}
}
