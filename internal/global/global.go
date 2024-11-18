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

	LogFieldKeyExtraActionType = "action_type"

	LogFieldKeyExtraRequest   = "request"
	LogFieldKeyExtraResponse  = "response"
	LogFieldKeyExtraReaction  = "reaction"
	LogFieldKeyExtraDataBytes = "bytes"

	LogFieldValueDefault               = "none"
	LogFieldValueService               = "osproxy"
	LogFieldValueComponentProxy        = "Proxy"
	LogFieldValueComponentActionWorker = "ActionWorker"
)

func GetLogCommonFields() map[string]any {
	return map[string]any{
		LogFieldKeyCommonService:   LogFieldValueService,
		LogFieldKeyCommonComponent: LogFieldValueDefault,
	}
}

func GetLogExtraFieldsProxy() map[string]any {
	return map[string]any{
		LogFieldKeyExtraError:     LogFieldValueDefault,
		LogFieldKeyExtraRequest:   LogFieldValueDefault,
		LogFieldKeyExtraResponse:  LogFieldValueDefault,
		LogFieldKeyExtraDataBytes: LogFieldValueDefault,
	}
}

func GetLogExtraFieldsActionWorker() map[string]any {
	return map[string]any{
		LogFieldKeyExtraError: LogFieldValueDefault,
	}
}

func ResetLogExtraFields(extra map[string]any) {
	for k := range extra {
		extra[k] = LogFieldValueDefault
	}
}

func ResetLogExtraField(extra map[string]any, key string) {
	if _, ok := extra[key]; ok {
		extra[key] = LogFieldValueDefault
	}
}

func SetLogExtraField(extra map[string]any, key string, val any) {
	if _, ok := extra[key]; ok {
		extra[key] = val
	}
}
