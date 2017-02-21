package main

import (
    "fmt"
    "flag"
    "os"
    ol "github.com/ossrs/go-oryx-lib/logger"
)

const (
    version string = "0.0.1"
)

func main() {
    fmt.Println(fmt.Sprintf("flv parser:%v, by panda of bravovcloud.com", version))

    var flvUrl string
    flag.StringVar(&flvUrl, "url", "./test.flv", "flv file to be parsed")

    ol.T(nil, "the input flv url is:", flvUrl)

    var f * os.File
    var err error
    if f, err = os.Open(flvUrl); err != nil {
        ol.T(nil, fmt.Sprint("open file:%v failed, err is %v", flvUrl, err))
        return
    }

    dec := &FlvDecoder{}
    if err = dec.ReadHeader(f); err != nil {
        ol.E(nil, fmt.Sprintf("decode flv header failed, err is %v", err))
        return
    }

    for {
        tag := &Tag{}
        if err = tag.Decode(f); err != nil {
            ol.E(nil, fmt.Sprintf("decode flv tag failed, err is %v", err))
            return
        }
    }
}
