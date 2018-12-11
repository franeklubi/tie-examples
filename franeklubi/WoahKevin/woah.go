package main

import (
    . "github.com/franeklubi/tie"
    . "./Mover"
    . "github.com/franeklubi/SimpleVector"
    . "github.com/franeklubi/SmartList"
    "time"
    "strconv"
)

const (
    // animation kinds
    BURST   int = 0
    E_BURST int = 1
)

var (
    // defining bounds
    platforms   []Bound = []Bound{
        Bound{ 12, 176,  134-12, 201-176, 500, 500},
        Bound{201, 176, 298-201, 201-176, 500, 500},
        Bound{365, 176, 487-365, 201-176, 500, 500},

        Bound{100, 266, 216-100, 291-266, 500, 500},
        Bound{283, 266, 399-283, 291-266, 500, 500},

        Bound{ 70, 356,  135-70, 381-356, 500, 500},
        Bound{217, 356, 282-217, 381-356, 500, 500},
        Bound{364, 356, 429-364, 381-356, 500, 500},

        Bound{115, 446, 384-115, 471-446, 500, 500},
    }
    walls       []Bound = []Bound{
        // Bound{  0, 0,    500-0,  25-0, 500, 500},
        Bound{  0, 0,     25-0, 477-0, 500, 500},
        Bound{475, 0,  500-475, 477-0, 500, 500},
    }
    lava        []Bound = []Bound{
        Bound{0, 477, 500-0, 500-477, 500, 500},
    }
    // Bound{x, y, x2-x, y2-y, 500, 500},

    // defining b (from r,g,b) colour of bounds - thats the distinction
    nothing_b   byte = 100
    platform_b  byte = 11
    lava_b      byte = 0
    wall_b      byte = 111

    // primary textures
    background_tx   Image
    neutral_tx      Image
    kevin_tx        Image

    // enemy textures
    egg_tx          Image
    br_egg_tx       Image
    bomb_tx         Image
    enemy_1_tx      Image
    explosion_tx    Image
    coin_tx         Image

    // player
    kevin           Mover

    // enemies
    enemies         SList
    eggos           SList
    bombs           SList
    coins           SList

    next_drop       float64
    drop_timestamp  time.Time

    bmb_minus       int     = 0
    egg_minus       int     = 0
    coins_needed    int     = 135

    thr_difusal     float64 = 6000
    shield_difusal  float64 = 2000
    coins_difusal   float64 = 20000

    explosion_size  float64

    // if to pickup at the end of the frame
    pickuping       bool    = false

    // animation list
    animations      SList

    // scoring vars
    rank_calculated bool    = false
    frames_rank     float64 = 0
    frames_offset   float64 = 0
)

func main() {
    Init(1280, 720, "Woah Kevin!", false)

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
    background_tx   = LoadImage("./assets/bg.png")
    neutral_tx      = LoadImage("./assets/bg_d.png")
    kevin_tx        = LoadImage("./assets/kevin.png")

    // enemy textures
    egg_tx          = LoadImage("./assets/egg.png")
    br_egg_tx       = LoadImage("./assets/broken_egg.png")
    bomb_tx         = LoadImage("./assets/bomb.png")
    enemy_1_tx      = LoadImage("./assets/dalek.png")
    explosion_tx    = LoadImage("./assets/explosion.png")
    coin_tx         = LoadImage("./assets/coin.png")
}

func setup() {
    // setting game's width and height
    if ( Width > Height ) {
        U_width = Height
        U_height = Height
        U_zero_x = (Width-Height)/2
        U_zero_y = 0
    } else {
        U_width = Width
        U_height = Width
        U_zero_x = 0
        U_zero_y = (Height-Width)/2
    }

    // generate player
    kevin = GenPlayer(SVector{U_width/2, U_height/4, 0}, U_width*0.1, 3, kevin_tx, background_tx,
        [][]Bound{platforms, walls, lava}, []byte{nothing_b, platform_b, lava_b, wall_b})

    explosion_size = kevin.Size*0.8

    next_drop = 1000
    drop_timestamp = time.Now()
}

