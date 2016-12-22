package main

import "errors"
import "sync"

type ThreadSafeStringQueue struct{
    mutex sync.Mutex
    data []string
    start int
    end int
    length int
}

func NewThreadSafeStringQueue(capacity uint) *ThreadSafeStringQueue{
    return &ThreadSafeStringQueue{sync.Mutex{},make([]string,capacity,capacity),0,0,0}
}

func (t *ThreadSafeStringQueue) push(str string) error{
    t.mutex.Lock()
    defer t.mutex.Unlock()
    if t.length==len(t.data){
        return errors.New("Cannot push onto full queue")
    }

    t.data[t.start]=str
    t.start=(t.start+1)%len(t.data)
    t.length++
    return nil
}

func (t *ThreadSafeStringQueue) push_slice(strs []string) error{
    t.mutex.Lock()
    defer t.mutex.Unlock()
    if t.length==len(t.data){
        return errors.New("Cannot push onto full queue")
    }

    if t.length+len(strs)>len(t.data){
        return errors.New("Not enougth space left in queue for slice push")
    }

    for _,str:=range strs{
        t.data[t.start]=str
        t.start=(t.start+1)%len(t.data)
    }
    t.length+=len(strs)

    return nil
}

func (t *ThreadSafeStringQueue) pull() (string,error){
    t.mutex.Lock()
    defer t.mutex.Unlock()
    if t.length==0{
        return "",errors.New("Cannot pull from empty queue")
    }

    return_string:=t.data[t.end]
    t.end=(t.end+1)%len(t.data)
    t.length--
    return return_string, nil
}

// func (t *ThreadSafeStringQueue) size() int{
//     t.mutex.Lock()
//     defer t.mutex.Unlock()

//     return t.length
// }

// func (t *ThreadSafeStringQueue) report() ([]string,int,int,int){
//     return t.data, t.start, t.end, t.length
// }