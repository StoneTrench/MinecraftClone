package rendering

import (
	"unsafe"

	. "github.com/StoneTrench/go-mc-clone/client/rendering/helpers"
	vk "github.com/vulkan-go/vulkan"
)

type vertex2D struct {
	Pos   [2]float32
	Color [3]float32
}

const SIZEOF_VERTEX2D = unsafe.Sizeof(vertex2D{})

var VERTICES = []vertex2D{
	{[2]float32{-0.5, -0.5}, [3]float32{1.0, 0.0, 0.0}},
	{[2]float32{0.5, -0.5}, [3]float32{0.0, 1.0, 0.0}},
	{[2]float32{0.5, 0.5}, [3]float32{0.0, 0.0, 1.0}},
	{[2]float32{-0.5, 0.5}, [3]float32{1.0, 1.0, 1.0}},
}

var INDICES = []uint16{
	0, 1, 2, 2, 3, 0,
}

func vertex2D_GetBindingDescription() []vk.VertexInputBindingDescription {
	desc := make([]vk.VertexInputBindingDescription, 1)

	desc[0] = vk.VertexInputBindingDescription{
		Binding:   0,
		Stride:    uint32(SIZEOF_VERTEX2D),
		InputRate: vk.VertexInputRateVertex,
	}

	return desc
}
func vertex2D_GetAttributeDescription() []vk.VertexInputAttributeDescription {
	desc := make([]vk.VertexInputAttributeDescription, 2)

	desc[0] = vk.VertexInputAttributeDescription{
		Binding:  0,
		Location: 0,
		Format:   SL_float2.ToVulkan(),
		Offset:   0,
	}

	desc[1] = vk.VertexInputAttributeDescription{
		Binding:  0,
		Location: 1,
		Format:   SL_float3.ToVulkan(),
		Offset:   SL_float2.Sizeof(),
	}

	return desc
}

func initVertexBuffer() (err error) {
	fr_LInfo("Init vertex buffer")
	var stagingBuffer vk.Buffer
	var stagingBufferMemory vk.DeviceMemory

	buffSize := vk.DeviceSize(int(SIZEOF_VERTEX2D) * len(VERTICES))

	err = InitBuffer(
		__physical_device,
		__logical_device,
		buffSize,
		vk.BufferUsageFlags(vk.BufferUsageTransferSrcBit),
		vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit|vk.MemoryPropertyHostCoherentBit),
		&stagingBuffer,
		&stagingBufferMemory,
	)
	if err != nil {
		return err
	}
	defer DeinitBuffer(__logical_device, &stagingBuffer, &stagingBufferMemory)

	var ppdata unsafe.Pointer
	res := vk.MapMemory(__logical_device, stagingBufferMemory, 0, buffSize, 0, &ppdata)
	err = HandleVKErrorResult(res, "failed to map memory")
	if err != nil {
		return err
	}
	gpuSlice := unsafe.Slice((*vertex2D)(ppdata), len(VERTICES))
	copy(gpuSlice, VERTICES)
	vk.UnmapMemory(__logical_device, stagingBufferMemory)

	__vertex_buffer = nil
	__vertex_buffer_memory = nil

	err = InitBuffer(
		__physical_device,
		__logical_device,
		buffSize,
		vk.BufferUsageFlags(vk.BufferUsageTransferDstBit|vk.BufferUsageVertexBufferBit),
		vk.MemoryPropertyFlags(vk.MemoryPropertyDeviceLocalBit),
		&__vertex_buffer,
		&__vertex_buffer_memory,
	)
	if err != nil {
		return err
	}

	err = CopyBuffer(__logical_device, __logical_graphics_queue, __command_pool, __vertex_buffer, stagingBuffer, buffSize)
	if err != nil {
		return err
	}

	return nil
}
func deinitVertexBuffer() {
	if __vertex_buffer != nil || __vertex_buffer_memory != nil {
		fr_LInfo("Deinit vertex buffer")
		DeinitBuffer(__logical_device, &__vertex_buffer, &__vertex_buffer_memory)
	}
}

/*
The previous chapter already mentioned that you should allocate multiple resources
like buffers from a single memory allocation, but in fact you should go a step further.
Driver developers recommend that you also store multiple buffers, like the vertex and
index buffer, into a single VkBuffer and use offsets in commands like
vkCmdBindVertexBuffers. The advantage is that your data is more cache friendly in that
case, because it's closer together. It is even possible to reuse the same chunk of
memory for multiple resources if they are not used during the same render operations,
provided that their data is refreshed, of course. This is known as aliasing and some
Vulkan functions have explicit flags to specify that you want to do this.
*/

func initIndexBuffer() (err error) {
	fr_LInfo("Init index buffer")
	var stagingBuffer vk.Buffer
	var stagingBufferMemory vk.DeviceMemory

	buffSize := vk.DeviceSize(int(2) * len(INDICES))

	err = InitBuffer(
		__physical_device,
		__logical_device,
		buffSize,
		vk.BufferUsageFlags(vk.BufferUsageTransferSrcBit),
		vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit|vk.MemoryPropertyHostCoherentBit),
		&stagingBuffer,
		&stagingBufferMemory,
	)
	if err != nil {
		return err
	}
	defer DeinitBuffer(__logical_device, &stagingBuffer, &stagingBufferMemory)

	var ppdata unsafe.Pointer
	res := vk.MapMemory(__logical_device, stagingBufferMemory, 0, buffSize, 0, &ppdata)
	err = HandleVKErrorResult(res, "failed to map memory")
	if err != nil {
		return err
	}
	gpuSlice := unsafe.Slice((*uint16)(ppdata), len(INDICES))
	copy(gpuSlice, INDICES)
	vk.UnmapMemory(__logical_device, stagingBufferMemory)

	__index_buffer = nil
	__index_buffer_memory = nil

	err = InitBuffer(
		__physical_device,
		__logical_device,
		buffSize,
		vk.BufferUsageFlags(vk.BufferUsageTransferDstBit|vk.BufferUsageIndexBufferBit),
		vk.MemoryPropertyFlags(vk.MemoryPropertyDeviceLocalBit),
		&__index_buffer,
		&__index_buffer_memory,
	)
	if err != nil {
		return err
	}

	err = CopyBuffer(__logical_device, __logical_graphics_queue, __command_pool, __index_buffer, stagingBuffer, buffSize)
	if err != nil {
		return err
	}

	return nil
}
func deinitIndexBuffer() {
	if __index_buffer_memory != nil || __index_buffer != nil {
		fr_LInfo("Deinit index buffer")
		DeinitBuffer(__logical_device, &__index_buffer, &__index_buffer_memory)
	}
}