func draw() {
    // draw background
    saturation := ReMap(float64(kevin.Money), 0, 135, 0, 1)
    alpha := ReMap(float64(kevin.Money), 0, 135, 255, 220)
    Background(255, 255, 255, 255)

    // drawing sides if window is not square
    if ( Width != Height ) {
        NoStroke()
        Fill(165, 11, 11, 255)

        if ( Width > Height ) {
            Rect(0, 0, U_zero_x, U_height)
            Rect(U_zero_x+U_width, 0, U_zero_x, U_height)
        } else {
            Rect(0, 0, U_width, U_zero_y)
            Rect(0, U_zero_y+U_height, U_width, U_zero_y)
        }
    }

    // center the game window
    Translate(U_zero_x, U_zero_y, 0)

    Fill(HsvToRgb(0, saturation, 1, byte(alpha)))
    PastePixels(neutral_tx, 0, 0, U_width, U_height)

    // draw money in the background
    drawMoney()

    // resetting hsv
    Fill(255, 255, 255, 255)

    // draw bounds
    PastePixels(background_tx, 0, 0, U_width, U_height)

    // draw player
    Fill(255, 255, 255, 255)
    kevin.Update()
    kevin.Draw()

    // drawing HUD
    kevin.Hud()

    // drawing enemies and checking for collisions with them
    drawEnemies()

    // drawing throwables
    drawEggos()
    drawBombs()

    // drawing coins
    drawCoins()

    // drawing shield
    drawShield()

    // playing animations
    for x := 0; x < animations.Len(); x++ {
        anim := animations.Get(x).(Animation)
        playAnim(anim)

        // if anim expired remove it
        if ( durToMili(time.Since(anim.Timestamp)) > anim.Difusal ) {
            animations.Remove(x)
            animations.Execute()
        }
    }

    // if 135 coins are collected player has won
    if ( kevin.Money >= coins_needed ) {

        // drawing the you won text
        Push()
            Translate(U_width/2, U_height*0.8, 0)
                Fill(255, 255, 0, 255)
                Text("YOU WON! :D", int(U_width*0.05), true)
        Pop()

        kevin.Dead = true
    }

    // drawing the restart screen
    if ( kevin.Dead ) {

        // calculating rank
        if ( !rank_calculated ) {
            frames_rank = Frames-frames_offset
            rank_calculated = true
        }

        // drawing the restart text
        Push()
            Translate(U_width/2, U_height*0.2, 0)
                Fill(0, 255, 0, 255)
                Text("RESTART?", int(U_width*0.1), true)

            Translate(0, -U_height*0.1, 0)
                Fill(255, 255, 0, 255)
                Text("YOUR SCORE: "+strconv.Itoa(int(frames_rank))+" FRAMES", int(U_width*0.038), true)
        Pop()

        // drawing the yes/no buttons
        Push()
            Translate(U_width/2-U_width*0.35, U_height*0.55, 0)
                Fill(0, 200, 0, 255)
                Rect(0, 0, U_width*0.7, U_width*0.1)

            Translate(U_width*0.35, U_width*0.018, 0)
                Fill(255, 255, 255, 255)
                Text("PRESS ENTER TO START", int(U_width*0.03), true)
        Pop()

        // make kevin drop the currently held item
        kevin.Holding = false
        // setting kevins position to the middle of the screen
        kevin.Pos = SVector{U_width/2, U_height/2, 0}

        // waiting for user to hit enter
        if ( Key == ENTER ) {
            eggos.Clr()
            bombs.Clr()
            enemies.Clr()
            kevin.Money = 0

            kevin.Dead = false
            kevin.Holding = false
            kevin.Lives = 3

            rank_calculated = false
            frames_offset = Frames
        }

    // spawn enemies if kevin is not dead
    } else {

        if ( durToMili(time.Since(drop_timestamp)) > next_drop ) {
            drop_timestamp = time.Now()

            offset := 1000.0

            next_drop = float64(Random(1001))+offset

            bmb_or_egg := int(Random(4))

            xpos := Random(U_width*0.8)+U_width*0.1

            if ( bmb_or_egg == 0 ) {
                addBmb(xpos)
            } else {
                addEgg(xpos)
            }
        }
    }

    // normalizing holding indexes
    kevin.Which_item += bmb_minus
    bmb_minus = 0
    kevin.Which_item += egg_minus
    egg_minus = 0

    // removing obsolete objects
    eggos.Execute()
    bombs.Execute()
    enemies.Execute()
    coins.Execute()

    // picking stuff up at the very end of the game pipeline
    // ( so nothing gets deleted while we pickin' it up etc. :) )
    if ( pickuping ) {
        pickUp()
        pickuping = false
    }
}

