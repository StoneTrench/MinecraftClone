package rendering

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	vk "github.com/vulkan-go/vulkan"
)

func initRenderPass() error {
	fr_LInfo("Init render pass")
	attchDesc := vk.AttachmentDescription{
		Flags:          0,
		Format:         __swapchain_format.Format,
		Samples:        vk.SampleCount1Bit,
		LoadOp:         vk.AttachmentLoadOpClear,
		StoreOp:        vk.AttachmentStoreOpStore,
		StencilLoadOp:  vk.AttachmentLoadOpDontCare,
		StencilStoreOp: vk.AttachmentStoreOpDontCare,
		InitialLayout:  vk.ImageLayoutUndefined,
		FinalLayout:    vk.ImageLayoutPresentSrc,
	}
	attachments := []vk.AttachmentDescription{attchDesc}

	colorAttchRefs := []vk.AttachmentReference{
		{
			Attachment: 0,
			Layout:     vk.ImageLayoutColorAttachmentOptimal,
		},
	}
	subpass := vk.SubpassDescription{
		Flags:                   0,
		PipelineBindPoint:       vk.PipelineBindPointGraphics,
		InputAttachmentCount:    0,
		PInputAttachments:       nil,
		ColorAttachmentCount:    uint32(len(colorAttchRefs)),
		PColorAttachments:       colorAttchRefs,
		PResolveAttachments:     nil,
		PDepthStencilAttachment: nil,
		PreserveAttachmentCount: 0,
		PPreserveAttachments:    nil,
	}
	subpasses := []vk.SubpassDescription{subpass}

	depInfo := vk.SubpassDependency{
		SrcSubpass:      vk.SubpassExternal,
		DstSubpass:      0,
		SrcStageMask:    vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
		DstStageMask:    vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
		SrcAccessMask:   0,
		DstAccessMask:   vk.AccessFlags(vk.AccessColorAttachmentWriteBit),
		DependencyFlags: 0,
	}
	dependencies := []vk.SubpassDependency{depInfo}
	renderPassInfo := vk.RenderPassCreateInfo{
		SType:           vk.StructureTypeRenderPassCreateInfo,
		PNext:           nil,
		Flags:           0,
		AttachmentCount: uint32(len(attachments)),
		PAttachments:    attachments,
		SubpassCount:    uint32(len(subpasses)),
		PSubpasses:      subpasses,
		DependencyCount: uint32(len(dependencies)),
		PDependencies:   dependencies,
	}

	__render_pass = nil
	res := vk.CreateRenderPass(__logical_device, &renderPassInfo, nil, &__render_pass)
	return HandleVKErrorResult(res, "failed to create render pass")
}
func deinitRenderPass() {
	if __render_pass != nil {
		fr_LInfo("Deinit render pass")
		vk.DestroyRenderPass(__logical_device, __render_pass, nil)
		__render_pass = nil
	}
}

func createShaderModule(name string) (vk.ShaderModule, error) {
	codeBytes, err := readShaderFile(name)
	if err != nil {
		return nil, err
	}

	if len(codeBytes)%4 != 0 {
		return nil, fmt.Errorf("shader '%s' size (%d) is not a multiple of 4 bytes (SPIR-V bytecode is divisible by 4)", name, len(codeBytes))
	}

	byteCount := uint(len(codeBytes))
	wordCount := byteCount / 4
	codeUint32 := make([]uint32, wordCount)

	reader := bytes.NewReader(codeBytes)
	err = binary.Read(reader, binary.NativeEndian, &codeUint32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode shader bytes: %w", err)
	}

	fr_LInfo(fmt.Sprintf("Shader size: %d bytes (%d words)", len(codeBytes), wordCount))

	smInfo := vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		PNext:    nil,
		Flags:    0,
		CodeSize: byteCount,
		PCode:    codeUint32,
	}

	var module vk.ShaderModule
	res := vk.CreateShaderModule(__logical_device, &smInfo, nil, &module)
	err = HandleVKErrorResult(res, "failed to create shader module")
	if err != nil {
		return nil, err
	}

	return module, nil
}
func readShaderFile(name string) ([]byte, error) {
	result, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read shader, %w", err)
	}
	return result, nil
}

