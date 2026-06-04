// Package rendering the entire rendering api.
package rendering

import (
	"fmt"

	game "github.com/StoneTrench/go-mc-clone/game"
	"github.com/veandco/go-sdl2/sdl"

	. "github.com/StoneTrench/go-mc-clone/client/rendering/helpers"
	vk "github.com/vulkan-go/vulkan"
)

func fr_LInfo(msg string) {
	game.LInfo(msg, "module", "rendering")
}

const (
	MAX_FRAMES_IN_FLIGHT = 2
)

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

	__command_pool vk.CommandPool

	__vertex_buffer        vk.Buffer
	__vertex_buffer_memory vk.DeviceMemory

	__command_buffers []vk.CommandBuffer

	__sync_semaphore_wait       []vk.Semaphore
	__sync_semaphore_signal     []vk.Semaphore
	__sync_fence_inflights      []vk.Fence
	__current_frame             int
	__force_frame_buffer_reinit bool
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

	err = initVertexBuffer()
	if err != nil {
		return fmt.Errorf("failed to init vertex buffer, %w", err)
	}

	err = initCommandBuffers()
	if err != nil {
		return fmt.Errorf("failed to init command buffers, %w", err)
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

	__running = false
	deinitSyncObjects()
	deinitCommandBuffers()
	deinitVertexBuffer()
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
			Stop()
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.GetKeyFromName("space") {
				fr_LInfo("SPACE PRESSED")
			} else if t.Keysym.Sym == sdl.GetKeyFromName("escape") {
				Stop()
			}
		case *sdl.WindowEvent:
			switch t.Type {
			case sdl.WINDOWEVENT_RESIZED:
				__force_frame_buffer_reinit = true
			case sdl.WINDOWEVENT_SIZE_CHANGED:
				__force_frame_buffer_reinit = true
			}
		}
	}
}
func DrawFrames() error {
	res := vk.WaitForFences(__logical_device, 1, []vk.Fence{__sync_fence_inflights[__current_frame]}, vk.True, vk.MaxUint64)
	if err := HandleVKErrorResult(res, "failed to wait for fences"); err != nil {
		return err
	}

	var imageIndex uint32
	res = vk.AcquireNextImage(__logical_device, __swapchain, vk.MaxUint64, __sync_semaphore_wait[__current_frame], nil, &imageIndex)
	if res == vk.ErrorOutOfDate {
		err := reinitSwapchains()
		if err != nil {
			return fmt.Errorf("failed to recreate swapchain %w", err)
		}
		return nil
	} else if err := HandleVKErrorResult(res, "failed to aquire next image"); err != nil {
		return err
	}

	res = vk.ResetFences(__logical_device, 1, []vk.Fence{__sync_fence_inflights[__current_frame]})
	if err := HandleVKErrorResult(res, "failed to reset fences"); err != nil {
		return err
	}

	res = vk.ResetCommandBuffer(__command_buffers[__current_frame], 0)
	if err := HandleVKErrorResult(res, "failed to reset command buffer"); err != nil {
		return err
	}

	err := recordCommandBuffer(__command_buffers[__current_frame], imageIndex)
	if err != nil {
		return err
	}

	wait_semaphores := []vk.Semaphore{__sync_semaphore_wait[__current_frame]}
	signal_semaphores := []vk.Semaphore{__sync_semaphore_signal[imageIndex]}

	submitInfo := vk.SubmitInfo{
		SType:                vk.StructureTypeSubmitInfo,
		PNext:                nil,
		WaitSemaphoreCount:   uint32(len(wait_semaphores)),
		PWaitSemaphores:      wait_semaphores,
		PWaitDstStageMask:    []vk.PipelineStageFlags{vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit)},
		CommandBufferCount:   1,
		PCommandBuffers:      []vk.CommandBuffer{__command_buffers[__current_frame]},
		SignalSemaphoreCount: uint32(len(signal_semaphores)),
		PSignalSemaphores:    signal_semaphores,
	}

	res = vk.QueueSubmit(__logical_graphics_queue, 1, []vk.SubmitInfo{submitInfo}, __sync_fence_inflights[__current_frame])
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
	if res == vk.ErrorOutOfDate || res == vk.Suboptimal || __force_frame_buffer_reinit {
		__force_frame_buffer_reinit = false
		err := reinitSwapchains()
		if err != nil {
			return fmt.Errorf("failed to recreate swapchain %w", err)
		}
	} else if err := HandleVKErrorResult(res, "failed to present queue"); err != nil {
		return err
	}

	__current_frame = (__current_frame + 1) % MAX_FRAMES_IN_FLIGHT

	return nil
}

func Stop() {
	__running = false
}
