package main

import "sync"






type StringSet struct{
    rw_mutex sync.RWMutex
    set map[string]bool
}

func NewStringSet() *StringSet{
    return &StringSet{sync.RWMutex{}, make(map[string]bool)}
}

func (t *StringSet) add(key string) bool {
    t.rw_mutex.Lock()
    defer t.rw_mutex.Unlock()

    contains:=t.set[key]
    t.set[key]=true

    return contains
}

func (t *StringSet) contains(key string) bool {
    t.rw_mutex.RLock()
    defer t.rw_mutex.RUnlock()

    return t.set[key]
}

func (t *StringSet) get_slice_without_contained(string_slice []string) []string{
    ret:=make([]string,0,len(string_slice))
    t.rw_mutex.RLock()
    defer t.rw_mutex.RUnlock()

    for _,s:=range string_slice{
        if !t.set[s]{
            ret=append(ret, s)
        }
    }

    return ret
}






type Uint64Set struct{
    mutex sync.Mutex
    set map[uint64]bool
}

func NewUint64Set() *Uint64Set{
    return &Uint64Set{sync.Mutex{}, make(map[uint64]bool)}
}

func (u *Uint64Set) add(key uint64) bool {
    u.mutex.Lock()
    defer u.mutex.Unlock()

    contains:=u.set[key]
    u.set[key]=true

    return contains
}