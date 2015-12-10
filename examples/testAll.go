package main

import (
    "github.com/stuphi/scroll-phat-go/scrollphat"
    "time"
    "os"
    "os/signal"
    "syscall"
)

func zero(s scrollphat.ScrollPhat) {
    s.Buffer = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    s.Offset = 0
    s.Update()
}

func main() {

    var sf scrollphat.ScrollPhat
    var x, y uint

    sf.Init()
    zero(sf)

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, syscall.SIGTERM)
    go func() {
        <-c
        zero(sf)
        os.Exit(1)
    }()

    for {
        for x = 0; x < 11; x++ { 
            for y = 0; y < 5; y++ {
                sf.SetPixel(x,y,1)
                sf.Update()
                time.Sleep(30 * time.Millisecond)
            }
        }
        for x = 0; x < 11; x++ { 
            for y = 0; y < 5; y++ {
                sf.SetPixel(x,y,0)
                sf.Update()
                time.Sleep(30 * time.Millisecond)
            }
        }
    }
}

