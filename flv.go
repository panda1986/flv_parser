package main

import (
    "io"
    "encoding/binary"
    ol "github.com/ossrs/go-oryx-lib/logger"
    "fmt"
    "bytes"
    "reflect"
)

type FlvDecoder struct {
    Version uint8
    ContainAudio bool
    ContainVideo bool
    DataOffset uint32
}

func (v *FlvDecoder) ReadHeader(r io.Reader) (err error) {
    tmp := [5]uint8{}
    if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
        ol.E(nil, fmt.Sprint("read flv header failed, err is %v", err))
        return
    }

    if string(tmp[0:3]) != "FLV" {
        ol.E(nil, fmt.Sprintf("exp flv header=FLV, actual=%v", string(tmp[0:3])))
    }

    v.Version = tmp[3]
    if tmp[4] & 0x04 == 0x04 {
        v.ContainAudio = true
    }
    if tmp[4] & 0x01 == 0x01 {
        v.ContainVideo = true
    }

    binary.Read(r, binary.BigEndian, &v.DataOffset)
    ol.T(nil, fmt.Sprintf("flv header:%+v", v))
    return
}

type Tag struct {
    PreviousTagSize uint32
    Filter EnumPreProcessing
    TagType EnumTagType
    Timestamp uint32
    DataSize uint32
}

func (v *Tag) Decode(r io.Reader) (err error) {
    binary.Read(r, binary.BigEndian, &v.PreviousTagSize)

    //Tag Header
    var tmp uint8
    if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
        ol.E(nil, fmt.Sprintf("read tag header failed, err is %v", err))
        return
    }
    if tmp & 0x20 == 0x20 {
        v.Filter = EnumPreProcessing(1)
    }
    v.TagType = EnumTagType(tmp & 0x1f)

    dataSize := make([]byte, 3)
    if _, err = io.ReadFull(r, dataSize); err != nil {
        ol.E(nil, fmt.Sprintf("read tag header data size failed, err is %v", err))
        return
    }
    v.DataSize = Bytes3ToUint32(dataSize)

    timeStamp := make([]byte, 4)
    if _, err = io.ReadFull(r, timeStamp); err != nil {
        ol.E(nil, fmt.Sprintf("read tag header timestamp failed, err is %v", err))
        return
    }
    v.Timestamp = GetTagTimestamp(timeStamp)

    streamId := make([]byte, 3)
    if _, err = io.ReadFull(r, streamId); err != nil {
        ol.E(nil, fmt.Sprintf("read tag header streamId failed, err is %v", err))
        return
    }

    // Tag Data
    data := make([]byte, v.DataSize)
    io.ReadFull(r, data)

    var td TadData
    switch v.TagType {
    case 8:
        td = &AudioTagData{}
    case 9:
        td = &VideoTagData{}
    case 18:
        td = &ScriptData{}
    }

    if td != nil {
        td.Decode(data)
    }

    ol.T(nil, fmt.Sprintf("decode a tag: %v, %v", v, td))
    return
}

func (v *Tag) String() string {
    return fmt.Sprintf("tag filter:%v, type:%v, data size:%v, timestamp:%v", v.Filter, v.TagType, v.DataSize, v.Timestamp)
}

type TadData interface {
    Decode(data []byte) (err error)
    String() string
}

type ScriptData struct {
}

func (v *ScriptData) readString(r io.Reader) (strLen uint16, strData []byte, err error) {
    if err = binary.Read(r, binary.BigEndian, &strLen); err != nil {
        return
    }
    strData = make([]byte, strLen)
    if err = binary.Read(r, binary.BigEndian, strData); err != nil {
        return
    }
    return
}

func (v *ScriptData) readValue(r io.Reader) (value interface{}, err error) {
    var valueType uint8
    if err = binary.Read(r, binary.BigEndian, &valueType); err != nil {
        return
    }
    switch valueType {
    case 0:
        var tmp float64
        if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
            return
        }
        value = tmp
    case 1:
        // Boolean, UI8
        var tmp uint8
        if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
            break
        }
        value = tmp
    case 2:
        var data []byte
        if _, data , err = v.readString(r); err != nil {
            return
        }
        value = string(data)
    }
    return
}

