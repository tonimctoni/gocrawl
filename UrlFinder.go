package main

import "regexp"
import "net/url"

type UrlFinder struct{
    re_urlfinder regexp.Regexp
}

func NewUrlFinder() *UrlFinder{
    return &UrlFinder{*regexp.MustCompile("(?:href=|src=|url=)[\"']?([^\"' <>]*)")}
}


const MAX_HOST_SHARING_URLS_PER_SITE = 3


func (u *UrlFinder) get_urls(base_url_s string, content []byte) []string{
    urls:=make([]string, 0, 64)
    hosts_counts:=make(map[string]uint)
    base_url, err:=url.Parse(base_url_s)
    if err!=nil{
        return urls
    }
    for _,potential_url:=range u.re_urlfinder.FindAllSubmatch(content,-1){
        potential_url, err:=url.Parse(string(potential_url[1]))
        if err!=nil{
            continue
        }

        potential_url=base_url.ResolveReference(potential_url)
        hostname:=potential_url.Hostname()

        if len(hostname)==0 || hosts_counts[hostname]>=MAX_HOST_SHARING_URLS_PER_SITE{
            continue
        }

        hosts_counts[hostname]+=1

        urls=append(urls, potential_url.String())
    }
    return urls
}