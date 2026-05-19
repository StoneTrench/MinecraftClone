// Package rendering the entire rendering api.
package rendering

import (
	"fmt"
	"unsafe"

	shared "github.com/StoneTrench/go-mc-clone/shared"
	"github.com/veandco/go-sdl2/sdl"

	vk "github.com/vulkan-go/vulkan"
)

func fr_LInfo(msg string) {
	shared.LInfo(msg, "module", "rendering")
}

func NewClearColorFloat(r, g, b, a float32) vk.ClearColorValue {
	var c vk.ClearColorValue
	ptr := (*[4]float32)(unsafe.Pointer(&c))
	ptr[0], ptr[1], ptr[2], ptr[3] = r, g, b, a
	return c
}

func HandleVKErrorResult(res vk.Result, msg string) error {
	if res != vk.Success {
		return fmt.Errorf("%s, vulkan error code %d", msg, res)
	}
	return nil
}

// State
var (
	_running bool

	__window          *sdl.Window
	__vulkan_instance vk.Instance
	__surface         vk.Surface
	__physical_device vk.PhysicalDevice
	__logical_device  vk.Device
	__queue_graphics  vk.Queue

	__extent         vk.Extent2D
	__surface_format vk.SurfaceFormat

	__swapchain        vk.Swapchain
	__swapchain_images []vk.Image

	__image_views []vk.ImageView
)

func Init() (err error) {
	defer func() {
		_running = true
		if err != nil {
			Deinit()
		}
	}()

	err = initWindow(800, 600)
	if err != nil {
		return fmt.Errorf("failed to init window, %w", err)
	}

	err = initVulkan()
	if err != nil {
		return fmt.Errorf("failed to init vulkan, %w", err)
	}

	err = initVkInstance()
	if err != nil {
		return fmt.Errorf("failed to init vulkan instance, %w", err)
	}

	err = initSurface()
	if err != nil {
		return fmt.Errorf("failed to init surface, %w", err)
	}

	err = initPhysicalDevices()
	if err != nil {
		return fmt.Errorf("failed to init physical devices, %w", err)
	}

	err = initLogicalDevice()
	if err != nil {
		return fmt.Errorf("failed to init logical device, %w", err)
	}

	err = initSwapchain()
	if err != nil {
		return fmt.Errorf("failed to init swapchain, %w", err)
	}

	err = initImageViews()
	if err != nil {
		return fmt.Errorf("failed to init image views, %w", err)
	}

	return nil
}
func Deinit() {
	_running = false
	deinitImageViews()
	deinitSwapchain()
	deinitLogicalDevice()
	deinitPhysicalDevices()
	deinitSurface()
	deinitVkInstance()
	deinitVulkan()
	deinitWindow()
}

func IsRunning() bool {
	return _running
}
func PollEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			_running = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.GetKeyFromName("space") {
				fr_LInfo("SPACE PRESSED")
			} else if t.Keysym.Sym == sdl.GetKeyFromName("escape") {
				_running = false
			}
		}
	}
}

func initWindow(width, height uint32) (err error) {
	fr_LInfo("Initializing SDL2")
	err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		return err
	}

	fr_LInfo("Loading vulkan dll")
	err = sdl.VulkanLoadLibrary("vulkan-1.dll")
	if err != nil {
		return err
	}

	fr_LInfo("Initializing window")
	__window, err = sdl.CreateWindow(
		fmt.Sprintf("%s  —  v%d.%d.%d", shared.PROJECT_NAME, shared.PROJECT_VERSION_MAJOR, shared.PROJECT_VERSION_MINOR, shared.PROJECT_VERSION_PATCH),
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(width), int32(height),
		sdl.WINDOW_VULKAN,
	)
	if err != nil {
		return err
	}

	return nil
}
func deinitWindow() {
	if __window != nil {
		fr_LInfo("Deinitializing window")
		__window.Destroy()
		__window = nil
	}
	fr_LInfo("Unloading vulkan dll")
	sdl.VulkanUnloadLibrary()
	fr_LInfo("Deinitializing SDL2")
	sdl.Quit()
}

func initVulkan() error {
	fr_LInfo("Initializing vulkan")
	procAddr := sdl.VulkanGetVkGetInstanceProcAddr()
	if procAddr == nil {
		return fmt.Errorf("failed to get VkGetInstanceProcAddr from SDL2")
	}
	vk.SetGetInstanceProcAddr(procAddr)

	err := vk.Init()
	if err != nil {
		return err
	}

	return nil
}
func deinitVulkan() {
	fr_LInfo("Deinitializing vulkan")
	// nothing
}

