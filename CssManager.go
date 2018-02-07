package main

import "regexp"
import "sync/atomic"
import "io/ioutil"
import "hash/fnv"
import "strings"
import "fmt"



type CssManager struct{
    counter uint64
    csshash_set Uint64Set
    re_comments *regexp.Regexp
    re_breaklines *regexp.Regexp
}

func NewCssManager() *CssManager {
    return &CssManager{0, *NewUint64Set(), regexp.MustCompile(`/\*(.|\n)*?\*/`), regexp.MustCompile(`\n{3,}`)}
}

const allowed_chars = "abcdefghijklmnopqrstuvwxzyABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789\n\t\r \"'(){}[]+-*/.,:;_@#%$!?=\\<>~^|&`"
func (c *CssManager) manage(content []byte) error{
    // Make sure all characters are allowed ones.
    outer: for _,char:=range content{
        for _,allowed_char:=range []byte(allowed_chars){
            if char==allowed_char{
                continue outer
            }
        }
        return nil
    }
    // Make some transformations for nicer text.
    content=c.re_comments.ReplaceAll(content, []byte(""))
    content=c.re_breaklines.ReplaceAll(content, []byte("\n\n"))
    content_as_string:=string(content)
    content_as_string=strings.ToLower(content_as_string)
    content_as_string=strings.TrimSpace(content_as_string)

    // Make sure it is a nice file.
    if strings.Count(content_as_string, "\n")<5 || len(content_as_string)<=50{
        return nil
    }
    content=[]byte(content_as_string)

    // Make sure the file is not repeated.
    h:=fnv.New64a()
    h.Write(content)
    content_hash:=h.Sum64()
    if c.csshash_set.add(content_hash){
        return nil
    }

    new_counter_value:=atomic.AddUint64(&c.counter, 1)
    return ioutil.WriteFile(fmt.Sprintf("css/css%06d.css", new_counter_value), content, 0644)
}