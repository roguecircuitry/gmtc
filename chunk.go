package main

type BlockBin struct {
	id uint32
}

type BlockData struct {
	index   uint32
	x, y, z int
	data    BlockBin
}

const CHUNK_SIDE_LEN int = 8

type Chunk struct {
	blocks [CHUNK_SIDE_LEN * CHUNK_SIDE_LEN * CHUNK_SIDE_LEN]BlockBin
}

func NewChunk() *Chunk {
	return new(Chunk)
}
func posToIdx(x, y, z int) int {
	return (x +
		y*CHUNK_SIDE_LEN +
		z*CHUNK_SIDE_LEN*CHUNK_SIDE_LEN)
}
func idxToPos(index int) (int, int, int) {
	return index / (CHUNK_SIDE_LEN * CHUNK_SIDE_LEN), (index / CHUNK_SIDE_LEN) % CHUNK_SIDE_LEN, index % CHUNK_SIDE_LEN
}
func (ch *Chunk) GetBlock(x, y, z int, out *BlockData) {
	if (x < 0 || x >= CHUNK_SIDE_LEN) ||
		(y < 0 || y >= CHUNK_SIDE_LEN) ||
		(z < 0 || z >= CHUNK_SIDE_LEN) {
		out.data.id = 0
		return
	}

	index := posToIdx(x, y, z)
	out.x = x
	out.y = y
	out.z = z
	out.data = ch.blocks[index]
}
func (ch *Chunk) SetBlockBin(x, y, z int, data BlockBin) {
	index := posToIdx(x, y, z)
	ch.blocks[index] = data
}
func (ch *Chunk) SetBlockData(block *BlockData) {
	index := posToIdx(block.x, block.y, block.z)
	ch.blocks[index] = block.data
}
func (ch *Chunk) ForEachBlock(cb func(block *BlockData)) {
	for i := 0; i < len(ch.blocks); i++ {
		x, y, z := idxToPos(i)
		data := ch.blocks[i]

		cb(&BlockData{
			x: x, y: y, z: z,
			data:  data,
			index: uint32(i),
		})
	}
}

type NeighborInfo struct {
	north, south, east, west, top, bottom BlockData
}

func (ch *Chunk) GetNeighbors(x, y, z int, out *NeighborInfo) {
	ch.GetBlock(x, y+1, z, &out.top)
	ch.GetBlock(x, y-1, z, &out.bottom)
	ch.GetBlock(x, y, z+1, &out.north)
	ch.GetBlock(x, y, z-1, &out.south)
	ch.GetBlock(x+1, y, z, &out.west)
	ch.GetBlock(x-1, y, z, &out.east)
}

func BlockRevealsNeighborFaces(block *BlockData) bool {
	return block.data.id == 0
}

func NeighborToCubeInfo(neighbor *NeighborInfo, out *CubeInfo) {
	out.top = BlockRevealsNeighborFaces(&neighbor.top)
	out.bottom = BlockRevealsNeighborFaces(&neighbor.bottom)
	out.north = BlockRevealsNeighborFaces(&neighbor.north)
	out.south = BlockRevealsNeighborFaces(&neighbor.south)
	out.east = BlockRevealsNeighborFaces(&neighbor.east)
	out.west = BlockRevealsNeighborFaces(&neighbor.west)
}

func (ch *Chunk) ForEachCubeInfo(cb func(cubeInfo *CubeInfo, block *BlockData)) {
	neighbors := &NeighborInfo{}
	cubeInfo := &CubeInfo{}

	ch.ForEachBlock(func(block *BlockData) {
		ch.GetNeighbors(block.x, block.y, block.z, neighbors)
		NeighborToCubeInfo(neighbors, cubeInfo)

		cubeInfo.min.Set(float32(block.x), float32(block.y), float32(block.z))
		cubeInfo.max.Copy(&cubeInfo.min).AddScalar(1)

		cb(cubeInfo, block)
	})
}
