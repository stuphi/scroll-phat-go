package main

import (
    "fmt"
//    "errors"
    "github.com/mrmorphic/hwio"
    "time"
)

const I2C_ADDR = 0x60
const CMD_SET_MODE = 0x00
const CMD_SET_BRIGHTNESS = 0x19
const MODE_5X11 = 3  //0b00000011

type ScrollPhat struct {
    buffer []byte
    device hwio.I2CDevice
    offset int
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

    s.offset = 0
}

func (s *ScrollPhat) Update() {
    var window []byte
    if s.offset > len(s.buffer){
        s.offset = 0
    }
    if s.offset + 11 <= len(s.buffer) {
        window = s.buffer[s.offset:s.offset + 11]
    } else {
        window = s.buffer[s.offset:]
        window = append(window, s.buffer[:11 - len(window)]...)
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

func main() {

    var sf ScrollPhat

    sf.Init()

    //sf.buffer = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10 }
    sf.buffer = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                       0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
                       0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
                       0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
                   }
    for i:=0;i<=100;i++ {
        sf.offset = sf.offset + 1
        sf.Update()
        time.Sleep(500 * time.Millisecond)
    }

    sf.buffer = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }
    sf.offset = 0
    sf.Update()
}

