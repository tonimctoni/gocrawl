package main

import "testing"


func TestStringSet(t *testing.T) {
    set:=NewStringSet()
    if set.contains("hello") || set.contains("") || set.contains("_"){
        t.Fail()
    }

    if set.add("hello"){
        t.Fail()
    }
    if !set.add("hello"){
        t.Fail()
    }
    if !set.contains("hello"){
        t.Fail()
    }
}


func TestNewUint64Set(t *testing.T) {
    set:=NewUint64Set()
    if set.add(14580346) || set.add(0) || set.add(8){
        t.Fail()
    }

    if set.add(73489){
        t.Fail()
    }
    if !set.add(73489){
        t.Fail()
    }
}
