package rendering

import (
	"fmt"
	vk "github.com/vulkan-go/vulkan"
)

func initSurface() error {
	fr_LInfo("Init surface")
	ptr, err := __window.VulkanCreateSurface(__vulkan_instance)
	if err != nil {
		return err
	}
	__surface = vk.SurfaceFromPointer(uintptr(ptr))

	return nil
}
func deinitSurface() {
	if __surface != nil {
		fr_LInfo("Deinit surface")
		vk.DestroySurface(__vulkan_instance, __surface, nil)
	}
}

func initSwapchain() error {
	fr_LInfo("Init swapchain")

	capabilities := vk.SurfaceCapabilities{}
	res := vk.GetPhysicalDeviceSurfaceCapabilities(__physical_device, __surface, &capabilities)
	err := HandleVKErrorResult(res, "failed to query surface capabilities")
	if err != nil {
		return err
	}

	formats_count := uint32(0)
	res = vk.GetPhysicalDeviceSurfaceFormats(__physical_device, __surface, &formats_count, nil)
	err = HandleVKErrorResult(res, "failed to query surface format count")
	if err != nil {
		return err
	}
	if formats_count == 0 {
		return fmt.Errorf("no surface formats found")
	}
	formats := make([]vk.SurfaceFormat, formats_count)
	res = vk.GetPhysicalDeviceSurfaceFormats(__physical_device, __surface, &formats_count, formats)
	err = HandleVKErrorResult(res, "failed to query surface formats")
	if err != nil {
		return err
	}

	presentModes_count := uint32(0)
	res = vk.GetPhysicalDeviceSurfacePresentModes(__physical_device, __surface, &presentModes_count, nil)
	err = HandleVKErrorResult(res, "failed to query surface present mode count")
	if err != nil {
		return err
	}
	if presentModes_count == 0 {
		return fmt.Errorf("no present modes found")
	}
	presentModes := make([]vk.PresentMode, presentModes_count)
	res = vk.GetPhysicalDeviceSurfacePresentModes(__physical_device, __surface, &presentModes_count, presentModes)
	err = HandleVKErrorResult(res, "failed to query surface present modes")
	if err != nil {
		return err
	}

	// Build swapchain

	// Find best format
	__swapchain_format = vk.SurfaceFormat{}
	{
		found := false
		for i := uint32(0); i < formats_count; i++ {
			formats[i].Deref()
			found = formats[i].Format == vk.FormatB8g8r8a8Srgb && formats[i].ColorSpace == vk.ColorSpaceSrgbNonlinear
			formats[i].Free()
			if found {
				__swapchain_format = formats[i]
				break
			}
		}
		if !found {
			__swapchain_format = formats[0]
		}
	}

	// Check if fifo is available
	present_mode := vk.PresentModeFifo
	{
		found := false
		for i := uint32(0); i < presentModes_count; i++ {
			if presentModes[i] == present_mode {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("failed to find eFifo present mode")
		}
	}

	__swapchain_extent = vk.Extent2D{}
	{
		capabilities.Deref()
		capabilities.CurrentExtent.Deref()
		capabilities.MinImageExtent.Deref()
		capabilities.MaxImageExtent.Deref()

		if capabilities.CurrentExtent.Height != vk.MaxUint32 {
			__swapchain_extent.Width = capabilities.CurrentExtent.Width
			__swapchain_extent.Height = capabilities.CurrentExtent.Height
		} else {
			w, h := __window.VulkanGetDrawableSize()
			__swapchain_extent.Width = min(max(uint32(w), capabilities.MinImageExtent.Width), capabilities.MaxImageExtent.Width)
			__swapchain_extent.Height = min(max(uint32(h), capabilities.MinImageExtent.Height), capabilities.MaxImageExtent.Height)
		}

		capabilities.MinImageExtent.Free()
		capabilities.MaxImageExtent.Free()
		capabilities.CurrentExtent.Free()
		capabilities.Free()
	}

	fr_LInfo("Surface extent:")
	fr_LInfo(fmt.Sprint(__swapchain_extent))

	image_count := max(3, capabilities.MinImageCount)
	if 0 < capabilities.MaxImageCount && capabilities.MaxImageCount < image_count {
		image_count = capabilities.MaxImageCount
	}

	/*
		The imageUsage bit field specifies what kind of operations we’ll use
		the images in the swap chain for. In this tutorial, we’re going to render
		directly to them, which means that they’re used as color attachment.
		It is also possible that you’ll render images to a separate image first to
		perform operations like post-processing. In that case you may use a
		value like vk::ImageUsageFlagBits::eTransferDst instead and use a memory
		operation to transfer the rendered image to a swap chain image.
	*/

	swInfo := vk.SwapchainCreateInfo{
		SType:                 vk.StructureTypeSwapchainCreateInfo,
		PNext:                 nil,
		Flags:                 0,
		Surface:               __surface,
		MinImageCount:         image_count,
		ImageFormat:           __swapchain_format.Format,
		ImageColorSpace:       __swapchain_format.ColorSpace,
		ImageExtent:           vk.Extent2D{Width: __swapchain_extent.Width, Height: __swapchain_extent.Height},
		ImageArrayLayers:      1,
		ImageUsage:            vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit),
		ImageSharingMode:      vk.SharingModeExclusive,
		QueueFamilyIndexCount: 0,
		PQueueFamilyIndices:   nil,
		PreTransform:          capabilities.CurrentTransform,
		CompositeAlpha:        vk.CompositeAlphaOpaqueBit,
		PresentMode:           present_mode,
		Clipped:               vk.True,
		OldSwapchain:          __swapchain,
	}

	__swapchain = nil
	res = vk.CreateSwapchain(__logical_device, &swInfo, nil, &__swapchain)
	err = HandleVKErrorResult(res, "failed to create swapchain")
	if err != nil {
		return err
	}
	sw_image_count := uint32(0)
	res = vk.GetSwapchainImages(__logical_device, __swapchain, &sw_image_count, nil)
	err = HandleVKErrorResult(res, "failed to get swapchain image count")
	if err != nil {
		return err
	}
	__swapchain_images = make([]vk.Image, sw_image_count)
	res = vk.GetSwapchainImages(__logical_device, __swapchain, &sw_image_count, __swapchain_images)
	err = HandleVKErrorResult(res, "failed to get swapchain images")
	if err != nil {
		return err
	}

	return nil
}
func deinitSwapchain() {
	if __swapchain != nil {
		fr_LInfo("Deinit swapchain")
		vk.DestroySwapchain(__logical_device, __swapchain, nil)
		__swapchain = nil
		__swapchain_images = nil
	}
}

func initImageViews() error {
	fr_LInfo("Init image views")

	__image_views = make([]vk.ImageView, len(__swapchain_images))
	for i, img := range __swapchain_images {
		ivInfo := vk.ImageViewCreateInfo{
			SType:    vk.StructureTypeImageViewCreateInfo,
			PNext:    nil,
			Flags:    0,
			Image:    img,
			ViewType: vk.ImageViewType2d,
			Format:   __swapchain_format.Format,
			Components: vk.ComponentMapping{
				R: vk.ComponentSwizzleIdentity,
				G: vk.ComponentSwizzleIdentity,
				B: vk.ComponentSwizzleIdentity,
				A: vk.ComponentSwizzleIdentity,
			},
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
				BaseMipLevel:   0,
				LevelCount:     1,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
		}
		res := vk.CreateImageView(__logical_device, &ivInfo, nil, &__image_views[i])
		err := HandleVKErrorResult(res, fmt.Sprintf("Failed to create image view No%d", i))
		if err != nil {
			return err
		}
	}

	return nil
}
func deinitImageViews() {
	if __image_views != nil {
		fr_LInfo("Deinit image views")
		for _, v := range __image_views {
			vk.DestroyImageView(__logical_device, v, nil)
		}
		__image_views = nil
	}
}
