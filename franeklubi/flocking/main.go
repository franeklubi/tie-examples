package main

import (
    . "github.com/franeklubi/tie"
    . "github.com/franeklubi/SimpleVector"
    . "./boid"
)

func main() {
    Init(720, 720, "boid", false)

    PassFunctions(
        preload,
        setup,
        draw,
    )

    Launch()
}

func preload() {

}

var (
    spks   []Boid
)

func setup() {
    // generating boids
    for x := 0; x < 50; x++ {
        x, y := float64(Random(Width)), float64(Random(Height))
        szpk := GenBoid(SVector{x, y, 0}, 0.1, 4)
        spks = append(spks, szpk)
    }
}

func Pipe(szpk *Boid) {
    szpk.Update()
    szpk.Wrap()
    szpk.Draw()
}

func draw() {
    Background(255, 255, 255, 255)

    for x := 0; x < len(spks); x++ {
        spks[x].Arrive(SVector{MouseX, MouseY, 0})
        spks[x].Separate(spks)

        Pipe(&spks[x])
    }
}
