package rendering

import (
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"
	. "github.com/StoneTrench/go-mc-clone/client/rendering/helpers"
)

func initFrameBuffers() error {
	fr_LInfo("Init framebuffers")
	__framebuffers = make([]vk.Framebuffer, len(__image_views))

	for i, v := range __image_views {
		attachments := []vk.ImageView{v}

		frameInfo := vk.FramebufferCreateInfo{
			SType:           vk.StructureTypeFramebufferCreateInfo,
			PNext:           nil,
			Flags:           0,
			RenderPass:      __render_pass,
			AttachmentCount: uint32(len(attachments)),
			PAttachments:    attachments,
			Width:           __swapchain_extent.Width,
			Height:          __swapchain_extent.Height,
			Layers:          1,
		}

		res := vk.CreateFramebuffer(__logical_device, &frameInfo, nil, &__framebuffers[i])
		err := HandleVKErrorResult(res, "failed to create framebuffer")
		if err != nil {
			return err
		}
	}

	return nil
}
func deinitFrameBuffers() {
	if __framebuffers != nil {
		fr_LInfo("Deinit framebuffers")
		for _, v := range __framebuffers {
			vk.DestroyFramebuffer(__logical_device, v, nil)
		}
		__framebuffers = nil
	}
}

func initCommandPool() error {
	fr_LInfo("Init command pool")
	queueFamilyIndex, err := getValidQueueFamilyIndex(__physical_device)
	if err != nil {
		return err
	}

	cmdPoolInfo := vk.CommandPoolCreateInfo{
		SType:            vk.StructureTypeCommandPoolCreateInfo,
		PNext:            nil,
		Flags:            vk.CommandPoolCreateFlags(vk.CommandPoolCreateResetCommandBufferBit),
		QueueFamilyIndex: queueFamilyIndex,
	}

	__command_pool = nil
	res := vk.CreateCommandPool(__logical_device, &cmdPoolInfo, nil, &__command_pool)
	err = HandleVKErrorResult(res, "failed to create command pool")
	if err != nil {
		return err
	}

	return nil
}
func deinitCommandPool() {
	if __command_pool != nil {
		fr_LInfo("Deinit command pool")
		vk.DestroyCommandPool(__logical_device, __command_pool, nil)
		__command_pool = nil
	}
}

func initCommandBuffers() error {
	fr_LInfo("Init command buffers")

	cmdBufInfo := vk.CommandBufferAllocateInfo{
		SType:              vk.StructureTypeCommandBufferAllocateInfo,
		PNext:              nil,
		CommandPool:        __command_pool,
		Level:              vk.CommandBufferLevelPrimary,
		CommandBufferCount: MAX_FRAMES_IN_FLIGHT,
	}

	__command_buffers = make([]vk.CommandBuffer, MAX_FRAMES_IN_FLIGHT)
	res := vk.AllocateCommandBuffers(__logical_device, &cmdBufInfo, __command_buffers)
	err := HandleVKErrorResult(res, "failed to allocate command buffers")
	if err != nil {
		return err
	}

	return nil
}
func deinitCommandBuffers() {
	if __command_buffers != nil {
		fr_LInfo("Deinit command buffers")
		__command_buffers = nil
	}
}

