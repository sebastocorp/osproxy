package global

import "osproxy/internal/utils"

const (
	HeaderContentType          = "Content-Type"
	HeaderContentTypeAppJson   = "application/json"
	HeaderContentTypeTextPlain = "text/plain"

	EndpointHealthz = "/healthz"
)

const (
	LogFieldKeyCommonService   = "service"
	LogFieldKeyCommonComponent = "component"

	LogFieldKeyExtraError = "error"

	LogFieldKeyExtraRequest   = "request"
	LogFieldKeyExtraRequestId = "request_id"
	LogFieldKeyExtraResponse  = "response"
	LogFieldKeyExtraReaction  = "reaction"
	LogFieldKeyExtraDataBytes = "bytes"

	LogFieldValueDefaultStr           = "none"
	LogFieldValueDefaultI64     int64 = 0
	LogFieldValueService              = "osproxy"
	LogFieldValueComponentProxy       = "Proxy"
)

func GetLogCommonFields() map[string]any {
	return map[string]any{
		LogFieldKeyCommonService:   LogFieldValueService,
		LogFieldKeyCommonComponent: LogFieldValueDefaultStr,
	}
}

func GetLogExtraFieldsProxy() map[string]any {
	return map[string]any{
		LogFieldKeyExtraRequestId: LogFieldValueDefaultStr,
		LogFieldKeyExtraRequest:   utils.DefaultRequestStruct(),
		LogFieldKeyExtraResponse:  utils.DefaultResponseStruct(),
		LogFieldKeyExtraReaction:  LogFieldValueDefaultStr,
		LogFieldKeyExtraDataBytes: LogFieldValueDefaultI64,
		LogFieldKeyExtraError:     LogFieldValueDefaultStr,
	}
}

func SetLogExtraField(extra map[string]any, key string, val any) {
	if _, ok := extra[key]; ok {
		extra[key] = val
	}
}