// draws shield
func drawShield() {
    since := durToMili(time.Since(kevin.Shielded.Timestamp))
    if ( since > shield_difusal && !kevin.Dead ) {
        kevin.Shielded.Active = false
        return
    }

    // drawing shield
    Fill(9, 160, 175, 100)
    Ellipse(kevin.Pos.X, kevin.Pos.Y-kevin.Size/2, kevin.Size*1.2, kevin.Size*1.2)
    Fill(255, 255, 255, 255)
}

// draws money representation in the background
func drawMoney() {

    size := kevin.Size*0.6
    hpos := U_height*0.05
    var modif float64 = 0

    offset := U_width*0.05

    coins_wrap := 15.0

    for x := 0.0; x < float64(kevin.Money); x++ {

        for (true) {
            if ( x-modif >= coins_wrap ) {
                modif += coins_wrap
                hpos += U_height*0.1
            } else {
                break
            }
        }

        xpos := offset+size*(x-modif)

        // making coins appear red
        Fill(255, 0, 0, 255)

        // pasting coin texture
        PastePixels(coin_tx, xpos, hpos, size, size)
    }
}

// drawing enemies
func drawEnemies() {

    for x := 0; x < enemies.Len(); x++ {
        enemy := enemies.Get(x).(Mover)

        // drawing and updating enemies
        enemy.Draw()
        enemy.Update()

        if ( enemy.Dead ) {
            // adding current enemy to deletion queue
            enemies.Remove(x)

            enemy.Callback(x)

            // generating coins
            coinage := Random(4)
            for x := 0.0; x < coinage; x++ {
                addCoin(enemy.Pos)
            }
        }

        if ( kevin.Pos.Distance(enemy.Pos) < kevin.Size/2 ) {
            kevin.NewChance()
        }

        // replacing updated enemy in the enemy list
        enemies.Sit(x, enemy)
    }

}

// drawing throwables
func drawThrowables(bmb_egg bool, elapsed_callback func(x int, thr Mover)) {

    var (
        iterate_len int
        CMP         bool
    )

    if ( bmb_egg ) {
        iterate_len = eggos.Len()
        CMP = EGG
    } else {
        iterate_len = bombs.Len()
        CMP = BMB
    }

    for x := 0; x < iterate_len; x++ {
        var thr Mover
        if ( bmb_egg ) {
            thr = eggos.Get(x).(Mover)
        } else {
            thr = bombs.Get(x).(Mover)
        }

        // moving currently held item above kevins head
        if ( kevin.Holding && x == kevin.Which_item && kevin.Bmb_egg == CMP ) {
            thr.Pos = kevin.Pos
            thr.Pos.Y -= kevin.Size
        }

        elapsed := durToMili(time.Since(thr.Timestamp))

        if ( elapsed > thr_difusal ) {

            if ( x == kevin.Which_item && kevin.Bmb_egg == CMP ) {
                kevin.Holding = false
            }

            elapsed_callback(x, thr)

            if ( kevin.Holding && x < kevin.Which_item && kevin.Bmb_egg == CMP ) {
                if ( bmb_egg ) {
                    egg_minus--
                } else {
                    bmb_minus--
                }
            }
        }

        // blinking eggs
        blinking(elapsed, thr_difusal, true, 255)

        // drawing and updating an thr
        thr.Draw()
        thr.Update()

        // if thr thrown
        if ( thr.Jumping.Magnitude() > 1 ) {
            if ( bmb_egg ) {
                explode(thr.Pos, thr.Size*1.5, x, false, false)
            } else {
                explode(thr.Pos, thr.Size*1.5, 100, false, true)
            }
        }

        // replacing updated thr in the throwables list
        if ( bmb_egg ) {
            eggos.Sit(x, thr)
        } else {
            bombs.Sit(x, thr)
        }
    }

    // resetting colour after drawing throwable
    Fill(255, 255, 255, 255)
}

// drawing eggs
func drawEggos() {
    drawThrowables(true, func(x int, thr Mover) {
        // adding current egg to deletion queue
        eggos.Remove(x)

        thr.Callback(x)
    })
}

