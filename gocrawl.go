package main

import "runtime"
import "time"
import "fmt"
import "sync/atomic"
import "io/ioutil"

const NUM_THREADS = 500


func main() {
    url_reservoir:=NewStringReservoir(1024*1024)
    url_reservoir.add("http://cssdb.co")
    visited_urls:=NewStringSet()
    client:=NewClient()
    url_finder:=NewUrlFinder()
    css_manager:=NewCssManager()

    error_counter:=uint32(0)
    counter:=uint32(0)
    for i:=0;i<NUM_THREADS;i++{
        go func(){
            defer func(){
                if r:=recover(); r!=nil {
                    new_error_counter:=atomic.AddUint32(&error_counter, 1)
                    ioutil.WriteFile(fmt.Sprintf("error_%03d.txt", new_error_counter), []byte(fmt.Sprintln(r)), 0644)
                }
            }()

            for{
                url, err:=url_reservoir.get()
                if err!=nil{
                    runtime.Gosched()
                    continue
                }

                if visited_urls.add(url){
                    continue
                }

                page_content,content_type_index,err:=client.get_if_content_type_is(url, "text/html", "text/css")
                if err!=nil{
                    continue
                }

                switch content_type_index{
                case 0:
                    found_urls:=url_finder.get_urls(url, page_content)
                    if len(found_urls)!=0{
                        found_urls=visited_urls.get_slice_without_contained(found_urls)
                        url_reservoir.add_slice(found_urls)
                    }
                    c:=atomic.AddUint32(&counter, 1)
                    fmt.Println(c, url)
                case 1:
                    err=css_manager.manage(page_content)
                    if err!=nil{
                        fmt.Println(err.Error())
                    }
                default:
                    panic("content_type_index has a value it should never have")
                }

                // time.Sleep(2*time.Second)
            }
        }()
    }
    for{
        time.Sleep(10*time.Hour)
    }
    // time.Sleep(4*time.Minute)
}
