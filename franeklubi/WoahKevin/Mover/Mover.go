package Player

import (
    . "github.com/franeklubi/tie"
    . "github.com/franeklubi/SimpleVector"
    "time"
    "strconv"
)

const (
    PULL        float64 = 10
    FRICTION    float64 = 0.9
    MOVER_WIDTH float64 = 2
    ENEMY_SIZE  float64 = 0.6

    // arm bytes
    CNTR_ARM    byte    = 0
    LEFT_ARM    byte    = 1
    RITE_ARM    byte    = 2

    BMB         bool    = true
    EGG         bool    = false

    // physics are designed to work at this resolution
    // so we have to scale to the resolution respectively when changed
    REF_X       float64 = 600
    REF_Y       float64 = 600
)

var (
    NOTHING     byte
    PLATFORM    byte
    LAVA        byte
    WALL        byte

    U_width     float64
    U_height    float64
    U_zero_x    float64
    U_zero_y    float64
)

// its the main mover struct
// im not gonna touch these and reorder them because i don't want to break
// every constructor i made
type Mover struct {
    Pos         SVector
    Size        float64
    Lives       int
    Dead        bool

    Avatar      Image
    Hitbox      Image
    Hitbox_crds [][]Bound

    Timestamp   time.Time

    Velocity    SVector
    Moving      SVector
    Jumping     SVector
    Jumped      bool
    Reversed    bool
    Holding     bool
    Shielded    Shield
    Which_item  int
    Money       int
    Bmb_egg     bool
    Callback    func(int)
}

// generate player
func GenPlayer(pos SVector, size float64, lives int, avatar, hitbox Image,
    hitbox_crds [][]Bound, bound_b []byte) (Mover) {

    if ( len(bound_b) < 4 ) {
        panic("not enough bound types")
    }
    NOTHING     = bound_b[0]
    PLATFORM    = bound_b[1]
    LAVA        = bound_b[2]
    WALL        = bound_b[3]

    return Mover{pos, size, lives, false,
        avatar, hitbox, hitbox_crds,
        time.Now(),
        SVector{0, 0, 0}, SVector{0, 0, 0}, SVector{0, 0, 0},
        false, false, false, Shield{false, time.Now()},
        0, 0, false,
        func(int){},
    }
}

// generate enemy
func GenEnemy(pla Mover, avatar Image, pos, dir SVector, callback func(int)) (Mover) {
    mover := Mover{pos, pla.Size*ENEMY_SIZE, 1, false,
        avatar, pla.Hitbox, pla.Hitbox_crds,
        time.Now(),
        SVector{0, 0, 0}, SVector{0, 0, 0}, SVector{0, 0, 0},
        false, false, false, Shield{false, time.Now()},
        0, 0, false,
        callback,
    }

    mover.Move(dir)

    return mover
}

// generate throwable object
func GenThrowable(pla Mover, avatar Image, pos SVector, callback func(int)) (Mover) {
    mover := Mover{pos, pla.Size*ENEMY_SIZE, 1, false,
        avatar, pla.Hitbox, pla.Hitbox_crds,
        time.Now(),
        SVector{0, 0, 0}, SVector{0, 10, 0}, SVector{0, 0, 0},
        false, false, false, Shield{false, time.Now()},
        0, 0, false,
        callback,
    }

    return mover
}

// draw a scene
func (m *Mover) Draw() {
    Push()

        // translating to mover position
        Translate(m.Pos.X, m.Pos.Y, 0)

            // checking if Moving the other way
            if ( m.Reversed ) {
                RotateY(180)
            }

            // drawing mover avatar
            m.Avatar.PastePixels(-m.Size/2, -m.Size, m.Size, m.Size)

    // popping transformations
    Pop()
}

func (m *Mover) Update() {

    // not updating when dead
    if ( m.Dead ) {
        return
    }

    // applying gravity
    m.Gravity()

    // adding Moving vector to Velocity
    m.Velocity = m.Velocity.Add(m.Moving)

    // adding Jumping vector to Velocity
    m.Velocity = m.Velocity.Add(m.Jumping)

    // adding FRICTION to Jumping
    m.Jumping = m.Jumping.MultiplyN(FRICTION)

    // mult to match ref resolution
    x_mult := ReMap(m.Velocity.X, 0, REF_X, 0, U_width)
    y_mult := ReMap(m.Velocity.Y, 0, REF_Y, 0, U_height)
    m.Velocity = SVector{x_mult, y_mult, 0}

    // applying Velocity
    m.Pos = m.Pos.Add(m.Velocity)

    // handle collisions
    m.collisionHandler()

    // clearing the Velocity
    m.Velocity = SVector{0, 0, 0}
}

func (m *Mover) Hud() {
    // drawing lives
    for x := 1; x <= m.Lives; x++ {
        m.Avatar.PastePixels(U_width-(m.Size*2*float64(x))/2-m.Size*0.8, m.Size*0.8, m.Size*0.8, m.Size*0.8)
    }

    // drawing things like money and etc
    Translate(U_width*0.5, 0, 0)
        Fill(126, 198, 12, 255)
        Text("Coins left: "+strconv.Itoa(135-m.Money), int(U_width*0.03), true)
    Translate(-U_width*0.5, 0, 0)

    // resetting fill
    Fill(255, 255, 255, 255)
}