// drawing bombs
func drawBombs() {
    drawThrowables(false, func(x int, thr Mover) {
        // adding current bmb to deletion queue
        bombs.Remove(x)

        thr.Callback(x)
    })
}

// blinking anim for coloring
func blinking(elapsed, difusal float64, raise bool, alpha int) {
    // defining a quick hsv helper function
    hsv := func() {
        var c float64
        if ( raise ) {
            c = ReMap(elapsed, 0, difusal, 0, 1)
        } else {
            c = 0
        }
        Fill(HsvToRgb(0, c, 1, 255))
    }

    // this for blinking mate
    if ( elapsed > (difusal/3)*2 ) {
        if ( Mod(elapsed, 500) < 250 ) {
            hsv()
        } else {
            Fill(255, 255, 255, byte(alpha))
        }
    } else {
        hsv()
    }
}

// draws coins
func drawCoins() {
    for x := 0; x < coins.Len(); x++ {
        // retrieving a coin
        coin := coins.Get(x).(Mover)

        if ( coin.Dead ) {
            coins.Remove(x)
        }

        elapsed := durToMili(time.Since(coin.Timestamp))

        if ( elapsed > coins_difusal ) {

            // adding current bmb to deletion queue
            coins.Remove(x)
        }

        if ( kevin.Pos.Distance(coin.Pos) < kevin.Size/2 ) {
            kevin.Money++
            kevin.Money = int(Limit(float64(kevin.Money), 0, 135))

            next_drop = ReMap(float64(kevin.Money), 0, 135, 1000, 400)
            coins_difusal = ReMap(float64(kevin.Money), 0, 135, 20000, 3000)

            coins.Remove(x)
        }

        // blinking the coins
        blinking(elapsed, coins_difusal, false, 0)

        // drawing and updating a coin
        coin.Draw()
        coin.Update()

        coins.Sit(x, coin)
    }
}

// adds a coin
func addCoin(pos SVector) {
    coin := GenThrowable(kevin, coin_tx, pos, func(x int){})

    jmp_v := SVector{float64(Random(20)-10), float64(-Random(50)), 0}

    coin.Jump(jmp_v)

    coins.Put(coin)
}

// add an egg
func addEgg(x float64) {
    eggos.Put(GenThrowable(
        kevin,
        egg_tx,
        SVector{x, -10, 0},
        func(index int){
            enemyFromEgg(index)
            thr := eggos.Get(index).(Mover)
            animations.Put(Animation{thr.Pos, explosion_size, E_BURST, time.Now(), 200})
        },
    ))
}

// spawn enemy from egg
func enemyFromEgg(index int) {
    egg := eggos.Get(index).(Mover)

    addEnemy(egg.Pos)
}

// adds an enemy, duh
func addEnemy(pos SVector) {

    enemies.Put(GenEnemy(kevin,
        enemy_1_tx,
        pos,
        getDir(pos),
        func(index int) {
            anim := Animation{enemies.Get(index).(Mover).Pos, explosion_size, BURST, time.Now(), 250}
            animations.Put(anim)
        },
    ))
}

// adding bomb
func addBmb(x float64) {
    bombs.Put(GenThrowable(
        kevin,
        bomb_tx,
        SVector{x, -10, 0},
        func(index int){
            // adding burst animation at the original pos
            anim := Animation{bombs.Get(index).(Mover).Pos, explosion_size, BURST, time.Now(), 250}
            animations.Put(anim)

            explode(bombs.Get(index).(Mover).Pos, explosion_size, 100, true, true)
        },
    ))
}

func explode(pos SVector, dist float64, except int, kill_player, kill_throwables bool) {

    // kill throwables
    for x := 0; x < eggos.Len() && kill_throwables; x++ {
        if ( x == except ) {
            continue
        }

        thr := eggos.Get(x).(Mover)

        if ( pos.Distance(thr.Pos) < dist ) {
            eggos.Remove(x)

            if ( x == kevin.Which_item && kevin.Bmb_egg == EGG ) {
                kevin.Holding = false
            }

            if ( kevin.Holding && x != kevin.Which_item && x < kevin.Which_item && kevin.Bmb_egg == EGG ) {
                egg_minus--
            }

            // generating coins
            coinage := Random(2)
            for x := 0.0; x < coinage; x++ {
                addCoin(thr.Pos)
            }

            // adding burst animation
            anim := Animation{thr.Pos, thr.Size, E_BURST, time.Now(), 200}
            animations.Put(anim)
        }
    }

    // kill enemies
    for x := 0; x < enemies.Len(); x++ {
        enm := enemies.Get(x).(Mover)

        if ( pos.Distance(enm.Pos) < dist ) {
            enemies.Remove(x)

            enm.Callback(x)

            // generating coins
            coinage := Random(3)
            for x := 0.0; x < coinage; x++ {
                addCoin(enm.Pos)
            }
        }
    }

    // kill player
    // if ( kevin.Pos.Distance(pos) < dist && kill_player ) {
    if ( kevin.Pos.Distance(pos) < dist*2 && kill_player ) {
        kevin.NewChance()
    }

}

