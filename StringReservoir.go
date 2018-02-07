package main

import "sync"
import "time"
import "math/rand"
import "errors"

type StringReservoir struct{
    mutex sync.Mutex
    rng rand.Rand
    strings []string
    capacity int
}


func NewStringReservoir(capacity int) *StringReservoir {
    return &StringReservoir{sync.Mutex{}, *rand.New(rand.NewSource(time.Now().UnixNano())), make([]string, 0, capacity), capacity}
}

func (s *StringReservoir) add(new_string string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    if len(s.strings)<s.capacity{
        s.strings=append(s.strings, new_string)
    } else{
        s.strings[s.rng.Intn(len(s.strings))]=new_string
    }
}

func (s *StringReservoir) add_slice(new_strings []string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    for _, new_string := range new_strings{
        if len(s.strings)<s.capacity{
            s.strings=append(s.strings, new_string)
        } else{
            s.strings[s.rng.Intn(len(s.strings))]=new_string
        }
    }
}

func (s *StringReservoir) get() (string,error){
    s.mutex.Lock()
    defer s.mutex.Unlock()
    if len(s.strings)==0{
        return "", errors.New("Reservoir is empty")
    } else{
        pos:=s.rng.Intn(len(s.strings))
        ret:=s.strings[pos]
        s.strings[pos]=s.strings[len(s.strings)-1]
        s.strings[len(s.strings)-1]=""
        s.strings=s.strings[:len(s.strings)-1]
        return ret, nil
    }
}