func initPipeline_ShaderStages() ([]vk.PipelineShaderStageCreateInfo, error) {

	// Shader Modules
	module, err := createShaderModule("assets/shaders/shader.spv")
	if err != nil {
		return nil, err
	}
	__pipeline_shader_module = module

	shaderVertStageInfo := vk.PipelineShaderStageCreateInfo{
		SType:               vk.StructureTypePipelineShaderStageCreateInfo,
		PNext:               nil,
		Flags:               0,
		Stage:               vk.ShaderStageVertexBit,
		Module:              module,
		PName:               "vertMain\x00",
		PSpecializationInfo: nil,
	}

	shaderFragStageInfo := vk.PipelineShaderStageCreateInfo{
		SType:               vk.StructureTypePipelineShaderStageCreateInfo,
		PNext:               nil,
		Flags:               0,
		Stage:               vk.ShaderStageFragmentBit,
		Module:              module,
		PName:               "fragMain\x00",
		PSpecializationInfo: nil,
	}

	shaderStages := []vk.PipelineShaderStageCreateInfo{
		shaderVertStageInfo,
		shaderFragStageInfo,
	}

	return shaderStages, nil
}

func initPipeline_PipelineLayout() error {
	pipelineLayoutInfo := vk.PipelineLayoutCreateInfo{
		SType:                  vk.StructureTypePipelineLayoutCreateInfo,
		PNext:                  nil,
		Flags:                  0,
		SetLayoutCount:         0,
		PSetLayouts:            nil,
		PushConstantRangeCount: 0,
		PPushConstantRanges:    nil,
	}

	__pipeline_layout = nil
	res := vk.CreatePipelineLayout(__logical_device, &pipelineLayoutInfo, nil, &__pipeline_layout)
	return HandleVKErrorResult(res, "failed to create pipeline layout")
}

func initPipeline_VertexInputState() vk.PipelineVertexInputStateCreateInfo {
	vertexInputStateInfo := vk.PipelineVertexInputStateCreateInfo{
		SType:                           vk.StructureTypePipelineVertexInputStateCreateInfo,
		PNext:                           nil,
		Flags:                           0,
		VertexBindingDescriptionCount:   0,
		PVertexBindingDescriptions:      nil,
		VertexAttributeDescriptionCount: 0,
		PVertexAttributeDescriptions:    nil,
	}
	return vertexInputStateInfo
}

func initPipeline_ViewportState() vk.PipelineViewportStateCreateInfo {
	viewport := vk.Viewport{
		X:        0,
		Y:        0,
		Width:    float32(__swapchain_extent.Width),
		Height:   float32(__swapchain_extent.Height),
		MinDepth: 0,
		MaxDepth: 1,
	}

	scissor := vk.Rect2D{
		Offset: vk.Offset2D{X: 0, Y: 0},
		Extent: __swapchain_extent,
	}

	viewportStateInfo := vk.PipelineViewportStateCreateInfo{
		SType:         vk.StructureTypePipelineViewportStateCreateInfo,
		PNext:         nil,
		Flags:         0,
		ViewportCount: 1,
		PViewports:    []vk.Viewport{viewport},
		ScissorCount:  1,
		PScissors:     []vk.Rect2D{scissor},
	}

	return viewportStateInfo
}

func initPipeline_DynamicState() vk.PipelineDynamicStateCreateInfo {
	dynamicStates := []vk.DynamicState{
		vk.DynamicStateViewport,
		vk.DynamicStateScissor,
	}

	dynamicStateInfo := vk.PipelineDynamicStateCreateInfo{
		SType:             vk.StructureTypePipelineDynamicStateCreateInfo,
		PNext:             nil,
		Flags:             0,
		DynamicStateCount: uint32(len(dynamicStates)),
		PDynamicStates:    dynamicStates,
	}

	return dynamicStateInfo
}

