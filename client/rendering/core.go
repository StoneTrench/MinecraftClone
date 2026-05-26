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

func VKResultToString(res vk.Result) string {
	switch res {
	case vk.Success:
		return "Success"
	case vk.NotReady:
		return "NotReady"
	case vk.Timeout:
		return "Timeout"
	case vk.EventSet:
		return "EventSet"
	case vk.EventReset:
		return "EventReset"
	case vk.Incomplete:
		return "Incomplete"
	case vk.ErrorOutOfHostMemory:
		return "ErrorOutOfHostMemory"
	case vk.ErrorOutOfDeviceMemory:
		return "ErrorOutOfDeviceMemory"
	case vk.ErrorInitializationFailed:
		return "ErrorInitializationFailed"
	case vk.ErrorDeviceLost:
		return "ErrorDeviceLost"
	case vk.ErrorMemoryMapFailed:
		return "ErrorMemoryMapFailed"
	case vk.ErrorLayerNotPresent:
		return "ErrorLayerNotPresent"
	case vk.ErrorExtensionNotPresent:
		return "ErrorExtensionNotPresent"
	case vk.ErrorFeatureNotPresent:
		return "ErrorFeatureNotPresent"
	case vk.ErrorIncompatibleDriver:
		return "ErrorIncompatibleDriver"
	case vk.ErrorTooManyObjects:
		return "ErrorTooManyObjects"
	case vk.ErrorFormatNotSupported:
		return "ErrorFormatNotSupported"
	case vk.ErrorFragmentedPool:
		return "ErrorFragmentedPool"
	case vk.ErrorOutOfPoolMemory:
		return "ErrorOutOfPoolMemory"
	case vk.ErrorInvalidExternalHandle:
		return "ErrorInvalidExternalHandle"
	case vk.ErrorSurfaceLost:
		return "ErrorSurfaceLost"
	case vk.ErrorNativeWindowInUse:
		return "ErrorNativeWindowInUse"
	case vk.Suboptimal:
		return "Suboptimal"
	case vk.ErrorOutOfDate:
		return "ErrorOutOfDate"
	case vk.ErrorIncompatibleDisplay:
		return "ErrorIncompatibleDisplay"
	case vk.ErrorValidationFailed:
		return "ErrorValidationFailed"
	case vk.ErrorInvalidShaderNv:
		return "ErrorInvalidShaderNv"
	case vk.ErrorInvalidDrmFormatModifierPlaneLayout:
		return "ErrorInvalidDrmFormatModifierPlaneLayout"
	case vk.ErrorFragmentation:
		return "ErrorFragmentation"
	case vk.ErrorNotPermitted:
		return "ErrorNotPermitted"
	// case vk.ResultBeginRange: return "ResultBeginRange"
	// case vk.ResultEndRange: return "ResultEndRange"
	case vk.ResultRangeSize:
		return "ResultRangeSize"
	case vk.ResultMaxEnum:
		return "ResultMaxEnum"
	default:
		panic("invalid vulkan result")
	}
}

func HandleVKErrorResult(res vk.Result, msg string) error {
	if res != vk.Success {
		return fmt.Errorf("%s, vulkan error code (%d) %s", msg, res, VKResultToString(res))
	}
	return nil
}

var (
	__running bool

	__window *sdl.Window

	__vulkan_instance vk.Instance

	__surface vk.Surface

	__physical_device vk.PhysicalDevice

	__logical_device         vk.Device
	__logical_graphics_queue vk.Queue

	__swapchain_extent vk.Extent2D
	__swapchain_format vk.SurfaceFormat
	__swapchain        vk.Swapchain
	__swapchain_images []vk.Image

	__image_views []vk.ImageView

	__render_pass vk.RenderPass

	__pipeline_shader_module vk.ShaderModule
	__pipeline_layout        vk.PipelineLayout
	__pipeline               vk.Pipeline
	__pipeline_cache         vk.PipelineCache

	__framebuffers []vk.Framebuffer

	__command_pool   vk.CommandPool
	__command_buffer vk.CommandBuffer

	__sync_semaphore_image_available vk.Semaphore
	__sync_semaphore_render_finished vk.Semaphore
	__sync_fence_inflight            vk.Fence
)

