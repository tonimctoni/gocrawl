package main

import "testing"

func TestStringReservoir(t *testing.T) {
    string_reservoir:=NewStringReservoir(10)
    if len(string_reservoir.strings)!=0 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }
    string_reservoir.add("hello")
    if len(string_reservoir.strings)!=1 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }
    s,e:=string_reservoir.get()
    if s!="hello" || e!=nil || len(string_reservoir.strings)!=0 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }
    s,e=string_reservoir.get()
    if s!="" || e==nil || len(string_reservoir.strings)!=0 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }

    for i:=0;i<12;i++{
        string_reservoir.add("hello")
    }
    if len(string_reservoir.strings)!=10 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }

    for i:=0;i<10;i++{
        s,e:=string_reservoir.get()
        if s!="hello" || e!=nil || len(string_reservoir.strings)!=10-i-1 || cap(string_reservoir.strings)!=10{
            t.Fail()
        }
    }
    s,e=string_reservoir.get()
    if s!="" || e==nil || len(string_reservoir.strings)!=0 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }

    string_reservoir.add_slice([]string{"hello","hello","hello"})
    if len(string_reservoir.strings)!=3 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }

    s,e=string_reservoir.get()
    if s!="hello" || e!=nil || len(string_reservoir.strings)!=2 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }

    string_reservoir.add_slice([]string{"hello","hello","hello"})
    if len(string_reservoir.strings)!=5 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }

    string_reservoir.add_slice([]string{"hello","hello","hello"})
    if len(string_reservoir.strings)!=8 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }

    string_reservoir.add_slice([]string{"hello","hello","hello"})
    if len(string_reservoir.strings)!=10 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }

    s,e=string_reservoir.get()
    if s!="hello" || e!=nil || len(string_reservoir.strings)!=9 || cap(string_reservoir.strings)!=10{
        t.Fail()
    }
}