func (v *ScriptData) Decode(data []byte) (err error) {
    r := bytes.NewReader(data)
    for {
        var metaType uint8
        if err = binary.Read(r, binary.BigEndian, &metaType); err != nil {
            break
        }

        switch metaType {
        case 2:
            //String, SCRIPTDATASTRING, video_file_format_spec_v10_1.pdf, page 83
            var strLen uint16
            var strData []byte
            if strLen, strData, err = v.readString(r); err != nil {
                break
            }
            ol.T(nil, fmt.Sprintf("meta type String: len=%v, data=%v", strLen, string(strData)))
        case 7:
            // Reference,UI16
            var metaValue uint16
            if err = binary.Read(r, binary.BigEndian, &metaValue); err != nil {
                break
            }
            ol.T(nil, fmt.Sprintf("meta type Reference:%v", metaValue))
        case 8:
            // ECMA array, SCRIPTDATAECMAARRAY
            var arrLen uint32 // ECMAArrayLength
            if err = binary.Read(r, binary.BigEndian, &arrLen); err != nil {
                return
            }
            for i := 0; i < int(arrLen); i ++ { // SCRIPTDATAOBJECTPROPERTY [ ]
                var strData []byte
                if _, strData, err = v.readString(r); err != nil {
                    break
                }

                propertyName := string(strData)
                var propertyValue interface{}
                if propertyValue, err = v.readValue(r); err != nil {
                    break
                }
                ol.T(nil, fmt.Sprintf("read one property: %v-%v, %v", propertyName, propertyValue, reflect.TypeOf(propertyValue)))
            }
        default:
            ol.W(nil, fmt.Sprintf("unsupported type:%v", metaType))
        }
    }
    
    if err == io.EOF {
        return nil
    }
    return
}

func (v *ScriptData) String() string {
    return ""
}

type AudioTagData struct {
    SoundFormat EnumSoundFormat
    SoundRate EnumSoundRate
    SoundSize EnumSoundSize
    SoundType EnumSoundType
    AACPacketType EnumAACPacketType
}

func (v *AudioTagData) Decode(data []byte) (err error) {
    tmp := data[0]

    v.SoundFormat = EnumSoundFormat((tmp & 0xf0) >> 4)
    v.SoundRate = EnumSoundRate((tmp & 0x0C) >> 2)
    v.SoundSize = EnumSoundSize((tmp & 0x02) >> 1)
    v.SoundType = EnumSoundType(tmp & 0x01)

    if v.SoundFormat == EnumSoundFormat(10) {
        v.AACPacketType = EnumAACPacketType(data[1])
    }
    return
}

func (v *AudioTagData) String() string {
    return fmt.Sprintf("audio format:%v, rate:%v, size:%v, type:%v, aac:%v", v.SoundFormat, v.SoundRate, v.SoundSize, v.SoundType, v.AACPacketType)
}

type VideoTagData struct {
    FrameType EnumFrameType
    CodecId EnumCodecId
    AVCPacketType EnumAVCPacketType
    CompositionTime int32
    NALULen uint32
}

func (v *VideoTagData) Decode(data []byte) (err error) {
    tmp := data[0]
    v.FrameType = EnumFrameType(tmp & 0xf0 >> 4)
    v.CodecId = EnumCodecId(tmp & 0x0f)
    if v.CodecId == EnumCodecId(7) {
        v.AVCPacketType = EnumAVCPacketType(data[1])
        v.CompositionTime = int32(Bytes3ToUint32(data[2:5])) // cts, if profile == baseline, cts always = 0, else not.
        if len(data) >= 6 {
            v.NALULen = binary.BigEndian.Uint32(data[5:9])
        }
    }
    return
}

func (v *VideoTagData) String() string {
    return fmt.Sprintf("video frame type:%v, codec:%v, avc packet type:%v, composition time:%v, NALU len:%v", v.FrameType, v.CodecId, v.AVCPacketType, v.CompositionTime, v.NALULen)
}