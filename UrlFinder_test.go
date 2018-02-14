package main

import "testing"

func TestUrlFinder(t *testing.T) {
    test:=func(html string, urls []string){
        url_finder:=NewUrlFinder()
        found_urls:=url_finder.get_urls("http://base.com/a/b/c/d", []byte(html))
        if found_urls==nil || len(urls)!=len(found_urls){
            t.Fail()
        }
        for i:=0;i<len(urls);i++{
            if urls[i]!=found_urls[i]{
                t.Fail()
            }
        }
    }

    test(`
        <link rel="stylesheet" href="https://www.google.com"> 
        <link rel="stylesheet" href="http://www.google.com/asd/qwe"> 
        <link rel="stylesheet" href="hello"> 
        <link rel="stylesheet" href="/hello"> 
        <link rel="stylesheet" href="../../hello"> 
        <a href="/hello">hello</a> <a href="/hello2">hello_</a> <a href="/hello_">hello3</a>
        <a href=/hello4>hello</a> 
        <a href=/hello5 style="">hello</a> 
        `, []string{
        "https://www.google.com",
        "http://www.google.com/asd/qwe",
        "http://base.com/a/b/c/hello",
        "http://base.com/hello",
        "http://base.com/a/hello",
    })

    test(`
        <a href="/hello">hello</a> <a href="/hello2">hello_</a> <a href="/hello_">hello3</a>
        `, []string{
        "http://base.com/hello",
        "http://base.com/hello2",
        "http://base.com/hello_",
    })

    test(`
        <a href=/hello4>hello</a> 
        <a href=/hello5 style="">hello</a> 
        `, []string{
        "http://base.com/hello4",
        "http://base.com/hello5",
    })

    test(`
        <a href=mailto:felix@mail.com>mailme</a> 
        `, []string{
    })
}
