package main

import "io/ioutil"
import "net/http"
import "regexp"
import "strings"
import "errors"

var re_urlfinder *regexp.Regexp
func initialize_global_regexp() {
    re_urlfinder=regexp.MustCompile("(?:href=|src=|url=)[\"']?([^\"' <>]*)")
}

func get_base_url(url string) string {
    search_start_index:=0
    if strings.HasPrefix(url, "http://"){
        search_start_index=len("http://")
    } else if strings.HasPrefix(url, "https://") {
        search_start_index=len("https://")
    } else {
        panic("Invalid URL")
    }

    end_index:=strings.IndexAny(url[search_start_index:], "/?#")
    if end_index==-1{
        return url
    }
    return url[:search_start_index+end_index]
}

type PageContent struct{
    base_url string
    full_url string
    content []byte
}

func NewPageContent(url string) (*PageContent, error) {
    response, err := http.Get(url)
    if err!=nil {
        return nil, err
    }
    defer response.Body.Close()

    content, err:=ioutil.ReadAll(response.Body)
    if err!=nil {
        return nil, err
    }

    return &PageContent{get_base_url(url), url, content}, nil
}

func NewPageContentIfContentType(url string, content_type string) (*PageContent, error) {
    response, err := http.Get(url)
    if err!=nil {
        return nil, err
    }
    defer response.Body.Close()

    if !strings.Contains(response.Header.Get("Content-Type"), content_type){
        return nil, errors.New("Content-Type does not contain text/html")
    }

    content, err:=ioutil.ReadAll(response.Body)
    if err!=nil {
        return nil, err
    }

    return &PageContent{get_base_url(url), url, content}, nil
}

func (p *PageContent) get_urls() []string{
    urls := make([]string, 0, 64)
    urls=append(urls, p.base_url)
    for _,value:=range re_urlfinder.FindAllSubmatch(p.content,-1){
        value:=string(value[1])
        hash_index:=strings.Index(value, "#")
        if hash_index!=-1{
            value=value[:hash_index]
        }
        if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://"){
            urls=append(urls, value)
        } else if value=="" || value=="/" || value=="\\" || value=="#"{
            continue
        } else if strings.HasPrefix(value, "//"){
            urls=append(urls, "http:"+value)
        } else if strings.HasPrefix(value, "/"){
            urls=append(urls, p.base_url+value)
        } else if strings.HasPrefix(value, "."){
            continue //too complicated, too much laziness
        } else if strings.HasSuffix(p.full_url, "/") {
            urls=append(urls, p.full_url+value) //full_url could also be done better
        } else {
            // fmt.Println("Having problems with:", p.full_url, value)
        }
    }
    return urls
}

func (p *PageContent) get_bytes() []byte{
    return p.content
}