func initPipeline_RasterizerState() vk.PipelineRasterizationStateCreateInfo {
	rasterizerStateInfo := vk.PipelineRasterizationStateCreateInfo{
		SType:                   vk.StructureTypePipelineRasterizationStateCreateInfo,
		PNext:                   nil,
		Flags:                   0,
		DepthClampEnable:        vk.False,
		RasterizerDiscardEnable: vk.False,
		PolygonMode:             vk.PolygonModeFill,
		CullMode:                vk.CullModeFlags(vk.CullModeBackBit),
		FrontFace:               vk.FrontFaceClockwise,
		LineWidth:               1,

		DepthBiasEnable:         vk.False,
		DepthBiasConstantFactor: 0,
		DepthBiasClamp:          0,
		DepthBiasSlopeFactor:    0,
	}

	return rasterizerStateInfo
}

func initPipeline_MultisampleState() vk.PipelineMultisampleStateCreateInfo {
	multisampleStateInfo := vk.PipelineMultisampleStateCreateInfo{
		SType:                 vk.StructureTypePipelineMultisampleStateCreateInfo,
		PNext:                 nil,
		Flags:                 0,
		RasterizationSamples:  vk.SampleCount1Bit,
		SampleShadingEnable:   vk.False,
		MinSampleShading:      0,
		PSampleMask:           nil,
		AlphaToCoverageEnable: vk.False,
		AlphaToOneEnable:      vk.False,
	}

	return multisampleStateInfo
}

func initPipeline_StencilState() vk.PipelineDepthStencilStateCreateInfo {
	stencilStateInfo := vk.PipelineDepthStencilStateCreateInfo{
		SType:                 vk.StructureTypePipelineDepthStencilStateCreateInfo,
		PNext:                 nil,
		Flags:                 0,
		DepthTestEnable:       vk.False,
		DepthWriteEnable:      vk.False,
		DepthCompareOp:        vk.CompareOpLessOrEqual,
		DepthBoundsTestEnable: vk.False,
		StencilTestEnable:     vk.False,
		Front:                 vk.StencilOpState{},
		Back:                  vk.StencilOpState{},
		MinDepthBounds:        0,
		MaxDepthBounds:        1,
	}

	return stencilStateInfo
}

func initPipeline_ColorBlendState() vk.PipelineColorBlendStateCreateInfo {

	// colorBlendAttachments := vk.PipelineColorBlendAttachmentState{
	// 	BlendEnable:         vk.True,
	// 	SrcColorBlendFactor: vk.BlendFactorSrcAlpha,
	// 	DstColorBlendFactor: vk.BlendFactorDstAlpha,
	// 	ColorBlendOp:        vk.BlendOpAdd,
	// 	SrcAlphaBlendFactor: vk.BlendFactorOne,
	// 	DstAlphaBlendFactor: vk.BlendFactorZero,
	// 	AlphaBlendOp:        vk.BlendOpAdd,
	// 	ColorWriteMask: vk.ColorComponentFlags(vk.ColorComponentRBit) |
	// 		vk.ColorComponentFlags(vk.ColorComponentGBit) |
	// 		vk.ColorComponentFlags(vk.ColorComponentBBit) |
	// 		vk.ColorComponentFlags(vk.ColorComponentABit),
	// }
	colorBlendAttachments := vk.PipelineColorBlendAttachmentState{
		BlendEnable: vk.False,
		ColorWriteMask: vk.ColorComponentFlags(vk.ColorComponentRBit) |
			vk.ColorComponentFlags(vk.ColorComponentGBit) |
			vk.ColorComponentFlags(vk.ColorComponentBBit) |
			vk.ColorComponentFlags(vk.ColorComponentABit),
	}

	colorBlendState := vk.PipelineColorBlendStateCreateInfo{
		SType:           vk.StructureTypePipelineColorBlendStateCreateInfo,
		PNext:           nil,
		Flags:           0,
		LogicOpEnable:   vk.False,
		LogicOp:         vk.LogicOpCopy,
		AttachmentCount: 1,
		PAttachments:    []vk.PipelineColorBlendAttachmentState{colorBlendAttachments},
		BlendConstants:  [4]float32{0, 0, 0, 0},
	}

	return colorBlendState
}

func initPipeline_InputAssemblyState() vk.PipelineInputAssemblyStateCreateInfo {
	inputAssemblyInfo := vk.PipelineInputAssemblyStateCreateInfo{
		SType:                  vk.StructureTypePipelineInputAssemblyStateCreateInfo,
		PNext:                  nil,
		Flags:                  0,
		Topology:               vk.PrimitiveTopologyTriangleFan,
		PrimitiveRestartEnable: vk.False,
	}

	return inputAssemblyInfo
}

