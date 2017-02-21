package main

import (
    "encoding/binary"
)

func Bytes3ToUint32(b []byte) uint32 {
    nb := []byte{}
    nb = append(nb, 0)
    nb = append(nb, b...)
    return binary.BigEndian.Uint32(nb)
}

func GetTagTimestamp(ts []byte) uint32 {
    nb := []byte{}
    nb = append(nb, ts[3])
    nb = append(nb, ts[0:3]...)
    return binary.BigEndian.Uint32(nb)
}