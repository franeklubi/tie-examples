package main

import (
    . "github.com/franeklubi/tie"
)

func main() {
    Init(700, 700, "chaos", false)

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

type coords struct {
    x, y float64
}

var (
    points []coords
    point_colors []Color
    iterator coords = coords{Random(Width), Random(Height)}

    last int = 0

    point_width float64 = 1
    multiplier float64 = 0.5

    point_count = 4
    points_a_frame = 500
)

func setup() {
    Background(0, 0, 0, 255)
    StrokeWidth(point_width)
    Stroke(255, 255, 255, 255)

    Translate(Width/2, Height/2, 0)

    for x := 0; x < point_count; x++ {
        angle := ((PI*2)/float64(point_count))*float64(x)
        points = append(points, coords{(Width/2)*Cos(angle), (Height/2)*Sin(angle)})

        point_colors = append(point_colors, Color{byte(Random(255)), byte(Random(255)), byte(Random(255)), 255})
        c := point_colors[x]
        Stroke(c.R, c.G, c.B, 255)
        Point(float64(points[x].x), float64(points[x].y))
    }
    Stroke(255, 255, 255, 255)
    Point(iterator.x, iterator.y)
}

func draw() {
    Translate(Width/2, Height/2, 0)

    for x := 0; x < points_a_frame; x++ {
        chosen := choosingFunc()
        c := point_colors[chosen]
        Stroke(c.R, c.G, c.B, 255)

        iterator.x = LinInt(iterator.x, points[chosen].x, multiplier)
        iterator.y = LinInt(iterator.y, points[chosen].y, multiplier)
        Point(float64(iterator.x), float64(iterator.y))
    }
}

func choosingFunc() (int) {

    var chosen int
    for {
        chosen = int(Random(float64(len(points))))

        if ( chosen - last != 2  && chosen - last != -2 ) {
            last = chosen
            break
        }
    }

    return chosen
}

func keyPressed() {
    if ( Key == ENTER ) {
        Redraw()
    }
}