// apply force
func (m *Mover) ApplyForce(sv SVector) {
    m.Velocity = m.Velocity.Add(sv)
}

// apply passed vector to the mover's Moving vector
func (m *Mover) Move(sv SVector) {
    m.Moving = m.Moving.Add(sv)
    if ( m.Moving.X < 0 ) {
        m.Reversed = true
    } else if ( m.Moving.X > 0 ) {
        m.Reversed = false
    }
}

// apply passed vector to the mover's Jumping vector
func (m *Mover) Jump(sv SVector) {
    if ( !m.Jumped ) {
        m.Jumping = m.Jumping.Add(sv)

        m.Jumped = true
    }
}

// take me to the ground boss
func (m *Mover) Gravity() {
    m.ApplyForce(SVector{0, PULL, 0})
}

// spit out bound name
func (m *Mover) boundName(sv SVector) (byte) {

    check_x := ReMap(sv.X, 0, U_width, 0, float64(m.Hitbox.W))
    check_y := ReMap(sv.Y, 0, U_height, 0, float64(m.Hitbox.H))

    var c Color = m.Hitbox.PixelAt( int(check_x), int(check_y) )
    if ( c.A != 0 ) {
        return c.B
    }

    return NOTHING
}

// check if crossed_bounds
func (m *Mover) collisionCheck() (bool, []Collision) {

    collided := false
    collisions := []Collision{}

    if ( m.Dead ) {
        return collided, collisions
    }

    for y := 0; y < len(m.Hitbox_crds); y++ {
        for x := 0; x < len(m.Hitbox_crds[y]); x++ {
            p := m.Hitbox_crds[y][x]
            x1 := ReMap(p.X, 0, float64(p.ScaleX), 0,  U_width)
            y1 := ReMap(p.Y, 0, float64(p.ScaleY), 0, U_height)
            x2 := ReMap(p.W, 0, float64(p.ScaleX), 0,  U_width)
            y2 := ReMap(p.H, 0, float64(p.ScaleY), 0, U_height)

            // Line(0, 0, m.Pos.X-m.Size/2, m.Pos.Y-m.Size/2)
            // left arm detector
            if ( insideBox(m.Pos.X-m.Size/MOVER_WIDTH, m.Pos.Y-m.Size/MOVER_WIDTH, x1, y1, x2, y2) ) {
                collided = true
                bound_name := m.boundName(m.Pos.AddN(-m.Size/MOVER_WIDTH))
                collisions = append(collisions, Collision{bound_name, LEFT_ARM, x1+x2})
            }

            // Line(0, 0, m.Pos.X+m.Size/MOVER_WIDTH, m.Pos.Y-m.Size/MOVER_WIDTH)
            // right arm detector
            if ( insideBox(m.Pos.X+m.Size/MOVER_WIDTH, m.Pos.Y-m.Size/MOVER_WIDTH, x1, y1, x2, y2) ) {
                collided = true
                bound_name := m.boundName(m.Pos.Add(SVector{m.Size/MOVER_WIDTH, -m.Size/MOVER_WIDTH, 0}))
                collisions = append(collisions, Collision{bound_name, RITE_ARM, x1})
            }

            // on-ground detector
            if ( insideBox(m.Pos.X, m.Pos.Y, x1, y1, x2, y2) ) {
                collided = true
                bound_name := m.boundName(m.Pos)
                collisions = append(collisions, Collision{bound_name, CNTR_ARM, y1})
            }
        }
    }

    return collided, collisions
}

// returns true if point is inside a box
func insideBox(x, y, b_x1, b_y1, b_x2, b_y2 float64) (bool) {
    if ( x > b_x1 && y > b_y1 && x < b_x1+b_x2 && y < b_y1+b_y2 ) {
        return true
    }
    return false
}

// collision handler
func (m *Mover) collisionHandler() {

    collided, collisions := m.collisionCheck()


    if ( !collided ) {
        return
    }

    // Println("collided")
    for x := 0; x < len(collisions); x++ {
        class   := collisions[x].class
        arm     := collisions[x].arm
        value   := collisions[x].correction

        m.Pos.X = Limit(m.Pos.X, U_width*0.1, U_width*0.91)

        switch(class) {
            case WALL:
                // m.Pos.X = m.Pos.X-m.Velocity.X
                m.Jumping.X = -m.Jumping.X*0.1

            case PLATFORM:
                // Println("platform")
                // if collided with a platform set player mover to y coord of the platform
                if ( arm == CNTR_ARM ) {
                    m.Pos.Y = value

                    // enabling jump
                    m.Jumped = false
                }

            case LAVA:
                if ( arm == CNTR_ARM ) {
                    m.Pos.Y = value

                    // enabling jump
                    m.Jumped = false
                }
                m.NewChance()

        }

    }
}

func (m *Mover) LoseLife() {
    if ( m.Lives > 1 ) {
        m.Lives--
        m.Pos = SVector{U_width/2, U_height/4, 0}
    } else {
        m.Lives--
        m.Die()
    }
}

func (m *Mover) NewChance() {
    if ( !m.Shielded.Active ) {
        m.LoseLife()
        m.Shielded = Shield{true, time.Now()}
    }
}

func (m *Mover) Die() {
    m.Dead = true
}