func recordCommandBuffer(cmdBuf vk.CommandBuffer, imageIndex uint32) error {
	beginInfo := vk.CommandBufferBeginInfo{
		SType:            vk.StructureTypeCommandBufferBeginInfo,
		PNext:            nil,
		Flags:            0,
		PInheritanceInfo: nil,
	}

	res := vk.BeginCommandBuffer(cmdBuf, &beginInfo)
	err := HandleVKErrorResult(res, "failed to begin command buffer")
	if err != nil {
		return err
	}

	rpBeginInfo := vk.RenderPassBeginInfo{
		SType:       vk.StructureTypeRenderPassBeginInfo,
		PNext:       nil,
		RenderPass:  __render_pass,
		Framebuffer: __framebuffers[imageIndex],
		RenderArea: vk.Rect2D{
			Offset: vk.Offset2D{X: 0, Y: 0},
			Extent: __swapchain_extent,
		},
		ClearValueCount: 1,
		PClearValues:    []vk.ClearValue{vk.ClearValue(NewClearColorFloat(0, 0, 0, 1))},
	}

	vk.CmdBeginRenderPass(cmdBuf, &rpBeginInfo, vk.SubpassContentsInline)
	vk.CmdBindPipeline(cmdBuf, vk.PipelineBindPointGraphics, __pipeline)

	viewport := vk.Viewport{
		X:        0,
		Y:        0,
		Width:    float32(__swapchain_extent.Width),
		Height:   float32(__swapchain_extent.Height),
		MinDepth: 0,
		MaxDepth: 1,
	}
	vk.CmdSetViewport(cmdBuf, 0, 1, []vk.Viewport{viewport})

	scissor := vk.Rect2D{
		Offset: vk.Offset2D{X: 0, Y: 0},
		Extent: __swapchain_extent,
	}
	vk.CmdSetScissor(cmdBuf, 0, 1, []vk.Rect2D{scissor})

	vk.CmdBindVertexBuffers(cmdBuf, 0, 1, []vk.Buffer{__vertex_buffer}, []vk.DeviceSize{0})

	vk.CmdDraw(cmdBuf, uint32(len(VERTICES)), 1, 0, 0)

	vk.CmdEndRenderPass(cmdBuf)
	res = vk.EndCommandBuffer(cmdBuf)
	err = HandleVKErrorResult(res, "failed to end (record) command buffer")
	if err != nil {
		return err
	}

	return nil
}

func initSyncObjects() error {
	fr_LInfo("Init synchronization objects")
	semaphoreInfo := vk.SemaphoreCreateInfo{
		SType: vk.StructureTypeSemaphoreCreateInfo,
		PNext: nil,
		Flags: 0,
	}

	fenceInfo := vk.FenceCreateInfo{
		SType: vk.StructureTypeFenceCreateInfo,
		PNext: nil,
		Flags: vk.FenceCreateFlags(vk.FenceCreateSignaledBit),
	}

	__sync_fence_inflights = make([]vk.Fence, MAX_FRAMES_IN_FLIGHT)
	__sync_semaphore_wait = make([]vk.Semaphore, MAX_FRAMES_IN_FLIGHT)
	__sync_semaphore_signal = make([]vk.Semaphore, len(__swapchain_images))

	for i := 0; i < MAX_FRAMES_IN_FLIGHT; i++ {
		res := vk.CreateSemaphore(__logical_device, &semaphoreInfo, nil, &__sync_semaphore_wait[i])
		err := HandleVKErrorResult(res, "failed to create image available semaphore")
		if err != nil {
			return err
		}

		res = vk.CreateFence(__logical_device, &fenceInfo, nil, &__sync_fence_inflights[i])
		err = HandleVKErrorResult(res, "failed to create fence")
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(__swapchain_images); i++ {
		res := vk.CreateSemaphore(__logical_device, &semaphoreInfo, nil, &__sync_semaphore_signal[i])
		err := HandleVKErrorResult(res, "failed to create render finished semaphore")
		if err != nil {
			return err
		}
	}

	return nil
}
func deinitSyncObjects() {
	if __sync_fence_inflights != nil || __sync_semaphore_signal != nil || __sync_semaphore_wait != nil {
		fr_LInfo("Deinit synchronization objects")
		for i := 0; i < MAX_FRAMES_IN_FLIGHT; i++ {
			vk.DestroyFence(__logical_device, __sync_fence_inflights[i], nil)
			vk.DestroySemaphore(__logical_device, __sync_semaphore_wait[i], nil)
		}
		for i := 0; i < len(__swapchain_images); i++ {
			vk.DestroySemaphore(__logical_device, __sync_semaphore_signal[i], nil)
		}
		__sync_fence_inflights = nil
		__sync_semaphore_wait = nil
		__sync_semaphore_signal = nil
	}
}

func reinitSwapchains() error {
	fr_LInfo("Reinit swapchain")
	flags := __window.GetFlags()
	for flags&sdl.WINDOW_MINIMIZED != 0 {
		flags = __window.GetFlags()
		sdl.PollEvent()
	}

	vk.DeviceWaitIdle(__logical_device)

	deinitFrameBuffers()
	deinitImageViews()
	deinitSwapchain()

	err := initSwapchain()
	if err != nil {
		return err
	}

	err = initImageViews()
	if err != nil {
		return err
	}

	err = initFrameBuffers()
	if err != nil {
		return err
	}

	return nil
}
