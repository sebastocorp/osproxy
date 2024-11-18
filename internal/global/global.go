package global

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
		LogFieldKeyExtraRequest:   LogFieldValueDefaultStr,
		LogFieldKeyExtraRequestId: LogFieldValueDefaultStr,
		LogFieldKeyExtraResponse:  LogFieldValueDefaultStr,
		LogFieldKeyExtraReaction:  LogFieldValueDefaultStr,
		LogFieldKeyExtraDataBytes: LogFieldValueDefaultI64,
		LogFieldKeyExtraError:     LogFieldValueDefaultStr,
	}
}

func ResetLogExtraFields(extra map[string]any) {
	for k := range extra {
		extra[k] = LogFieldValueDefaultStr
	}
}

func ResetLogExtraField(extra map[string]any, key string) {
	if _, ok := extra[key]; ok {
		extra[key] = LogFieldValueDefaultStr
	}
}

func SetLogExtraField(extra map[string]any, key string, val any) {
	if _, ok := extra[key]; ok {
		extra[key] = val
	}
}
