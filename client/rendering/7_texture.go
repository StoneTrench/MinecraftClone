package rendering

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"unsafe"

	. "github.com/StoneTrench/go-mc-clone/client/rendering/helpers"
	vk "github.com/vulkan-go/vulkan"
)

func initTextureImage() (err error) {
	file, err := os.Open("assets/textures/image.png")
	if err != nil {
		return fmt.Errorf("failed to open image, %w", err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image, %w", err)
	}
	file.Close()

	img_bounds := img.Bounds()
	img_height := img_bounds.Dy()
	img_width := img_bounds.Dx()
	img_memory_size := img_width * img_height * 4

	buffSize := vk.DeviceSize(img_memory_size)

	rgba_image := image.NewRGBA(img_bounds)

	var stagingBuffer vk.Buffer
	var stagingBufferMemory vk.DeviceMemory

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

	var gpuMappedPtr unsafe.Pointer
	res := vk.MapMemory(__logical_device, stagingBufferMemory, 0, buffSize, 0, &gpuMappedPtr)
	err = HandleVKErrorResult(res, "failed to map memory")
	if err != nil {
		return err
	}

	draw.Draw(rgba_image, img_bounds, img, image.Pt(0, 0), draw.Op(image.YCbCrSubsampleRatio410))

	gpuSlice := unsafe.Slice((*uint8)(gpuMappedPtr), img_memory_size)
	copy(gpuSlice, rgba_image.Pix)
	vk.UnmapMemory(__logical_device, stagingBufferMemory)

	//// Create texture

	err = InitImage2D(
		__physical_device,
		__logical_device,
		uint32(img_width),
		uint32(img_height),
		vk.FormatR8g8b8a8Srgb,
		vk.ImageTilingOptimal,
		vk.ImageUsageFlags(vk.ImageUsageTransferDstBit|vk.ImageUsageSampledBit),
		vk.MemoryPropertyFlags(vk.MemoryPropertyDeviceLocalBit),
		&__texture_image,
		&__index_buffer_memory,
	)
	if err != nil {
		return err
	}

	//// Allocate memory for image

	return nil
}
func deinitTextureImage() {
	if __texture_image != nil {
		vk.DestroyImage(__logical_device, __texture_image, nil)
		__texture_image = nil
	}
}
