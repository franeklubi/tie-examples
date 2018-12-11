package main

import (
    . "github.com/franeklubi/tie"
)

func main() {
    Init(600, 600, "genimg", false)

    PassFunctions(
        preload,
        setup,
        draw,
        keyPressed,
    )

    Launch()
}

func preload() {

}

var (
    doppelganger Image = Image{[]byte{}, 2, 2}
    original Image = LoadImage("./4.png")
)

func setup() {
    doppelganger.PushPixel(Color{0, 255, 0, 255})
    doppelganger.PushPixel(Color{255, 0, 0, 255})
    doppelganger.PushPixel(Color{255, 255, 0, 255})
    doppelganger.PushPixel(Color{0, 0, 255, 255})
}

func draw() {
    Background(234, 123, 54, 255)

    Translate(MouseX, MouseY, 0)
        PastePixels(doppelganger, 0, 0, 100, 100)
        PastePixels(original, -100, 0, 100, 100)

    Translate(0, -30, 0)
        Text("generated", 15, false)
    Translate(-100, 0, 0)
        Text("loaded", 15, false)
}

func keyPressed() {

}
