package main

type EnumPreProcessing int
func (v EnumPreProcessing) String() (s string) {
    switch v {
    case 0:
        s = "NonPreProcessing"
    case 1:
        s = "PreProcessing"
    }
    return
}

type EnumTagType int
func (v EnumTagType) String() (s string) {
    switch v {
    case 8:
        s = "audio"
    case 9:
        s = "video"
    case 18:
        s = "script"
    }
    return
}

type EnumSoundFormat int
func (v EnumSoundFormat) String() (s string) {
    switch v {
    case 1:
        s = "ADPCM"
    case 2:
        s = "MP3"
    case 10:
        s = "AAC"
    case 11:
        s = "Speex"
    }
    return
}

type EnumSoundRate int
func (v EnumSoundRate) String() (s string) {
    switch v {
    case 0:
        s = "5.5kHz"
    case 1:
        s = "11kHz"
    case 2:
        s = "22kHz"
    case 3:
        s = "44kHz"
    }
    return
}

type EnumSoundSize int
func (v EnumSoundSize) String() (s string) {
    switch v {
    case 0:
        s = "8bit"
    case 1:
        s = "16bit"
    }
    return
}

type EnumSoundType int
func (v EnumSoundType) String() (s string) {
    switch v {
    case 0:
        s = "Mono"
    case 1:
        s = "Stereo"
    }
    return
}

// if sound format = 10
type EnumAACPacketType int
func (v EnumAACPacketType) String() (s string) {
    switch v {
    case 0:
        s = "Sequence header"
    case 1:
        s = "raw"
    }
    return
}

type EnumFrameType int
func (v EnumFrameType) String() (s string) {
    switch v {
    case 1:
        s = "key frame"
    case 2:
        s = "inter frame"
    case 3:
        s = "disposable inter frame"
    case 4:
        s = "generated key frame"
    case 5:
        s = "video info/command frame"
    }
    return
}

type EnumCodecId int
func (v EnumCodecId) String() (s string) {
    switch v {
    case 2:
        s = "H263"
    case 3:
        s = "Screen video"
    case 4:
        s = "VP6"
    case 5:
        s = "VP6 with alpha channel"
    case 6:
        s = "Screen video 2"
    case 7:
        s = "AVC"
    }
    return
}

type EnumAVCPacketType int
func (v EnumAVCPacketType) String() (s string) {
    switch v {
    case 0:
        s = "Sequence header"
    case 1:
        s = "NALU"
    case 2:
        s = "End of sequence"
    }
    return
}




