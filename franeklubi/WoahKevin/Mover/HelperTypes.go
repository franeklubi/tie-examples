package Player

import (
    . "github.com/franeklubi/SimpleVector"
    "time"
)

type Bound struct {
    X, Y, W, H, ScaleX, ScaleY float64
}

type Collision struct {
    class       byte
    arm         byte
    correction  float64
}

type Animation struct {
    Pos         SVector
    Size        float64
    Class       int
    Timestamp   time.Time
    Difusal     float64
}

type Shield struct {
    Active      bool
    Timestamp   time.Time
}
