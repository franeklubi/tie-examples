package main

import (
    . "github.com/franeklubi/tie"
)

func main() {
    Init(500, 500, "2dtest", false)

    PassFunctions(
        preload,
        setup,
        draw,
    )

    Launch()
}

func preload(){

}

var (
    zdj Image
    zdjwyc Image
)

func setup() {
    zdj = LoadImage("./image.png")
    // Println(len(zdj.Pixels))
    Println(zdj.PixelAt(9, 0))
    zdjwyc = zdj.GetPixels(100, 100, 200, 200)
}

func draw() {

    Background(85, 207, 153, 255)

    Fill(109, 119, 170, 255)
    Rect(Width/2, Height/2, 100, 100)

    Fill(177, 179, 51, 100)
    Rect(Width/2+50, Height/2+50, 100, 100)

    Fill(255, 255, 255, 100)
    Translate(MouseX, MouseY, 0)
    PastePixels(zdj, 0, 0, Width/4, Height/4)

    PastePixels(zdjwyc, 0, -Height/4, Width/4, Height/4)

    // Background(89, 4, 85, 255)
    Fill(80, 74, 210, 255)
    buffer := CopyPixels()
    PastePixels(buffer, -Width/4, 0, Width/4, Height/4)
}
