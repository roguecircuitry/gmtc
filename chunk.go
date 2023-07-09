package main

type BlockBin struct {
	id uint32
}

type BlockData struct {
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
			x, y, z,
			data,
		})
	}
}
