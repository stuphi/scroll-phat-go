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


func main(){
    
    m, e := hwio.GetModule("i2c")
    if e != nil {
        fmt.Printf("could not get i2c module: %s\n", e)
        return
    }
    i2c := m.(hwio.I2CModule)

    // Uncomment on Raspberry pi, which doesn't automatically enable i2c bus. BeagleBone does,
    // as the default device tree enables it.

    i2c.Enable()
    defer i2c.Disable()

    device := i2c.GetDevice(I2C_ADDR)

    e = device.WriteByte(CMD_SET_MODE, MODE_5X11)
    if e != nil {
        fmt.Printf("could not write Byte: %s\n", e)
    }

    e = device.WriteByte(CMD_SET_BRIGHTNESS, 2)
    if e != nil {
        fmt.Printf("could not set brightness: %s\n", e)
    }

    myData := []byte{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF }

    e = device.Write(0x01, myData)
    if e != nil {
        fmt.Printf("could not write Byte Array: %s\n", e)
    }

    e = device.WriteByte(0x0C, 0)
    if e != nil {
        fmt.Printf("could not update column register: %s\n", e)
    }

    time.Sleep(500 * time.Millisecond)

    myData = []byte{0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00}

    e = device.Write(0x01, myData)
    if e != nil {
        fmt.Printf("could not write Byte Array: %s\n", e)
    }

    e = device.WriteByte(0x0C, 0)
    if e != nil {
        fmt.Printf("could not update column register: %s\n", e)
    }

    time.Sleep(500 * time.Millisecond)

    myData = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

    e = device.Write(0x01, myData)
    if e != nil {
        fmt.Printf("could not write Byte Array: %s\n", e)
    }

    e = device.WriteByte(0x0C, 0)
    if e != nil {
        fmt.Printf("could not update column register: %s\n", e)
    }


}

