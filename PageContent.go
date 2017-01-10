package main

import "io/ioutil"
import "net/http"
import "regexp"
import "strings"
import "errors"
import "time"

var re_urlfinder *regexp.Regexp
func initialize_global_regexp() {
    re_urlfinder=regexp.MustCompile("(?:href=|src=|url=)[\"']?([^\"' <>]*)")
}

var client_pagecontent *http.Client
func initialize_global_client() {
    client_pagecontent=&http.Client{Timeout: 5*time.Second}
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
    full_url string
    content []byte
}


// type WrongMimeTypeError error
func NewPageContentIfContentType(url string, content_type... string) (*PageContent, int, error) {
    response, err := client_pagecontent.Get(url)
    if err!=nil {
        return nil, -1, err
    }
    defer response.Body.Close()

    content_type_index:=int(-1)
    if func()bool{
        for index,ct:=range content_type{
            if strings.Contains(response.Header.Get("Content-Type"), ct){
                content_type_index=index
                return false
            }
        }
        return true
    }(){
        return nil, -1, errors.New("Content-Type does not contain mime type")
    }

    content, err:=ioutil.ReadAll(response.Body)
    if err!=nil {
        return nil, -1, err
    }

    return &PageContent{url, content}, content_type_index, nil
}

func (p *PageContent) get_urls() []string{
    urls := make([]string, 0, 64)
    base_url:=get_base_url(p.full_url)
    urls=append(urls, base_url)
    for _,value:=range re_urlfinder.FindAllSubmatch(p.content,-1){
        value:=string(value[1])
        hash_index:=strings.Index(value, "#")
        if hash_index!=-1{
            value=value[:hash_index]
        }
        quest_index:=strings.Index(value, "?")//The web is too large to crawl it all anyways
        if quest_index!=-1{
            value=value[:quest_index]
        }
        if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://"){
            urls=append(urls, value)
        } else if value=="" || value=="/" || value=="\\" || value=="#"{
            continue
        } else if strings.HasPrefix(value, "//"){
            urls=append(urls, "http:"+value)
        } else if strings.HasPrefix(value, "/"){
            urls=append(urls, base_url+value)
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