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
func InitBuffer(
	pdevice vk.PhysicalDevice,
	device vk.Device,
	size vk.DeviceSize,
	usage vk.BufferUsageFlags,
	properties vk.MemoryPropertyFlags,
	out_buffer *vk.Buffer,
	out_memory *vk.DeviceMemory,
) error {
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

	res := vk.CreateBuffer(device, &buffInfo, nil, out_buffer)
	err := HandleVKErrorResult(res, "failed to create buffer")
	if err != nil {
		return err
	}

	var memReq vk.MemoryRequirements
	vk.GetBufferMemoryRequirements(device, *out_buffer, &memReq)
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

	/*
	   It should be noted that in a real world application,
	   you're not supposed to actually call vkAllocateMemory
	   for every individual buffer. The maximum number of
	   simultaneous memory allocations is limited by the
	   maxMemoryAllocationCount physical device limit, which
	   may be as low as 4096 even on high end hardware like
	   an NVIDIA GTX 1080. The right way to allocate memory
	   for a large number of objects at the same time is to
	   create a custom allocator that splits up a single
	   allocation among many different objects by using the
	   offset parameters that we've seen in many functions.
	*/

	res = vk.AllocateMemory(device, &memInfo, nil, out_memory)
	err = HandleVKErrorResult(res, "failed to allocate buffer memory")
	if err != nil {
		return err
	}

	res = vk.BindBufferMemory(device, *out_buffer, *out_memory, 0)
	err = HandleVKErrorResult(res, "failed to bind buffer memory")
	if err != nil {
		return err
	}

	return nil
}
func DeinitBuffer(device vk.Device, buffer *vk.Buffer, memory *vk.DeviceMemory) {
	if buffer != nil {
		vk.DestroyBuffer(device, *buffer, nil)
		buffer = nil
	}
	if memory != nil {
		vk.FreeMemory(device, *memory, nil)
		memory = nil
	}
}
func CopyBuffer(device vk.Device, queue vk.Queue, cmdPool vk.CommandPool, dst vk.Buffer, src vk.Buffer, size vk.DeviceSize) (err error) {
	allocInfo := vk.CommandBufferAllocateInfo{
		SType:              vk.StructureTypeCommandBufferAllocateInfo,
		PNext:              nil,
		CommandPool:        cmdPool,
		Level:              vk.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}

	cmdBuf := make([]vk.CommandBuffer, 1)
	res := vk.AllocateCommandBuffers(device, &allocInfo, cmdBuf)
	err = HandleVKErrorResult(res, "failed to allocate command buffer when copying buffer")
	if err != nil {
		return err
	}

	beginInfo := vk.CommandBufferBeginInfo{
		SType:            vk.StructureTypeCommandBufferBeginInfo,
		PNext:            nil,
		Flags:            vk.CommandBufferUsageFlags(vk.CommandBufferUsageOneTimeSubmitBit),
		PInheritanceInfo: nil,
	}
	res = vk.BeginCommandBuffer(cmdBuf[0], &beginInfo)
	err = HandleVKErrorResult(res, "failed to begin command buffer when copying buffer")
	if err != nil {
		return err
	}

	copyRegion := vk.BufferCopy{
		SrcOffset: 0,
		DstOffset: 0,
		Size:      size,
	}

	vk.CmdCopyBuffer(cmdBuf[0], src, dst, 1, []vk.BufferCopy{copyRegion})
	res = vk.EndCommandBuffer(cmdBuf[0])
	err = HandleVKErrorResult(res, "failed to end command buffer when copying buffer")
	if err != nil {
		return err
	}

	/*
		We could use a fence and wait with vkWaitForFences,
		or simply wait for the transfer queue to become idle
		with vkQueueWaitIdle. A fence would allow you to schedule
		multiple transfers simultaneously and wait for all
		of them complete, instead of executing one at a time.
		That may give the driver more opportunities to optimize.
	*/

	submitInfo := vk.SubmitInfo{
		SType:                vk.StructureTypeSubmitInfo,
		PNext:                nil,
		WaitSemaphoreCount:   0,
		PWaitSemaphores:      nil,
		PWaitDstStageMask:    nil,
		CommandBufferCount:   1,
		PCommandBuffers:      cmdBuf,
		SignalSemaphoreCount: 0,
		PSignalSemaphores:    nil,
	}

	res = vk.QueueSubmit(queue, 1, []vk.SubmitInfo{submitInfo}, vk.NullFence)
	err = HandleVKErrorResult(res, "failed to submit queue when copying buffer")
	if err != nil {
		return err
	}
	res = vk.QueueWaitIdle(queue)
	err = HandleVKErrorResult(res, "failed to wait queue idle when copying buffer")
	if err != nil {
		return err
	}
	vk.FreeCommandBuffers(device, cmdPool, 1, cmdBuf)

	return nil
}

func InitImage2D(
	pdevice vk.PhysicalDevice,
	device vk.Device,
	width uint32,
	height uint32,
	format vk.Format,
	tiling vk.ImageTiling,
	usage vk.ImageUsageFlags,
	properties vk.MemoryPropertyFlags,

	image *vk.Image,
	image_memory *vk.DeviceMemory,
) (err error) {
	imageInfo := vk.ImageCreateInfo{
		SType:                 vk.StructureTypeImageCreateInfo,
		PNext:                 nil,
		Flags:                 0,
		ImageType:             vk.ImageType2d,
		Format:                format,
		MipLevels:             1,
		ArrayLayers:           1,
		Samples:               vk.SampleCount1Bit,
		Tiling:                tiling,
		Usage:                 usage,
		SharingMode:           vk.SharingModeExclusive,
		QueueFamilyIndexCount: 0,
		PQueueFamilyIndices:   nil,
		InitialLayout:         vk.ImageLayoutUndefined,
		Extent: vk.Extent3D{
			Width:  uint32(width),
			Height: uint32(height),
			Depth:  1,
		},
	}

	res := vk.CreateImage(device, &imageInfo, nil, image)
	err = HandleVKErrorResult(res, "failed to create image")
	if err != nil {
		return err
	}

	var memReq vk.MemoryRequirements
	vk.GetImageMemoryRequirements(device, *image, &memReq)

	memtype, err := findMemoryType(pdevice, memReq.MemoryTypeBits, properties)
	if err != nil {
		return err
	}

	allocInfo := vk.MemoryAllocateInfo{
		SType:           vk.StructureTypeMemoryAllocateInfo,
		PNext:           nil,
		AllocationSize:  memReq.Size,
		MemoryTypeIndex: memtype,
	}

	res = vk.AllocateMemory(device, &allocInfo, nil, image_memory)
	err = HandleVKErrorResult(res, "failed to allocate image memory")
	if err != nil {
		return err
	}

	res = vk.BindImageMemory(device, *image, *image_memory, 0)
	err = HandleVKErrorResult(res, "failed to bind image memory")
	if err != nil {
		return err
	}

	return nil
}
func DeinitImage2D(device vk.Device, image *vk.Image, memory *vk.DeviceMemory) {
	if image != nil {
		vk.DestroyImage(device, *image, nil)
		image = nil
	}
	if memory != nil {
		vk.FreeMemory(device, *memory, nil)
		memory = nil
	}
}