func initVkInstance() error {
	fr_LInfo("Initializing vulkan instance")

	application_info := vk.ApplicationInfo{
		SType:              vk.StructureTypeApplicationInfo,
		PNext:              nil,
		PApplicationName:   shared.PROJECT_NAME + "\x00",
		ApplicationVersion: vk.MakeVersion(shared.PROJECT_VERSION_MAJOR, shared.PROJECT_VERSION_MINOR, shared.PROJECT_VERSION_PATCH),
		PEngineName:        "No Engine\x00",
		EngineVersion:      vk.MakeVersion(1, 0, 0),
		ApiVersion:         vk.ApiVersion10,
	}

	layers := []string{
		// "VK_LAYER_KHRONOS_validation",
	}
	extensions := __window.VulkanGetInstanceExtensions()

	fr_LInfo("Layers: " + fmt.Sprint(layers))
	fr_LInfo("Extensions: " + fmt.Sprint(extensions))

	instance_info := vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PNext:                   nil,
		Flags:                   0,
		PApplicationInfo:        &application_info,
		EnabledLayerCount:       uint32(len(layers)),
		PpEnabledLayerNames:     layers,
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledExtensionNames: extensions,
	}
	res := vk.CreateInstance(&instance_info, nil, &__vulkan_instance)
	err := HandleVKErrorResult(res, "failed to create vulkan instance")
	if err != nil {
		return err
	}

	err = vk.InitInstance(__vulkan_instance)
	if err != nil {
		return err
	}

	return nil
}
func deinitVkInstance() {
	if __vulkan_instance != nil {
		fr_LInfo("Deinitializing vulkan instance")
		vk.DestroyInstance(__vulkan_instance, nil)
		__vulkan_instance = nil
	}
}

func initPhysicalDevices() error {
	fr_LInfo("Initializing physical devices")

	count := uint32(0)
	res := vk.EnumeratePhysicalDevices(__vulkan_instance, &count, nil)
	err := HandleVKErrorResult(res, "failed to enumerate physical devices")
	if err != nil {
		return err
	}
	devices := make([]vk.PhysicalDevice, count)
	res = vk.EnumeratePhysicalDevices(__vulkan_instance, &count, devices)
	err = HandleVKErrorResult(res, "failed to enumerate physical devices")
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("failed to find gpus with vulkan support")
	}

	// Pick physical device
	__physical_device = nil
	for i := uint32(0); i < count; i++ {
		pdev := devices[i]

		if isDeviceSuitable(pdev) {
			__physical_device = pdev
			break
		}
	}

	if __physical_device == nil {
		return fmt.Errorf("failed to find suitable physical device")
	}

	return nil
}
func isDeviceSuitable(pdev vk.PhysicalDevice) bool {
	properties := vk.PhysicalDeviceProperties{}
	vk.GetPhysicalDeviceProperties(pdev, &properties)
	properties.Deref()
	defer properties.Free()

	if properties.ApiVersion < vk.ApiVersion10 {
		return false
	}

	if properties.DeviceType != vk.PhysicalDeviceTypeDiscreteGpu {
		return false
	}

	features := vk.PhysicalDeviceFeatures{}
	vk.GetPhysicalDeviceFeatures(pdev, &features)
	features.Deref()
	defer features.Free()

	if features.GeometryShader == vk.False {
		return false
	}

	_, err := getValidQueueFamilyIndex(pdev)
	return err == nil
}
func getValidQueueFamilyIndex(pdev vk.PhysicalDevice) (uint32, error) {
	queue_props_count := uint32(0)
	vk.GetPhysicalDeviceQueueFamilyProperties(pdev, &queue_props_count, nil)
	queue_props := make([]vk.QueueFamilyProperties, queue_props_count)
	vk.GetPhysicalDeviceQueueFamilyProperties(pdev, &queue_props_count, queue_props)

	for i := uint32(0); i < queue_props_count; i++ {
		queue_props[i].Deref()
		var is_surface_support vk.Bool32 = vk.False
		vk.GetPhysicalDeviceSurfaceSupport(pdev, i, __surface, &is_surface_support)
		is_valid := queue_props[i].QueueFlags&vk.QueueFlags(vk.QueueGraphicsBit) != 0 && (is_surface_support == vk.True)

		queue_props[i].Free()

		if is_valid {
			return i, nil
		}
	}

	return 0, fmt.Errorf("failed to find valid family index")
}
func deinitPhysicalDevices() {
	if __physical_device != nil {
		fr_LInfo("Deinitializing physical devices")
		__physical_device = nil
	}
}