// pick up nearest object
func pickUp() {

    if ( kevin.Holding ) {
        kevin.Holding = false

        var thr Mover

        if ( kevin.Bmb_egg == BMB ) {
            thr = bombs.Get(kevin.Which_item).(Mover)
        } else {
            thr = eggos.Get(kevin.Which_item).(Mover)
        }

        var jmp_vec SVector

        if ( kevin.Reversed ) {
            jmp_vec = SVector{-30, -20, 0}
        } else {
            jmp_vec = SVector{ 30, -20, 0}
        }

        // the moving and jumping will affect throwable's trajectory
        jmp_vec = jmp_vec.Add(kevin.Moving.MultiplyN(4).Add(kevin.Jumping))

        thr.Jump(jmp_vec)

        if ( kevin.Bmb_egg == BMB ) {
            bombs.Sit(kevin.Which_item, thr)
        } else {
            eggos.Sit(kevin.Which_item, thr)
        }

        return
    }

    if ( bombs.Len() == 0 && eggos.Len() == 0 ) {
        return
    }

    // setting this up to player's size/2 so nothing farther will be caught
    var nearest_item_dist float64 = kevin.Size

    // for bombs
    for x := 0; x < bombs.Len(); x++ {
        thr := bombs.Get(x).(Mover)
        dist := kevin.Pos.Distance(thr.Pos)

        if ( dist < nearest_item_dist ) {
            nearest_item_dist = dist

            kevin.Which_item = x
            kevin.Bmb_egg = BMB
            kevin.Holding = true
        }
    }

    // for eggos
    for x := 0; x < eggos.Len(); x++ {
        thr := eggos.Get(x).(Mover)
        dist := kevin.Pos.Distance(thr.Pos)

        if ( dist < nearest_item_dist ) {
            nearest_item_dist = dist

            kevin.Which_item = x
            kevin.Bmb_egg = EGG
            kevin.Holding = true
        }
    }

}

// dir based on position
func getDir(pos SVector) (SVector) {

    var dir SVector = SVector{0, 0, 0}

    if ( pos.X > U_width/2 ) {
        dir.X -= 1
    } else {
        dir.X += 1
    }

    return dir
}

func playAnim(a Animation) {
    switch(a.Class) {
        case BURST:
            size_mult := ReMap(durToMili(time.Since(a.Timestamp)), 0, a.Difusal, 0.2, 1.5)
            burst(explosion_tx, a.Pos.Subtract(SVector{0, a.Size/2, 0}), a.Size*size_mult)

        case E_BURST:
            size_mult := ReMap(durToMili(time.Since(a.Timestamp)), 0, a.Difusal, 0.2, 1.5)
            burst(br_egg_tx, a.Pos.Subtract(SVector{0, a.Size/2, 0}), a.Size*size_mult)
    }
}

func durToMili(t time.Duration) (float64) {
    return float64(t/time.Millisecond)
}

func burst(avatar Image, pos SVector, size float64) {
    PastePixels(avatar, pos.X-size/2, pos.Y-size/2, size, size)
}

func keyPressed() {
    switch(Key) {
        case "A":   kevin.Move( SVector{-6,   0, 0} )
        case "D":   kevin.Move( SVector{ 6,   0, 0} )
        case "W":   kevin.Jump( SVector{ 0, -30, 0} )
        case " ":   pickuping = true
        case "F": ToggleFullscreen()
    }
}

func keyReleased() {
    switch(Key) {
        case "A":   kevin.Move( SVector{ 6,  0, 0} )
        case "D":   kevin.Move( SVector{-6,  0, 0} )
    }
}
