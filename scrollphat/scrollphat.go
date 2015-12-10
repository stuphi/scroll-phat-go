package scrollphat

import (
    "fmt"
//    "errors"
    "github.com/mrmorphic/hwio"
    "strings"
)

const I2C_ADDR = 0x60
const CMD_SET_MODE = 0x00
const CMD_SET_BRIGHTNESS = 0x19
const MODE_5X11 = 3  //0b00000011



type ScrollPhat struct {
    Buffer []byte
    device hwio.I2CDevice
    Offset int
}

func (s *ScrollPhat) Init() {

    m, e := hwio.GetModule("i2c")
    if e != nil {
        fmt.Printf("could not get i2c module: %s\n", e)
        return
    }
    i2c := m.(hwio.I2CModule)

    // Uncomment on Raspberry pi, which doesn't automatically enable i2c bus. BeagleBone does,
    // as the default device tree enables it.

    i2c.Enable()
    //defer i2c.Disable()

    s.device = i2c.GetDevice(I2C_ADDR)

    e = s.device.WriteByte(CMD_SET_MODE, MODE_5X11)
    if e != nil {
        fmt.Printf("could not write Byte: %s\n", e)
    }

    e = s.device.WriteByte(CMD_SET_BRIGHTNESS, 2)
    if e != nil {
        fmt.Printf("could not set brightness: %s\n", e)
    }
    s.Buffer = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    s.Offset = 0
}

func (s *ScrollPhat) Update() {
    var window []byte
    if s.Offset > len(s.Buffer){
        s.Offset = 0
    }
    if s.Offset + 11 <= len(s.Buffer) {
        window = s.Buffer[s.Offset:s.Offset + 11]
    } else {
        window = s.Buffer[s.Offset:]
        window = append(window, s.Buffer[:11 - len(window)]...)
    }
    //window = append(window, 0xFF)

    e := s.device.Write(0x01, window)
    if e != nil {
        fmt.Printf("could not write Byte Array: %s\n", e)
    }
    
    e = s.device.WriteByte(0x0C, 0)
    if e != nil {
        fmt.Printf("could not update column register: %s\n", e)
    }
}

func (s *ScrollPhat) SetPixel(x, y, val uint) {
    if val == 1 {
        s.Buffer[x] |= (1 << y)
    } else {
        s.Buffer[x] &= ^(1 << y)
    }
}

func StringToBuf(s string) ([]byte){
    
    var ret []byte
    var c byte

    s = strings.ToUpper(s)

    letters := map[string][]byte{
        "A": {0x1E, 0x09, 0x1E, 0x00},
        "B": {0x1F, 0x15, 0x0A, 0x00},
        "C": {0x0E, 0x11, 0x0A, 0x00},
        "D": {0x1F, 0x11, 0x0E, 0x00},
        "E": {0x1F, 0x15, 0x15, 0x00},
        "F": {0x1F, 0x05, 0x05, 0x00},
        "G": {0x0E, 0x11, 0x1D, 0x00},
        "H": {0x1F, 0x04, 0x1F, 0x00},
        "I": {0x11, 0x1F, 0x11, 0x00},
        "J": {0x09, 0x11, 0x0F, 0x00},
        "K": {0x1F, 0x04, 0x1B, 0x00},
        "L": {0x1F, 0x10, 0x10, 0x00},
        "M": {0x1F, 0x02, 0x04, 0x02, 0x1F, 0x00},
        "N": {0x1F, 0x02, 0x0C, 0x1F, 0x00},
        "O": {0x0E, 0x11, 0x0E, 0x00},
        "P": {0x1F, 0x09, 0x06, 0x00},
        "Q": {0x0E, 0x11, 0x09, 0x16, 0x00},
        "R": {0x1F, 0x09, 0x16, 0x00},
        "S": {0x12, 0x15, 0x09, 0x00},
        "T": {0x01, 0x1F, 0x01, 0x00},
        "U": {0x0F, 0x10, 0x10, 0x0F, 0x00},
        "V": {0x0F, 0x10, 0x0F, 0x00},
        "W": {0x0F, 0x10, 0x08, 0x10, 0x0F, 0x00},
        "X": {0x1B, 0x04, 0x1B, 0x00},
        "Y": {0x03, 0x1C, 0x03, 0x00},
        "Z": {0x19, 0x15, 0x13, 0x00},
        "0": {0x0E, 0x15, 0x0E, 0x00},
        "1": {0x12, 0x1F, 0x10, 0x00},
        "2": {0x19, 0x15, 0x12, 0x00},
        "3": {0x11, 0x15, 0x0A, 0x00},
        "4": {0x0E, 0x09, 0x1C, 0x00},
        "5": {0x17, 0x15, 0x09, 0x00},
        "6": {0x0E, 0x15, 0x08, 0x00},
        "7": {0x19, 0x05, 0x03, 0x00},
        "8": {0x0A, 0x15, 0x0A, 0x00},
        "9": {0x02, 0x15, 0x0E, 0x00},
        " ": {0x00, 0x00},
        ":": {0x0A, 0x00},
        ".": {0x10, 0x00},
        "-": {0x04, 0x04, 0x04, 0x00},

    }
    for i:=0; i<len(s); i++ {
        c = s[i]
        ret = append(ret, letters[string([]byte{c})]...)
    }

    return ret
}