func initLogicalDevice() error {
	fr_LInfo("Initializing logical device")

	family_index, err := getValidQueueFamilyIndex(__physical_device)
	if err != nil {
		return err
	}

	devQueueInfo := vk.DeviceQueueCreateInfo{
		SType:            vk.StructureTypeDeviceQueueCreateInfo,
		PNext:            nil,
		Flags:            0,
		QueueFamilyIndex: family_index,
		QueueCount:       1,
		PQueuePriorities: []float32{0.5},
	}

	deviceFeatures := vk.PhysicalDeviceFeatures{}

	extensions := []string{"VK_KHR_swapchain\x00"}

	devInfo := vk.DeviceCreateInfo{
		SType:                   vk.StructureTypeDeviceCreateInfo,
		PNext:                   nil,
		Flags:                   0,
		QueueCreateInfoCount:    1,
		PQueueCreateInfos:       []vk.DeviceQueueCreateInfo{devQueueInfo},
		EnabledLayerCount:       0,
		PpEnabledLayerNames:     nil,
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledExtensionNames: extensions,
		PEnabledFeatures:        []vk.PhysicalDeviceFeatures{deviceFeatures},
	}

	res := vk.CreateDevice(__physical_device, &devInfo, nil, &__logical_device)
	err = HandleVKErrorResult(res, "failed to create logical device")
	if err != nil {
		return err
	}

	vk.GetDeviceQueue(__logical_device, family_index, 0, &__queue_graphics)

	return nil
}
func deinitLogicalDevice() {
	if __logical_device != nil {
		fr_LInfo("Deinitializing logical device")
		vk.DestroyDevice(__logical_device, nil)
		__logical_device = nil
	}
}

func initSurface() error {
	fr_LInfo("Initializing surface")
	ptr, err := __window.VulkanCreateSurface(__vulkan_instance)
	if err != nil {
		return err
	}
	__surface = vk.SurfaceFromPointer(uintptr(ptr))

	return nil
}
func deinitSurface() {
	if __surface != nil {
		fr_LInfo("Deinitializing surface")
		vk.DestroySurface(__vulkan_instance, __surface, nil)
	}
}

func initSwapchain() error {
	fr_LInfo("Initializing swapchain")

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
	format := vk.SurfaceFormat{}
	{
		found := false
		for i := uint32(0); i < formats_count; i++ {
			formats[i].Deref()
			found = formats[i].Format == vk.FormatB8g8r8a8Srgb && formats[i].ColorSpace == vk.ColorSpaceSrgbNonlinear
			formats[i].Free()
			if found {
				format = formats[i]
				break
			}
		}
		if !found {
			format = formats[0]
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

	extent := vk.Extent2D{}
	{
		capabilities.CurrentExtent.Deref()
		if capabilities.CurrentExtent.Height != vk.MaxUint32 {
			extent.Width = capabilities.CurrentExtent.Width
			extent.Height = capabilities.CurrentExtent.Height
		} else {
			w, h := __window.VulkanGetDrawableSize()
			extent.Width = min(max(uint32(w), capabilities.MinImageExtent.Width), capabilities.MaxImageExtent.Width)
			extent.Height = min(max(uint32(h), capabilities.MinImageExtent.Height), capabilities.MaxImageExtent.Height)
		}

		capabilities.CurrentExtent.Free()
	}

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

	__extent = extent
	__surface_format = format

	swInfo := vk.SwapchainCreateInfo{
		SType:                 vk.StructureTypeSwapchainCreateInfo,
		PNext:                 nil,
		Flags:                 0,
		Surface:               __surface,
		MinImageCount:         image_count,
		ImageFormat:           format.Format,
		ImageColorSpace:       format.ColorSpace,
		ImageExtent:           extent,
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
		fr_LInfo("Deinitializing swapchain")
		vk.DestroySwapchain(__logical_device, __swapchain, nil)
		__swapchain = nil
		__swapchain_images = nil
	}
}

func initImageViews() error {
	return nil
}
func deinitImageViews() {

}
