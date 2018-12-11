package main

import (
    . "github.com/franeklubi/tie"
)

func main() {
    Init(500, 500, "3dqueue", false)

    PassFunctions(
        preload,
        setup,
        draw,
        keyPressed,
    )

    Launch()
}

func preload(){

}

func setup() {
    // NoStroke()
}

var (
    fidelity uint32 = 10
)

func draw() {
    Background(231, 23, 74, 255)
    Translate(0, 0, -Width*0.8)

    Translate(0, Height/2, 0)
        RotateX(ReMap(MouseY+Height/2, 0, Height, 0, 360))
    Translate(0, -Height/2, 0)
    Translate(Width/2, 0, 0)
        RotateY(ReMap(MouseX+Width/2, 0, Width, 0, 360))
    Translate(-Width/2, 0, 0)

    Fill(130, 23, 74, 255)
    Ellipse(Width/2, Height/2, Width*2, Height*2)
    Fill(43, 254, 230, 255)
    Rect(0, 0, Width, Height)

    Fill(131, 191, 221, 255)
    // Translate(Width/2, Height/2, 0)
    // Cube(0.5, 0.5, 0.5)
    DepthRefreshOff()
        Sphere(1, fidelity)
    DepthRefreshOn()

    // Translate(200, 0, 0)
    // Cube(200)
    // Translate(-Width/2-200, -Height/2, 0)

    Point(Width/2, Height/2)
}

func keyPressed() {
    if ( Key == UP ) {
        fidelity++
        Println(fidelity)
    }
    if ( Key == DOWN ) {
        fidelity--
        Println(fidelity)
    }
}
