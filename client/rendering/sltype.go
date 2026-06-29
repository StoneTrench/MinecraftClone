package rendering

import (
	"fmt"

	// . "github.com/StoneTrench/go-mc-clone/client/rendering/helpers"
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

	SL_float4x4 SLType = SLType(vk.FormatAstc4x4SrgbBlock)
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
