package client

import (
	"fmt"
	"unsafe"

	"github.com/StoneTrench/go-mat-lib/vec"
	"github.com/StoneTrench/go-mc-clone/game/regentries"
	"github.com/StoneTrench/go-mc-clone/game/resources"
	"github.com/StoneTrench/go-mc-clone/game/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var FACE_INDICIES = [6][6]int32{
	{1, 2, 6, 6, 5, 1}, // back
	{3, 0, 4, 4, 7, 3}, // front
	{7, 4, 5, 5, 6, 7}, // top
	{0, 3, 2, 2, 1, 0}, // bottom
	{0, 1, 5, 5, 4, 0}, // left
	{2, 3, 7, 7, 6, 2}, // right
}
var VOXEL_VERTICES = [8][3]float32{
	{0, 0, 0}, // 0
	{0, 0, 1}, // 1
	{1, 0, 1}, // 2
	{1, 0, 0}, // 3

	{0, 1, 0}, // 4
	{0, 1, 1}, // 5
	{1, 1, 1}, // 6
	{1, 1, 0}, // 7
}

/*

  Y
  |
  o--X
 /
/
Z

Counter Clockwise (This is what I use)
/-\
V |
Clockwise (This is not what I use)
/-\
| V

      4 --- --- 7
     /|        /|
    / |       / |
  5 --- --- 6   |
  |   |     |   |
  |   0 --- |-- 3
  |  /      |  /
  | /       | /
  1 --- --- 2

*/

func GenerateChunkMesh(level *world.Level, chunk_pos vec.Vector3[int32], reg_blocks *resources.Registry[regentries.BlockType, regentries.BlockId]) (rl.Mesh, error) {
	mesh := rl.Mesh{}

	count_vertex := int32(0)
	count_triangle := int32(0)
	arr_vertices := []float32{}
	arr_color := []uint8{}

	chunk_world_pos := world.ChunkToWorld(chunk_pos)

	chunk := level.GetChunk(chunk_pos, false)
	if chunk == nil {
		return mesh, fmt.Errorf("no chunk returned")
	}

	for x := range int32(world.CHUNK_SIZE) {
		for y := range int32(world.CHUNK_SIZE) {
			for z := range int32(world.CHUNK_SIZE) {
				local_pos := vec.InitVector3(x, y, z)

				block_id := chunk.Blocks[world.LocalToIndex(local_pos)]
				block_type, err := reg_blocks.GetById(block_id)
				if err != nil {
					return mesh, fmt.Errorf("failed to get block type for mesh, %w", err)
				}

				if block_type.IsSolid {
					for f := range 6 {
						neigh_array := world.NEIGHBOURS[f]
						face := FACE_INDICIES[f]

						can_draw_face := true

						neigh_pos := local_pos.Offset(neigh_array[0], neigh_array[1], neigh_array[2]).Add(chunk_world_pos)
						neigh_id, err := level.GetBlock(neigh_pos, false)
						if err == nil {
							neigh_type, err := reg_blocks.GetById(neigh_id)
							if err != nil {
								return mesh, fmt.Errorf("failed to get neighbour block type for mesh, %w", err)
							}
							can_draw_face = !neigh_type.IsSolid
						}

						if can_draw_face {
							for v := range 6 {
								vert := VOXEL_VERTICES[face[v]]

								vx := (float32(local_pos.X) + vert[0])
								vy := (float32(local_pos.Y) + vert[1])
								vz := (float32(local_pos.Z) + vert[2])

								vx /= 2
								vy /= 2
								vz /= 2

								vx -= 8
								vz -= 8

								arr_vertices = append(arr_vertices, vx, vy, vz)
								arr_color = append(arr_color, block_type.Color.R, block_type.Color.G, block_type.Color.B, block_type.Color.A)
							}
							count_vertex += 6
							count_triangle += 2
						}
					}
				}
			}
		}
	}

	mesh.VertexCount = count_vertex
	mesh.TriangleCount = count_triangle
	mesh.Vertices = unsafe.SliceData(arr_vertices)
	mesh.Colors = unsafe.SliceData(arr_color)
	rl.UploadMesh(&mesh, false)

	return mesh, nil
}
