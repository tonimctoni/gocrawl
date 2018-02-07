package main

import "testing"

// type UrlFinder struct{
//     re_urlfinder regexp.Regexp
// }

// func NewUrlFinder() *UrlFinder{
//     return &UrlFinder{*regexp.MustCompile("(?:href=|src=|url=)[\"']?([^\"' <>]*)")}
// }


// const MAX_HOST_SHARING_URLS_PER_SITE = 3


// func (u *UrlFinder) get_urls(base_url_s string, content []byte) []string{
//     urls:=make([]string, 0, 64)
//     hosts_counts:=make(map[string]uint)
//     base_url, err:=url.Parse(base_url_s)
//     if err!=nil{
//         return urls
//     }
//     for _,potential_url:=range u.re_urlfinder.FindAllSubmatch(content,-1){
//         potential_url, err:=url.Parse(string(potential_url[1]))
//         if err!=nil{
//             continue
//         }

//         potential_url=base_url.ResolveReference(potential_url)
//         hostname:=potential_url.Hostname()

//         if len(hostname)==0 || hosts_counts[hostname]>=MAX_HOST_SHARING_URLS_PER_SITE{
//             continue
//         }

//         hosts_counts[hostname]+=1

//         urls=append(urls, potential_url.String())
//     }
//     return urls
// }




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