func Init() (err error) {
	defer func() {
		__running = true
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

	err = initRenderPass()
	if err != nil {
		return fmt.Errorf("failed to init render pass, %w", err)
	}

	err = initPipeline()
	if err != nil {
		return fmt.Errorf("failed to init vulkan pipeline, %w", err)
	}

	err = initFrameBuffers()
	if err != nil {
		return fmt.Errorf("failed to init framebuffers, %w", err)
	}

	err = initCommandPool()
	if err != nil {
		return fmt.Errorf("failed to init command pool, %w", err)
	}

	err = initSyncObjects()
	if err != nil {
		return fmt.Errorf("failed to init synchronization objects, %w", err)
	}

	return nil
}
func Deinit() {
	if __logical_device != nil {
		vk.DeviceWaitIdle(__logical_device)
	}
	// if __logical_device != nil && __sync_fence_inflight != nil {
	// 	vk.WaitForFences(__logical_device, 1, []vk.Fence{__sync_fence_inflight}, vk.True, 1000e+6)
	// }

	__running = false
	deinitSyncObjects()
	deinitCommandPool()
	deinitFrameBuffers()
	deinitPipeline()
	deinitRenderPass()
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
	return __running
}
func PollEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			__running = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.GetKeyFromName("space") {
				fr_LInfo("SPACE PRESSED")
			} else if t.Keysym.Sym == sdl.GetKeyFromName("escape") {
				__running = false
			}
		}
	}
}
func DrawFrames() error {
	res := vk.WaitForFences(__logical_device, 1, []vk.Fence{__sync_fence_inflight}, vk.True, vk.MaxUint64)
	if err := HandleVKErrorResult(res, "failed to wait for fences"); err != nil {
		return err
	}
	res = vk.ResetFences(__logical_device, 1, []vk.Fence{__sync_fence_inflight})
	if err := HandleVKErrorResult(res, "failed to reset fences"); err != nil {
		return err
	}

	var imageIndex uint32
	res = vk.AcquireNextImage(__logical_device, __swapchain, vk.MaxUint64, __sync_semaphore_image_available, nil, &imageIndex)
	// if err := HandleVKErrorResult(res, "failed to acquire next image"); err != nil {
	// 	return err
	// }

	res = vk.ResetCommandBuffer(__command_buffer, 0)
	if err := HandleVKErrorResult(res, "failed to reset command buffer"); err != nil {
		return err
	}

	err := recordCommandBuffer(__command_buffer, imageIndex)
	if err != nil {
		return err
	}

	wait_semaphores := []vk.Semaphore{__sync_semaphore_image_available}
	signal_semaphores := []vk.Semaphore{__sync_semaphore_render_finished}

	submitInfo := vk.SubmitInfo{
		SType:                vk.StructureTypeSubmitInfo,
		PNext:                nil,
		WaitSemaphoreCount:   uint32(len(wait_semaphores)),
		PWaitSemaphores:      wait_semaphores,
		PWaitDstStageMask:    []vk.PipelineStageFlags{vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit)},
		CommandBufferCount:   1,
		PCommandBuffers:      []vk.CommandBuffer{__command_buffer},
		SignalSemaphoreCount: uint32(len(signal_semaphores)),
		PSignalSemaphores:    signal_semaphores,
	}

	res = vk.QueueSubmit(__logical_graphics_queue, 1, []vk.SubmitInfo{submitInfo}, __sync_fence_inflight)
	if err := HandleVKErrorResult(res, "failed to submit draw command buffer to queue"); err != nil {
		return err
	}

	// End of the drawframe function
	presentInfo := vk.PresentInfo{
		SType:              vk.StructureTypePresentInfo,
		PNext:              nil,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    signal_semaphores,
		SwapchainCount:     1,
		PSwapchains:        []vk.Swapchain{__swapchain},
		PImageIndices:      []uint32{imageIndex},
		PResults:           nil,
	}

	res = vk.QueuePresent(__logical_graphics_queue, &presentInfo)
	// if err := HandleVKErrorResult(res, "failed to present the queue"); err != nil {
	// 	return err
	// }

	return nil
}
