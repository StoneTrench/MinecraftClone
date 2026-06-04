package helpers

import (
	"fmt"

	vk "github.com/vulkan-go/vulkan"
)

func findMemoryType(pdevice vk.PhysicalDevice, typefilter uint32, properties vk.MemoryPropertyFlags) (uint32, error) {
	var props vk.PhysicalDeviceMemoryProperties
	vk.GetPhysicalDeviceMemoryProperties(pdevice, &props)
	props.Deref()
	defer props.Free()

	for i := uint32(0); i < props.MemoryTypeCount; i++ {
		memType := props.MemoryTypes[i]
		memType.Deref()
		defer memType.Free()
		if (typefilter&(1<<i) != 0) && (memType.PropertyFlags&properties == properties) {
			return i, nil
		}
	}

	return 0, fmt.Errorf("failed to find valid memory type")
}
func InitBuffer(pdevice vk.PhysicalDevice, device vk.Device, size vk.DeviceSize, usage vk.BufferUsageFlags, properties vk.MemoryPropertyFlags, buffer *vk.Buffer, memory *vk.DeviceMemory) error {
	buffInfo := vk.BufferCreateInfo{
		SType:                 vk.StructureTypeBufferCreateInfo,
		PNext:                 nil,
		Flags:                 0,
		Size:                  size,
		Usage:                 usage,
		SharingMode:           vk.SharingModeExclusive,
		QueueFamilyIndexCount: 0,
		PQueueFamilyIndices:   nil,
	}

	res := vk.CreateBuffer(device, &buffInfo, nil, buffer)
	err := HandleVKErrorResult(res, "failed to create buffer")
	if err != nil {
		return err
	}

	var memReq vk.MemoryRequirements
	vk.GetBufferMemoryRequirements(device, *buffer, &memReq)
	memReq.Deref()
	defer memReq.Free()

	memoryTypeIndex, err := findMemoryType(pdevice, memReq.MemoryTypeBits, properties)
	if err != nil {
		return err
	}

	memInfo := vk.MemoryAllocateInfo{
		SType:           vk.StructureTypeMemoryAllocateInfo,
		PNext:           nil,
		AllocationSize:  memReq.Size,
		MemoryTypeIndex: memoryTypeIndex,
	}

	res = vk.AllocateMemory(device, &memInfo, nil, memory)
	err = HandleVKErrorResult(res, "failed to allocate buffer memory")
	if err != nil {
		return err
	}

	res = vk.BindBufferMemory(device, *buffer, *memory, 0)
	err = HandleVKErrorResult(res, "failed to bind buffer memory")
	if err != nil {
		return err
	}

	return nil
}
func DeinitBuffer(device vk.Device, buffer vk.Buffer, memory vk.DeviceMemory) {
	if buffer != nil {
		vk.DestroyBuffer(device, buffer, nil)
		buffer = nil
	}
	if memory != nil {
		vk.FreeMemory(device, memory, nil)
		memory = nil
	}
}
