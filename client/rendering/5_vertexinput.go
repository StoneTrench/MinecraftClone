package rendering

import (
	"fmt"
	"unsafe"

	. "github.com/StoneTrench/go-mc-clone/client/rendering/helpers"
	vk "github.com/vulkan-go/vulkan"
)

type SLType vk.Format

const (
	SL_uint  SLType = SLType(vk.FormatR32Uint)
	SL_uint2 SLType = SLType(vk.FormatR32g32Uint)
	SL_uint3 SLType = SLType(vk.FormatR32g32b32Uint)
	SL_uint4 SLType = SLType(vk.FormatR32g32b32a32Uint)

	SL_int  SLType = SLType(vk.FormatR32Sint)
	SL_int2 SLType = SLType(vk.FormatR32g32Sint)
	SL_int3 SLType = SLType(vk.FormatR32g32b32Sint)
	SL_int4 SLType = SLType(vk.FormatR32g32b32a32Sint)

	SL_float  SLType = SLType(vk.FormatR32Sfloat)
	SL_float2 SLType = SLType(vk.FormatR32g32Sfloat)
	SL_float3 SLType = SLType(vk.FormatR32g32b32Sfloat)
	SL_float4 SLType = SLType(vk.FormatR32g32b32a32Sfloat)

	SL_double  SLType = SLType(vk.FormatR64Sfloat)
	SL_double2 SLType = SLType(vk.FormatR64g64Sfloat)
	SL_double3 SLType = SLType(vk.FormatR64g64b64Sfloat)
	SL_double4 SLType = SLType(vk.FormatR64g64b64a64Sfloat)
)

func (s SLType) ToVulkan() vk.Format {
	return vk.Format(s)
}

func (s SLType) Sizeof() uint32 {
	switch s {
	case SL_uint:
		return 4
	case SL_uint2:
		return 8
	case SL_uint3:
		return 12
	case SL_uint4:
		return 16
	case SL_int:
		return 4
	case SL_int2:
		return 8
	case SL_int3:
		return 12
	case SL_int4:
		return 16
	case SL_float:
		return 4
	case SL_float2:
		return 8
	case SL_float3:
		return 12
	case SL_float4:
		return 16
	case SL_double:
		return 8
	case SL_double2:
		return 16
	case SL_double3:
		return 24
	case SL_double4:
		return 32
	}

	panic(fmt.Errorf("invalid SLType: %d", s))
}

type vertex2D struct {
	Pos   [2]float32
	Color [3]float32
}

var SIZEOF_VERTEX2D = SL_float2.Sizeof() + SL_float3.Sizeof()

var VERTICES = []vertex2D{
	{[2]float32{-0.5, -0.5}, [3]float32{0.0, 0.0, 1.0}},
	{[2]float32{0.5, 0.5}, [3]float32{0.0, 0.0, 1.0}},
	{[2]float32{-0.5, 0.5}, [3]float32{0.0, 0.0, 1.0}},

	{[2]float32{0.5, -0.5}, [3]float32{1.0, 1.0, 0.0}},
	{[2]float32{0.5, 0.5}, [3]float32{1.0, 1.0, 0.0}},
	{[2]float32{-0.5, -0.5}, [3]float32{1.0, 1.0, 0.0}},
}

func vertex2D_GetBindingDescription() []vk.VertexInputBindingDescription {
	desc := make([]vk.VertexInputBindingDescription, 1)

	desc[0] = vk.VertexInputBindingDescription{
		Binding:   0,
		Stride:    SIZEOF_VERTEX2D,
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

func initVertexBuffer() error {
	fr_LInfo("Init vertex buffer")
	__vertex_buffer = nil
	__vertex_buffer_memory = nil

	buffSize := vk.DeviceSize(int(SIZEOF_VERTEX2D) * len(VERTICES))

	err := InitBuffer(
		__physical_device,
		__logical_device,
		buffSize,
		vk.BufferUsageFlags(vk.BufferUsageVertexBufferBit),
		vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit|vk.MemoryPropertyHostCoherentBit),
		&__vertex_buffer,
		&__vertex_buffer_memory,
	)
	if err != nil {
		return err
	}

	var ppdata unsafe.Pointer
	res := vk.MapMemory(__logical_device, __vertex_buffer_memory, 0, buffSize, 0, &ppdata)
	err = HandleVKErrorResult(res, "failed to map memory")
	if err != nil {
		return err
	}
	gpuSlice := unsafe.Slice((*vertex2D)(ppdata), len(VERTICES))
	copy(gpuSlice, VERTICES)
	vk.UnmapMemory(__logical_device, __vertex_buffer_memory)

	return nil
}
func deinitVertexBuffer() {
	if __vertex_buffer != nil || __vertex_buffer_memory != nil {
		fr_LInfo("Deinit vertex buffer")
		DeinitBuffer(__logical_device, __vertex_buffer, __vertex_buffer_memory)
	}
}
