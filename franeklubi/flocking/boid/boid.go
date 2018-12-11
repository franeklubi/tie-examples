package boid

import (
    . "github.com/franeklubi/tie"
    . "github.com/franeklubi/SimpleVector"
)

type Boid struct {
    Pos, Vel, Acc       SVector
    Radius, MaxF, MaxS  float64
}

func GenBoid(pos SVector, maxF, maxS float64) (Boid) {
    empty := SVector{0, 0, 0}

    return Boid{
        pos, empty, empty,
        25, maxF, maxS,
    }
}

func (b *Boid) ApplyForce(force SVector) {
    b.Acc = b.Acc.Add(force)
}

func (b *Boid) Update() {
    // apply acceleration to velocity
    b.Vel = b.Vel.Add(b.Acc)

    // limit velocity
    b.Vel = b.Vel.Limit(b.MaxS)

    // apply velocity to position
    b.Pos = b.Pos.Add(b.Vel)

    // reset acceleration to 0
    b.Acc = b.Acc.MultiplyN(0)
}

func (b *Boid) Draw() {

    pos := b.Pos
    rad := b.Radius/3

    Push()
        Translate(pos.X, pos.Y, 0)
            // drawing velocity vector
            vel_line := b.Vel.MultiplyN(10)
            Stroke(236, 36, 78, 255)
            Line(0, 0, vel_line.X, vel_line.Y)
            Stroke(0, 0, 0, 255)

            straight := SVector{0, -1, 0}
            between := straight.Angle(b.Vel)
            if ( b.Vel.X < 0 ) {
                between = -between
            }

            Rotate(RadToDeg(between))

            // drawing boid fill
            Fill(60, 218, 208, 255)
            BeginShape(TRIANGLE_FAN);
                Vertex(0, -rad*2);
                Vertex(-rad, rad*2);
                Vertex(rad, rad*2);
            EndShape();

            // drawing boid outline
            Fill(0, 0, 0, 255)
            BeginShape(LINE_LOOP);
                Vertex(0, -rad*2);
                Vertex(-rad, rad*2);
                Vertex(rad, rad*2);
            EndShape();
    Pop()
}

func (b *Boid) Wrap() {
    if (b.Pos.X < -b.Radius) {
        b.Pos.X = Width+b.Radius;
    }
    if (b.Pos.Y < -b.Radius) {
        b.Pos.Y = Height+b.Radius;
    }
    if (b.Pos.X > Width+b.Radius) {
        b.Pos.X = -b.Radius;
    }
    if (b.Pos.Y > Height+b.Radius) {
        b.Pos.Y = -b.Radius;
    }
}

// steering functions

func (b *Boid) Seek(target SVector) {
    // A vector pointing from the position to the target
    desired := target.Subtract(b.Pos);

    // Scale to maximum speed
    desired = desired.Normalize()
    desired = desired.MultiplyN(b.MaxS)

    // Steering = Desired minus velocity
    steer := desired.Subtract(b.Vel);
    steer = steer.Limit(b.MaxF);  // Limit to maximum steering force

    b.ApplyForce(steer);
}

func (b *Boid) Arrive(target SVector) {
    // A vector pointing from the position to the target
    desired := target.Subtract(b.Pos);

    d_mag := desired.Magnitude()

    if ( d_mag < 100 ) {
        new_max := ReMap(d_mag, 0, 100, 0, b.MaxS)

        desired = desired.Normalize()
        desired = desired.MultiplyN(new_max)

    } else {
        // Scale to maximum speed
        desired = desired.Normalize()
        desired = desired.MultiplyN(b.MaxS)
    }


    // Steering = Desired minus velocity
    steer := desired.Subtract(b.Vel);
    steer = steer.Limit(b.MaxF);  // Limit to maximum steering force

    b.ApplyForce(steer);
}

func (b *Boid) Separate(boids []Boid) {
    desired_separation := b.Radius*2
    var sum SVector
    count := 0

    // For every boid in the system, check if it's too close
    for x := 0; x < len(boids); x++ {
        d := b.Pos.Distance(boids[x].Pos)
        // If the distance is greater than 0 and less than an arbitrary amount (0 when you are yourself)
        if ( d > 0 && d < desired_separation ) {
            // Calculate vector pointing away from neighbor
            diff := b.Pos.Subtract(boids[x].Pos)
            diff = diff.Normalize()
            diff = diff.DivideN(d)          // Weight by distance
            sum = sum.Add(diff)
            count++             // Keep track of how many
        }
    }

    // Average -- divide by how many
    if (count > 0) {
        // Our desired vector is moving away maximum speed
        sum = sum.Normalize()
        sum = sum.MultiplyN(b.MaxS)

        // Implement Reynolds: Steering = Desired - Velocity
        steer := sum.Subtract(b.Vel)
        steer = steer.Limit(b.MaxF*2)
        b.ApplyForce(steer)
    }
}
