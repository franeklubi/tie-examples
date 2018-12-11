// this example is copied from https://github.com/shiffman/LearningProcessing/blob/master/chp14_transformations/example_14_17_nested_push_pop/example_14_17_nested_push_pop.pde
// it's been made by the man himself - the remarkable Dan Shiffman, and then
// rewritten by me to comply with tie's standards
package main

import (
    . "github.com/franeklubi/tie"
)

func main() {
    Init(480, 480, "nepp", false)

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

var (
    theta float64 = 0
)

func draw() {
    Background(255, 255, 255, 255);
    Stroke(0, 0, 0, 255);

    // Translate to center of window
    Translate(Width/2, Height/2, 0);

    // Loop from 0 to 360 degrees (2*PI radians)
    for i := 0.0; i < 2*PI; i += 0.2 {

        // Push, rotate and draw a line!
        // The transformation state is saved at the beginning of each cycle through the for loop and restored at the end.
        // Try commenting out these lines to see the difference!
        Push();
        Rotate(RadToDeg(theta + i));
        Line(0, 0, 100, 0);

        // Loop from 0 to 360 degrees (2*PI radians)
        for j := 0.0; j < 2*PI; j += 0.5 {
            // Push, translate, rotate and draw a line!
            Push();
            Translate(100, 0, 0);
            Rotate(RadToDeg(-theta-j));
            Line(0, 0, 50, 0);
            // We're done with the inside loop, pop!
            Pop();
        }

        // We're done with the outside loop, pop!
        Pop();
    }

    // Increment theta
    theta += 0.01;
}