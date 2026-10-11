package client

import (
	"math"

	"github.com/StoneTrench/go-mc-clone/src/game/generation"
	"github.com/StoneTrench/go-mc-clone/src/game/world"
)

func BaseTerrain(i int, r generation.PhaseRuntime) world.Chunk {
	air := r.LookupId("core:air")
	stone := r.LookupId("core:stone")
	dirt := r.LookupId("core:dirt")
	grass := r.LookupId("core:grass")
	water := r.LookupId("core:water")
	sand := r.LookupId("core:sand")
	snow := r.LookupId("core:snow")

	res := world.Chunk{}

	for p := range world.IterateChunk() {
		pos := r.ChunkOffset.Add(p)
		f := pos.ToFloat64()

		continent := math.Sin(f.X*0.01) * math.Cos(f.Z*0.01) * 24
		hills := math.Sin(f.X*0.05+1.3) * math.Cos(f.Z*0.05-0.7) * 6
		detail := math.Sin(f.X*0.21-2.1) * math.Cos(f.Z*0.19+0.9) * 1.5
		h := continent + hills + detail

		const seaLevel = 2.0

		var block world.BlockId
		switch {
		case f.Y > h:
			if f.Y <= seaLevel {
				block = water
			} else {
				block = air
			}
		case f.Y > h-1:
			switch {
			case h <= seaLevel+1:
				block = sand
			case h > 18:
				block = snow
			default:
				block = grass
			}
		case f.Y > h-4:
			block = dirt
		default:
			block = stone
		}

		res.SetBlock(p, block)
	}

	return res
}

func CarveCaves(i int, r generation.PhaseRuntime) world.Chunk {
	air := r.LookupId("core:air")
	stone := r.LookupId("core:stone")
	dirt := r.LookupId("core:dirt")
	grass := r.LookupId("core:grass")

	res := r.Results[i-1]

	for p := range world.IterateChunk() {
		b, err := res.GetBlock(p)
		if err != nil {
			continue
		}
		if b != stone && b != dirt && b != grass {
			continue
		}

		pos := r.ChunkOffset.Add(p)
		f := pos.ToFloat64()

		a := math.Sin(f.X*0.09)*math.Cos(f.Z*0.09)*math.Sin(f.Y*0.13) +
			math.Cos(f.X*0.13)*math.Sin(f.Z*0.11)*math.Cos(f.Y*0.09)
		bb := math.Sin(f.X*0.07+4.2) * math.Sin(f.Y*0.11-1.1)

		if math.Abs(a) < 0.10 && math.Abs(bb) < 0.35 {
			res.SetBlock(p, air)
		}
	}

	return res
}

func Decorate(i int, r generation.PhaseRuntime) world.Chunk {
	air := r.LookupId("core:air")
	stone := r.LookupId("core:stone")
	grass := r.LookupId("core:grass")
	iron := r.LookupId("core:iron_ore")
	log := r.LookupId("core:log")
	leaves := r.LookupId("core:leaves")

	res := r.Results[i-1]

	for p := range world.IterateChunk() {
		b, err := res.GetBlock(p)
		if err != nil {
			continue
		}

		pos := r.ChunkOffset.Add(p)
		f := pos.ToFloat64()

		switch b {
		case stone:
			v := math.Sin(f.X*0.4)*math.Cos(f.Z*0.4)*math.Sin(f.Y*0.4) +
				math.Cos(f.X*0.23)*math.Sin(f.Y*0.27)
			if f.Y < 12 && v > 0.85 {
				res.SetBlock(p, iron)
			}

		case grass:
			above := p
			above.Y++
			if b2, err := res.GetBlock(above); err != nil || b2 != air {
				continue
			}

			h := math.Sin(f.X*12.9898+f.Z*78.233) * 43758.5453
			frac := h - math.Floor(h)
			if frac < 0.985 {
				continue
			}

			for dz := int32(1); dz <= 4; dz++ {
				q := p
				q.Y += dz
				res.SetBlock(q, log)
			}

			// offsets: (dx, dy, dz) as int32
			offsets := [][3]int32{
				{-1, 0, 5}, {1, 0, 5}, {0, -1, 5}, {0, 1, 5},
				{-1, -1, 5}, {1, 1, 5}, {-1, 1, 5}, {1, -1, 5},
				{0, 0, 6},
			}
			for _, d := range offsets {
				q := p
				q.X += d[0]
				q.Z += d[1]
				q.Y += d[2]
				if b2, err := res.GetBlock(q); err == nil && b2 == air {
					res.SetBlock(q, leaves)
				}
			}
		}
	}

	return res
}
