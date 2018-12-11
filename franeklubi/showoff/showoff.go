package main

import (
    . "github.com/franeklubi/tie"
)

const (
    pass_width  = 500
    pass_height = 500
)

func main() {

    Init(pass_width, pass_height, "showoff", false)

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
    // Background(255, 163, 249, 255)
    StrokeWidth(10)
}

func draw() {

    Translate(-Width/2, -Height/2, 0)
    Translate(MouseX, MouseY, 0)

    Background(255, 163, 249, 255)

    Fill(255, 0, 0, 255)
    Ellipse(Width/2, Height/2, Width/2, Height/2)

    StrokeWidth(10)
    x := Cos(DegToRad(float64(Frames))) * Width/4
    y := Sin(DegToRad(float64(Frames))) * Height/4
    Line(Width/2, Height/2, x+Width/2, y+Height/2)
    StrokeWidth(1)

    Fill(29, 189, 45, 255)
    Ellipse(Mod(Frames, Width), Height/2, 300, Sin(float64(Frames)/100)*Height/2)

    Fill(131, 88, 255, 100)
    Rect(Mod(Frames*2, Width), Height/2, Width/4, Sin(float64(Frames)/100)*Height/2)

    // RotateZ(70)
    // RotateY(70)
    // RotateX(70)
    Fill(255, 255, 255, 255)
    Rect(Width/2+Width/4, Height/2, 100, 100)
    // RotateX(-70)
    // RotateY(-70)
    // RotateZ(-70)

    Frames++
}
