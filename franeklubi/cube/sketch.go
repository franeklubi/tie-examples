package main

import (
    . "github.com/franeklubi/tie"
)

const (
    pass_width  = 500
    pass_height = 500
)

func main() {
    Init(pass_width, pass_height, "defnot", false)

    PassFunctions(
        preload,
        setup,
        draw,
        keyPressed,
        keyReleased,
    )

    Launch()
}

func preload() {

}

func setup() {
    // Background(255, 255, 255, 255);
}

var (
    value_Z float64 = 0
    value_LR float64 = 0
)

func draw() {
    Background(246, 149, 222, 255)

    Translate(0, 0, -Width/2)

    Translate(0, Height/2, 0)
        RotateX((ReMap(MouseY+Height/2, 0, Height, 0, 360)))
    Translate(0, -Height/2, 0)
    Translate(Width/2, 0, 0)
        RotateY((ReMap(MouseX+Width/2, 0, Width, 0, 360)))
    Translate(-Width/2, 0, 0)
    // RotateY(Frames)
    // Println(-float64(int(Frames)%1000)/10)
    // Translate(0, 0, -float64(int(Frames)%1000)/10)
    Translate(value_LR, 0, value_Z)

    Fill(255, 255, 255, 255)
    Rect(0, 0, Width, Height)

    // SetCursorPos(Width/2, Height/2)

    Push()
        Translate(50, 50, 0)
        Rotate(Frames)
        RotateX(Frames)
        Fill(217, 100, 116, 100)
        Rect(-50, -25, 100, 50)
    Pop()

    Push()
        Translate(-100, 300, 0)
        Scale(1, 1.5, 1)
        Rotate(-Frames * 3)
        RotateY(-Frames * 3)
        Fill(157, 122, 222, 255)
        Rect(-50, -25, 100, 50)
    Pop()

    Rect(Width, Height/2, 100, 100)

    // PastePixels(pixels, 0, 0, Width, Height)
    Fill(85, 153, 252, 255)
    Ellipse(Width/2, Height/2, Width, Height)

    DepthRefreshOff()

        Fill(141, 205, 41, 255)
        Cube(0.5, 0.5, 0.5)

    DepthRefreshOn()

    StrokeWidth(4)
        Point(Width/2, Height/2)
    StrokeWidth(1)

    pixels := CopyPixels()

    Translate(Width/2, 0, 0)
        RotateY(-(ReMap(MouseX+Width/2, 0, Width, 0, 360)))
    Translate(-Width/2, 0, 0)
    Translate(0, Height/2, 0)
        RotateX(-(ReMap(MouseY+Height/2, 0, Height, 0, 360)))
    Translate(0, -Height/2, 0)

    Translate(0, 0, Width/2)

    Translate(0, 0, -1)

    Fill(54, 230, 198, 255)
    PastePixels(pixels, Width/4*3, Height/4*3, Width/4, Height/4)
}

func keyPressed() {
    Println(Key)
    switch ( Key ) {
        case UP:
            value_Z+=10
        case DOWN:
            value_Z-=10
        case LEFT:
            value_LR-=10
        case RIGHT:
            value_LR+=10
        case ENTER:
            ToggleFullscreen()
        case F1: Println("F1")
        case F2: Println("F2")
        case F3: Println("F3")
        case F4: Println("F4")
        case F5: Println("F5")
        case F6: Println("F6")
        case F7: Println("F7")
        case F8: Println("F8")
        case F9: Println("F9")
        case F10: Println("F10")
        case F11: Println("F11")
        case F12: Println("F12")
    }
    HideMouse()
}

func keyReleased() {
    ShowMouse()
}
