package main

import (
    . "github.com/franeklubi/tie"
)

func main() {
    Init(600, 600, "replace", false)

    PassFunctions(
        preload,
        setup,
        draw,
    )

    Launch()
}

var (
    img Image
)

func preload() {
    img = LoadImage("./img.png")
}

func setup() {
    pixel_to_replace := Color{0, 255, 0, 255}
    img.Replace(1, 1, pixel_to_replace)
    Println(img.PixelAt(1, 1))
}

func draw() {
    Background(54, 200, 91, 255)

    Translate(MouseX, MouseY, 0)
    PastePixels(img, 0, 0, 200, 200)
    PastePixels(img.GetPixels(0, 0, 5, 5), 0, -200, 200, 200)
}
