package helpers

import "fmt"
import (
	vk "github.com/vulkan-go/vulkan"
)

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
		return "Incomplete or ResultEndRange"
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
		return "ErrorFragmentedPool or ResultBeginRange"
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
	// case vk.ResultBeginRange:
	// 	return "ResultBeginRange"
	// case vk.ResultEndRange:
	// 	return "ResultEndRange"
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
		return fmt.Errorf("%s, vulkan error code (%d : %s)", msg, res, VKResultToString(res))
	}
	return nil
}
