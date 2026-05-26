package rendering

import (
	"fmt"

	shared "github.com/StoneTrench/go-mc-clone/shared"
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"
)

func initWindow(width, height uint32) (err error) {
	fr_LInfo("Init SDL2")
	err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		return err
	}

	fr_LInfo("Loading vulkan dll")
	err = sdl.VulkanLoadLibrary("vulkan-1.dll")
	if err != nil {
		return err
	}

	fr_LInfo("Init window")
	__window, err = sdl.CreateWindow(
		fmt.Sprintf("%s  —  v%d.%d.%d", shared.PROJECT_NAME, shared.PROJECT_VERSION_MAJOR, shared.PROJECT_VERSION_MINOR, shared.PROJECT_VERSION_PATCH),
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(width), int32(height),
		sdl.WINDOW_VULKAN|sdl.WINDOW_RESIZABLE,
	)
	if err != nil {
		return err
	}

	return nil
}
func deinitWindow() {
	if __window != nil {
		fr_LInfo("Deinit window")
		__window.Destroy()
		__window = nil
	}
	fr_LInfo("Unloading vulkan dll")
	sdl.VulkanUnloadLibrary()
	fr_LInfo("Deinit SDL2")
	sdl.Quit()
}

func initVulkan() error {
	fr_LInfo("Init vulkan")
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
	fr_LInfo("Deinit vulkan")
	// nothing
}

func initVkInstance() error {
	fr_LInfo("Init vulkan instance")

	application_info := vk.ApplicationInfo{
		SType:              vk.StructureTypeApplicationInfo,
		PNext:              nil,
		PApplicationName:   shared.PROJECT_NAME + "\x00",
		ApplicationVersion: vk.MakeVersion(shared.PROJECT_VERSION_MAJOR, shared.PROJECT_VERSION_MINOR, shared.PROJECT_VERSION_PATCH),
		PEngineName:        "No Engine\x00",
		EngineVersion:      vk.MakeVersion(1, 0, 0),
		ApiVersion:         vk.ApiVersion11,
	}

	layers := []string{
		"VK_LAYER_KHRONOS_validation\x00",
	}
	extensions := __window.VulkanGetInstanceExtensions()
	// extensions = append(extensions, "VK_KHR_shader_draw_parameters\x00")

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
		fr_LInfo("Deinit vulkan instance")
		vk.DestroyInstance(__vulkan_instance, nil)
		__vulkan_instance = nil
	}
}

func initPhysicalDevices() error {
	fr_LInfo("Init physical devices")

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
		fr_LInfo("Deinit physical devices")
		__physical_device = nil
	}
}

func initLogicalDevice() error {
	fr_LInfo("Init logical device")

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
	extensions := []string{
		"VK_KHR_swapchain\x00",
		vk.KhrShaderDrawParametersExtensionName + "\x00",
	}

	PdevQueueInfos := []vk.DeviceQueueCreateInfo{devQueueInfo}
	PdeviceFeatures := []vk.PhysicalDeviceFeatures{deviceFeatures}

	devInfo := vk.DeviceCreateInfo{
		SType:                   vk.StructureTypeDeviceCreateInfo,
		PNext:                   nil,
		Flags:                   0,
		QueueCreateInfoCount:    uint32(len(PdevQueueInfos)),
		PQueueCreateInfos:       PdevQueueInfos,
		EnabledLayerCount:       0,
		PpEnabledLayerNames:     nil,
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledExtensionNames: extensions,
		PEnabledFeatures:        PdeviceFeatures,
	}

	__logical_device = nil
	res := vk.CreateDevice(__physical_device, &devInfo, nil, &__logical_device)
	err = HandleVKErrorResult(res, "failed to create logical device")
	if err != nil {
		return err
	}

	vk.GetDeviceQueue(__logical_device, family_index, 0, &__logical_graphics_queue)

	return nil
}
func deinitLogicalDevice() {
	if __logical_device != nil {
		fr_LInfo("Deinit logical device")
		vk.DestroyDevice(__logical_device, nil)
		__logical_device = nil
	}
}
