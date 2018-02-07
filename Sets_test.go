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

    if set.add("it") || set.add("is") || set.add("me"){
        t.Fail()
    }

    slice:=[]string{"hello", "1", "2", "3", "it", "it", "is", "me", "4"}
    slice=set.get_slice_without_contained(slice)
    if len(slice)!=4 || slice[0]!="1" || slice[1]!="2" || slice[2]!="3" || slice[3]!="4"{
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
