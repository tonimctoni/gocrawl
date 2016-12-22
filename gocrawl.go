package main

import "runtime"
import "time"
import "fmt"
import "sync/atomic"

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


//Goal: make list of visited websites and their contained urls
func main() {
    initialize_global_regexp()
    visited_sites:=NewThreadSafeStringSet()
    url_queue:=NewThreadSafeStringQueue(100000000)
    url_queue.push("http://golang.org")
    output_file,err:=NewMyThreadSafeFile("output_file.txt")
    if err!=nil{
        panic("Could not open (create) output_file.txt")
    }
    defer output_file.close()
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

                page_content,err:=NewPageContentIfHTML(url)
                if err!=nil{
                    continue
                }

                new_urls:=page_content.get_urls()
                url_queue.push_slice(new_urls)
                output_file.write(url, new_urls)
                atomic.AddUint32(&counter, 1)
                fmt.Println(counter, thread_num, url)
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
}