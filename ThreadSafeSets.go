package main

import "sync"

// func consumer2() {
//     //Fuck go's restrictions
//     fmt.Println("")
//     strings.HasPrefix("", "")
//     http.Get("")
//     ioutil.ReadAll(nil)
//     errors.New("")
//     regexp.MustCompile("")
// }

//If I make a real set: use RWMutex.
type ThreadSafeStringSet struct{
    mutex sync.Mutex
    set map[string]bool
}

func NewThreadSafeStringSet() *ThreadSafeStringSet{
    return &ThreadSafeStringSet{sync.Mutex{}, make(map[string]bool)}
}

func (t *ThreadSafeStringSet) add(key string) bool {
    t.mutex.Lock()
    defer t.mutex.Unlock()

    contains:=t.set[key]
    t.set[key]=true

    return contains
}

//If I make a real set: use RWMutex.
type ThreadSafeUint32Set struct{
    mutex sync.Mutex
    set map[uint32]bool
}

func NewThreadSafeUint32Set() *ThreadSafeUint32Set{
    return &ThreadSafeUint32Set{sync.Mutex{}, make(map[uint32]bool)}
}

func (t *ThreadSafeUint32Set) add(key uint32) bool {
    t.mutex.Lock()
    defer t.mutex.Unlock()

    contains:=t.set[key]
    t.set[key]=true

    return contains
}