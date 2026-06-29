package rendering

import (
	"math"
	"unsafe"

	. "github.com/StoneTrench/go-mat-lib"
	. "github.com/StoneTrench/go-mc-clone/client/rendering/helpers"
	vk "github.com/vulkan-go/vulkan"
)

/*
Vulkan expects the data in your structure to be aligned in memory in a specific way, for example:

	Scalars have to be aligned by N (= 4 bytes given 32 bit floats).
	A vec2 must be aligned by 2N (= 8 bytes)
	A vec3 or vec4 must be aligned by 4N (= 16 bytes)
	A nested structure must be aligned by the base alignment of its members rounded up to a multiple of 16.
	A mat4 matrix must have the same alignment as a vec4.
*/
type uniformBufferObject struct {
	model Matrix4x4[float32]
	view  Matrix4x4[float32]
	proj  Matrix4x4[float32]
}

const SIZEOF_uniformBufferObject = unsafe.Sizeof(uniformBufferObject{})

func initDescriptorSetLayout() (err error) {
	fr_LInfo("Init descriptor set layout")
	uboLayoutBinding := vk.DescriptorSetLayoutBinding{
		Binding:            0,
		DescriptorType:     vk.DescriptorTypeUniformBuffer,
		DescriptorCount:    1,
		StageFlags:         vk.ShaderStageFlags(vk.ShaderStageVertexBit), // Where it will be used
		PImmutableSamplers: nil,
	}

	layoutInfo := vk.DescriptorSetLayoutCreateInfo{
		SType:        vk.StructureTypeDescriptorSetLayoutCreateInfo,
		PNext:        nil,
		Flags:        0,
		BindingCount: 1,
		PBindings:    []vk.DescriptorSetLayoutBinding{uboLayoutBinding},
	}

	__descriptor_set_layout = nil
	res := vk.CreateDescriptorSetLayout(__logical_device, &layoutInfo, nil, &__descriptor_set_layout)
	err = HandleVKErrorResult(res, "failed to create descriptor set layout")
	if err != nil {
		return err
	}

	return nil
}
func deinitDescriptorSetLayout() {
	if __descriptor_set_layout != nil {
		fr_LInfo("Deinit descriptor set layout")
		vk.DestroyDescriptorSetLayout(__logical_device, __descriptor_set_layout, nil)
		__descriptor_set_layout = nil
	}
}

func initUniformBuffer() (err error) {
	fr_LInfo("Init uniform buffer")

	__uniform_buffer = make([]vk.Buffer, MAX_FRAMES_IN_FLIGHT)
	__uniform_buffer_memory = make([]vk.DeviceMemory, MAX_FRAMES_IN_FLIGHT)
	__uniform_buffer_mapped = make([]unsafe.Pointer, MAX_FRAMES_IN_FLIGHT)

	for i := range MAX_FRAMES_IN_FLIGHT {
		err = InitBuffer(__physical_device,
			__logical_device,
			vk.DeviceSize(SIZEOF_uniformBufferObject),
			vk.BufferUsageFlags(vk.BufferUsageUniformBufferBit),
			vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit|vk.MemoryPropertyHostCoherentBit),
			&__uniform_buffer[i],
			&__uniform_buffer_memory[i],
		)
		if err != nil {
			return err
		}

		res := vk.MapMemory(__logical_device, __uniform_buffer_memory[i], 0, vk.DeviceSize(SIZEOF_uniformBufferObject), 0, &__uniform_buffer_mapped[i])
		err = HandleVKErrorResult(res, "failed to map memory")
		if err != nil {
			return err
		}

		// We do not have to unmap the memory, because we keep using it to write to the uniform buffers.
	}

	return nil
}
func deinitUniformBuffer() {
	if __uniform_buffer != nil {

		for i := range MAX_FRAMES_IN_FLIGHT {
			DeinitBuffer(__logical_device, &__uniform_buffer[i], &__uniform_buffer_memory[i])
		}

		__uniform_buffer = nil
		__uniform_buffer_memory = nil
		__uniform_buffer_mapped = nil
	}
}
func updateUniformBuffer(frame_index int, time_current, time_delta float64) {
	mat := Matrix4x4[float32]{}
	ubo := uniformBufferObject{
		model: mat.InitRotZ(float32(time_current * math.Pi / 2)),
		view:  mat.InitLookAt(InitVector3F(2, 2, 2), InitVector3F(0, 0, 0), InitVector3F(0, 0, 1)),
		proj:  mat.InitPerspectiveProjection(45*math.Pi/180, float32(__swapchain_extent.Width), float32(__swapchain_extent.Height), 0.1, 10),
	}

	*(*uniformBufferObject)(__uniform_buffer_mapped[frame_index]) = ubo
}

