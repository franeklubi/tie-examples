package main

import (
    . "github.com/franeklubi/tie"
)

func main() {
    Init(800, 800, "textTest", false)

    PassFunctions(
        preload,
        setup,
        draw,
    )

    Launch()
}

func preload() {

}

func setup() {

}

func draw() {
    Background(174, 222, 16, 255)

    txt := "! yes ~ nyaaa ~"

    var txt_len float64 = float64(len(txt))
    txt_size :=  (Width-MouseX)/txt_len

    Translate(MouseX, MouseY, 0)

    Fill(145, 64, 130, 255)
    Text(txt, int(txt_size), false)
}
