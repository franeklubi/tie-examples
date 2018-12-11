package main

// import the package
import (
    "github.com/franeklubi/tie"
)

func main() {
    // initialize engine in main
    tie.Init(500, 500, "window_name", false)
    //           width, height, window_name, is_resizable

    // pass all the functions you want to be used by the engine
    tie.PassFunctions(
        preload,
        setup,
        draw,
    )

    // launch the engine
    tie.Launch()
}

// called only once, before setup, nothing can be drawn here
func preload() {

}

// called only once, before draw, you can draw here
func setup() {

}

// called once every frame
func draw() {

}