func initDescriptorPool() (err error) {
	fr_LInfo("Init descriptor pool")

	poolSize := vk.DescriptorPoolSize{
		Type:            vk.DescriptorTypeUniformBuffer,
		DescriptorCount: MAX_FRAMES_IN_FLIGHT,
	}

	poolInfo := vk.DescriptorPoolCreateInfo{
		SType:         vk.StructureTypeDescriptorPoolCreateInfo,
		PNext:         nil,
		Flags:         0,
		MaxSets:       MAX_FRAMES_IN_FLIGHT,
		PoolSizeCount: 1,
		PPoolSizes:    []vk.DescriptorPoolSize{poolSize},
	}

	res := vk.CreateDescriptorPool(__logical_device, &poolInfo, nil, &__descriptor_pool)
	err = HandleVKErrorResult(res, "failed to create descriptor pool")
	if err != nil {
		return err
	}

	return nil
}
func deinitDescriptorPool() {
	if __descriptor_pool != nil {
		fr_LInfo("Deinit descriptor pool")
		vk.DestroyDescriptorPool(__logical_device, __descriptor_pool, nil)
	}
}

func initDescriptorSets() (err error) {
	layouts := make([]vk.DescriptorSetLayout, MAX_FRAMES_IN_FLIGHT)
	for i := range MAX_FRAMES_IN_FLIGHT {
		layouts[i] = __descriptor_set_layout
	}
	allocInfo := vk.DescriptorSetAllocateInfo{
		SType:              vk.StructureTypeDescriptorSetAllocateInfo,
		PNext:              nil,
		DescriptorPool:     __descriptor_pool,
		DescriptorSetCount: MAX_FRAMES_IN_FLIGHT,
		PSetLayouts:        layouts,
	}

	__descriptor_sets = make([]vk.DescriptorSet, MAX_FRAMES_IN_FLIGHT)

	res := vk.AllocateDescriptorSets(__logical_device, &allocInfo, &__descriptor_sets[0])
	err = HandleVKErrorResult(res, "failed to allocate descriptor sets")
	if err != nil {
		return err
	}

	for i := range MAX_FRAMES_IN_FLIGHT {
		bufferInfo := vk.DescriptorBufferInfo{
			Buffer: __uniform_buffer[i],
			Offset: 0,
			Range:  vk.DeviceSize(SIZEOF_uniformBufferObject),
		}

		descriptorWrite := vk.WriteDescriptorSet{
			SType:            vk.StructureTypeWriteDescriptorSet,
			PNext:            nil,
			DstSet:           __descriptor_sets[i],
			DstBinding:       0,
			DstArrayElement:  0,
			DescriptorCount:  1,
			DescriptorType:   vk.DescriptorTypeUniformBuffer,
			PImageInfo:       nil,
			PBufferInfo:      []vk.DescriptorBufferInfo{bufferInfo},
			PTexelBufferView: nil,
		}

		vk.UpdateDescriptorSets(__logical_device, 1, []vk.WriteDescriptorSet{descriptorWrite}, 0, nil)
	}

	return nil
}
func deinitDescriptorSets() {
	if __descriptor_sets != nil {
		__descriptor_sets = nil
	}
}