func initPipeline_PipelineCache() error {
	cacheInfo := vk.PipelineCacheCreateInfo{
		SType:           vk.StructureTypePipelineCacheCreateInfo,
		PNext:           nil,
		Flags:           0,
		InitialDataSize: 0,
		PInitialData:    nil,
	}

	__pipeline_cache = nil
	res := vk.CreatePipelineCache(__logical_device, &cacheInfo, nil, &__pipeline_cache)
	return HandleVKErrorResult(res, "failed to create pipeline cache")
}

func initPipeline() (err error) {
	fr_LInfo("Init shader module")

	shaderStages, err := initPipeline_ShaderStages()
	if err != nil {
		return err
	}

	err = initPipeline_PipelineLayout()
	if err != nil {
		return err
	}

	err = initPipeline_PipelineCache()
	if err != nil {
		return err
	}

	inputAssemblyStateInfo := initPipeline_InputAssemblyState()
	vertexInputStateInfo := initPipeline_VertexInputState()
	viewportStateInfo := initPipeline_ViewportState()
	rasterizerStateInfo := initPipeline_RasterizerState()
	multisampleStateInfo := initPipeline_MultisampleState()
	stencilStateInfo := initPipeline_StencilState()
	colorBlendStateInfo := initPipeline_ColorBlendState()

	// fr_LInfo(fmt.Sprintf("PViewports nil? %t", viewportStateInfo.PViewports == nil))
	// fr_LInfo(fmt.Sprintf("PScissors nil? %t", viewportStateInfo.PScissors == nil))

	dynamicStateInfo := initPipeline_DynamicState()

	// Final
	graphicsPipelineInfo := vk.GraphicsPipelineCreateInfo{
		SType:               vk.StructureTypeGraphicsPipelineCreateInfo,
		PNext:               nil,
		Flags:               0,
		StageCount:          uint32(len(shaderStages)),
		PStages:             shaderStages,
		PVertexInputState:   &vertexInputStateInfo,
		PInputAssemblyState: &inputAssemblyStateInfo,
		PTessellationState:  nil,
		PViewportState:      &viewportStateInfo,
		PRasterizationState: &rasterizerStateInfo,
		PMultisampleState:   &multisampleStateInfo,
		PDepthStencilState:  &stencilStateInfo,
		PColorBlendState:    &colorBlendStateInfo,
		PDynamicState:       &dynamicStateInfo,
		Layout:              __pipeline_layout,
		RenderPass:          __render_pass,
		Subpass:             0,
		BasePipelineHandle:  nil,
		BasePipelineIndex:   -1,
	}

	pipelines := make([]vk.Pipeline, 1)
	fr_LInfo(fmt.Sprint(viewportStateInfo))
	fr_LInfo(fmt.Sprint(__logical_device, __pipeline_cache, 1, []vk.GraphicsPipelineCreateInfo{graphicsPipelineInfo}, nil, pipelines))
	res := vk.CreateGraphicsPipelines(__logical_device, __pipeline_cache, 1, []vk.GraphicsPipelineCreateInfo{graphicsPipelineInfo}, nil, pipelines)
	err = HandleVKErrorResult(res, "failed to create graphics pipelines")
	if err != nil {
		return err
	}
	__pipeline = pipelines[0]

	return nil
}
func deinitPipeline() {
	if __pipeline != nil {
		fr_LInfo("Deinit pipeline")
		vk.DestroyPipeline(__logical_device, __pipeline, nil)
		__pipeline = nil
	}
	if __pipeline_cache != nil {
		fr_LInfo("Deinit pipeline cache")
		vk.DestroyPipelineCache(__logical_device, __pipeline_cache, nil)
		__pipeline_cache = nil
	}
	if __pipeline_layout != nil {
		fr_LInfo("Deinit pipeline layout")
		vk.DestroyPipelineLayout(__logical_device, __pipeline_layout, nil)
		__pipeline_layout = nil
	}
	if __pipeline_shader_module != nil {
		fr_LInfo("Deinit shader module")
		vk.DestroyShaderModule(__logical_device, __pipeline_shader_module, nil)
		__pipeline_shader_module = nil
	}
}
