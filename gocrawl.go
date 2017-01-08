package main

import "runtime"
import "time"
import "fmt"
import "sync/atomic"
import "io/ioutil"
import "hash/fnv"
import "regexp"
import "strings"

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

var re_comments *regexp.Regexp
var re_breaklines *regexp.Regexp
func initialize_global_css_manager_regexp() {
    re_comments=regexp.MustCompile(`/\*(.|\n)*?\*/`)
    re_breaklines=regexp.MustCompile(`\n{3,}`)
}

const allowed_chars = "abcdefghijklmnopqrstuvwxzyABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789\n\t\r \"'(){}[]+-*/.,:;_@#%$!?=\\<>~^|&`"
func (c *CssManager) manage(content []byte) error{
    h:=fnv.New32a()
    h.Write(content)
    content_hash:=h.Sum32()
    if c.csshash_set.add(content_hash){
        return nil
    }
    //Make sure all characters are allowed ones
    outer: for _,char:=range content{
        for _,allowed_char:=range []byte(allowed_chars){
            if char==allowed_char{
                continue outer
            }
        }
        return nil
    }
    //Make some transformations for nicer text
    content=re_comments.ReplaceAll(content, []byte(""))
    content=re_breaklines.ReplaceAll(content, []byte("\n\n"))
    content_as_string:=string(content)
    content_as_string=strings.ToLower(content_as_string)
    content_as_string=strings.TrimSpace(content_as_string)
    //Make sure it is a nice file
    if strings.Count(content_as_string, "\n")<5 || len(content_as_string)<=50{
        return nil
    }
    content=[]byte(content_as_string)

    new_counter_value:=atomic.AddUint32(&c.counter, 1)
    return ioutil.WriteFile(fmt.Sprintf("css/css%06d.css", new_counter_value), content, 0644)
}


//Goal: make list of visited websites and their contained urls
//also get css files :D
func main() {
    initialize_global_regexp()
    initialize_global_css_manager_regexp()
    // http.DefaultTransport.ResponseHeaderTimeout=5*time.Second
    initialize_global_client()
    visited_sites:=NewThreadSafeStringSet()
    url_queue:=NewThreadSafeStringQueue(10000000)// Ten million
    url_queue.push("http://cssdb.co")
    // output_file,err:=NewMyThreadSafeFile("output_file.txt")
    // if err!=nil{
    //     panic("Could not open (create) output_file.txt")
    // }
    // defer output_file.close()
    css_manager:=NewCssManager()
    counter:=uint32(0)

    for i:=0;i<300;i++{
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
                // output_file.write(url, new_urls)
                new_counter_value:=atomic.AddUint32(&counter, 1)
                fmt.Println(new_counter_value, thread_num, url)
            }
        }(i)
    }
    for{
        time.Sleep(10*time.Hour)
    }
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


// panic: runtime error: slice bounds out of range

// goroutine 78 [running]:
// panic(0x619380, 0xc42000c0e0)
//         /usr/local/go/src/runtime/panic.go:500 +0x1a1
// bufio.(*Reader).ReadSlice(0xc5d539d740, 0x41140a, 0xc525183520, 0x20, 0x20, 0x6313e0, 0xc57f5b5a78)
//         /usr/local/go/src/bufio/bufio.go:308 +0x1ee
// net/http/internal.readChunkLine(0xc5d539d740, 0xc52f95bc00, 0x7fdf64f40580, 0x0, 0x0, 0x0)
//         /usr/local/go/src/net/http/internal/chunked.go:110 +0x34
// net/http/internal.(*chunkedReader).beginChunk(0xc46349f740)
//         /usr/local/go/src/net/http/internal/chunked.go:47 +0x32
// net/http/internal.(*chunkedReader).Read(0xc46349f740, 0xc4927c4000, 0x2000, 0x2000, 0xc4e8fe45a0, 0x67758e, 0x479473)
//         /usr/local/go/src/net/http/internal/chunked.go:77 +0x7b
// net/http.(*body).readLocked(0xc4b061f7c0, 0xc4927c4000, 0x2000, 0x2000, 0xc495ad04e0, 0xc52edab780, 0x7fdf60596010)
//         /usr/local/go/src/net/http/transfer.go:651 +0x61
// net/http.bodyLocked.Read(0xc4b061f7c0, 0xc4927c4000, 0x2000, 0x2000, 0x74cd20, 0xc4200e2000, 0xed0009599)
//         /usr/local/go/src/net/http/transfer.go:850 +0x55
// io/ioutil.devNull.ReadFrom(0x0, 0x74d9e0, 0xc4b061f7c0, 0x61da20, 0x47f801, 0x7fdf60596010)
//         /usr/local/go/src/io/ioutil/ioutil.go:151 +0x85
// io/ioutil.(*devNull).ReadFrom(0xc42000c2b0, 0x74d9e0, 0xc4b061f7c0, 0xc57f5b5cf0, 0x1, 0x68ce4473ba26180)
//         <autogenerated>:9 +0x66
// io.copyBuffer(0x74d8e0, 0xc42000c2b0, 0x74d9e0, 0xc4b061f7c0, 0x0, 0x0, 0x0, 0x4a3a81, 0xc57f5b5d80, 0xc57f5b5d90)
//         /usr/local/go/src/io/io.go:384 +0x323
// io.Copy(0x74d8e0, 0xc42000c2b0, 0x74d9e0, 0xc4b061f7c0, 0x0, 0x65aeb4, 0xc52f95bc00)
//         /usr/local/go/src/io/io.go:360 +0x68
// net/http.(*body).Close(0xc4b061f7c0, 0x0, 0x0)
//         /usr/local/go/src/net/http/transfer.go:814 +0x20d
// net/http.(*cancelTimerBody).Close(0xc525183520, 0xc4cf21e510, 0x0)
//         /usr/local/go/src/net/http/client.go:667 +0x34
// main.NewPageContentIfContentType(0xc44fd609c0, 0x3a, 0x65aeb4, 0x9, 0x0, 0x74c960, 0xc4fcd000f0)
//         /home/toni/Laboratorio/gocrawl/PageContent.go:66 +0x35e
// main.main.func1(0xc420022100, 0xc42000d040, 0xc4200102c0, 0xc429982020, 0x2c)
//         /home/toni/Laboratorio/gocrawl/gocrawl.go:111 +0xd8
// created by main.main
//         /home/toni/Laboratorio/gocrawl/gocrawl.go:130 +0x2ec
