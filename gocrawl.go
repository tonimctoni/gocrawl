package main

import "runtime"
import "time"
import "fmt"
import "sync/atomic"
import "io/ioutil"
import "hash/fnv"

// func consumer() {
//     //Fuck go's restrictions
//     fmt.Println("")
//     strings.HasPrefix("", "")
//     http.Get("")
//     ioutil.ReadAll(nil)
//     errors.New("")
//     regexp.MustCompile("")
//     log.Fatal("")
//     runtime.Gosched()
// }

import "sync"
import "os"
type MyThreadSafeFile struct{
    mutex sync.Mutex
    file *os.File
}

func NewMyThreadSafeFile(filename string) (*MyThreadSafeFile, error){
    file, err:=os.Create(filename)
    if err!=nil{
        return nil, err
    }
    return &MyThreadSafeFile{sync.Mutex{}, file},nil
}

func (m *MyThreadSafeFile) close() {
    m.file.Close()
}

func (m *MyThreadSafeFile) write(url string, urls []string) error{
    m.mutex.Lock()
    defer m.mutex.Unlock()

    _,err:=m.file.WriteString(url+"\n")
    if err!=nil{
        return err
    }
    for _,url:=range urls{
        _,err:=m.file.WriteString("\t"+url+"\n")
        if err!=nil{
            return err
        }
    }
    return nil
}


type CssManager struct{
    counter uint32
    csshash_set ThreadSafeUint32Set
}

func NewCssManager() *CssManager{
    return &CssManager{0,*NewThreadSafeUint32Set()}
}

func (c *CssManager) manage(content []byte) error{
    h:=fnv.New32a()
    h.Write(content)
    content_hash:=h.Sum32()
    if c.csshash_set.add(content_hash){
        return nil
    }
    new_counter_value:=atomic.AddUint32(&c.counter, 1)
    return ioutil.WriteFile(fmt.Sprintf("css/css%05d.css", new_counter_value), content, 0644)
}


//Goal: make list of visited websites and their contained urls
//also get css files :D
func main() {
    initialize_global_regexp()
    visited_sites:=NewThreadSafeStringSet()
    url_queue:=NewThreadSafeStringQueue(1000000)
    url_queue.push("http://cssdb.co")
    output_file,err:=NewMyThreadSafeFile("output_file.txt")
    if err!=nil{
        panic("Could not open (create) output_file.txt")
    }
    defer output_file.close()
    css_manager:=NewCssManager()
    counter:=uint32(0)

    for i:=0;i<120;i++{
        go func(thread_num int){
            for {//j:=0;j<10;j++{
                url,err:=url_queue.pull()
                if err!=nil{
                    runtime.Gosched()
                    continue
                }

                if visited_sites.add(url){
                    continue
                }

                page_content,err:=NewPageContentIfContentType(url, "text/html")
                if err!=nil{
                    page_content,err:=NewPageContentIfContentType(url, "text/css")
                    if err!=nil{
                        continue
                    }
                    err=css_manager.manage(page_content.get_bytes())
                    if err!=nil{
                        fmt.Println(err.Error())
                    }
                    continue
                }

                new_urls:=page_content.get_urls()
                url_queue.push_slice(new_urls)
                output_file.write(url, new_urls)
                new_counter_value:=atomic.AddUint32(&counter, 1)
                fmt.Println(new_counter_value, thread_num, url)
            }
        }(i)
    }
    time.Sleep(10*time.Minute)
    // sq:=NewThreadSafeStringQueue(2)
    // fmt.Println(sq.push("1"))
    // fmt.Println(sq.push("2"))
    // fmt.Println(sq.push("3"))
    // fmt.Println(sq.pull())
    // fmt.Println(sq.pull())
    // fmt.Println(sq.pull())
    // fmt.Println(sq.pull())
    // tsss:=NewThreadSafeStringSet()
    // fmt.Println(tsss.add("hello"))
    // fmt.Println(tsss.add("hello"))
    // fmt.Println(tsss.add("qwe"))
    // sq:=NewThreadSafeStringQueue(4)
    // a:=[]string{"1","2","3"}
    // fmt.Println(sq.report())
    // fmt.Println(sq.push_slice(a))
    // fmt.Println(sq.report())
    // fmt.Println(sq.push_slice(a))
    // fmt.Println(sq.pull())
    // fmt.Println(sq.pull())
    // fmt.Println(sq.pull())
    // fmt.Println(sq.pull())
    // css_manager:=NewCssManager()
    // fmt.Println(css_manager.manage([]byte("hello")))
    // fmt.Println(css_manager.manage([]byte("hello")))
    // fmt.Println(css_manager.manage([]byte("hella")))
    // fmt.Println(css_manager.manage([]byte("hella")))
}