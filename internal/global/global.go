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

	LogFieldKeyExtraRequest    = "request"
	LogFieldKeyExtraResponse   = "response"
	LogFieldKeyExtraStatusCode = "code"
	LogFieldKeyExtraDataBytes  = "bytes"

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
		LogFieldKeyExtraError: LogFieldValueDefault,
	}
}

func GetLogExtraFieldsActionWorker() map[string]any {
	return map[string]any{
		LogFieldKeyExtraError: LogFieldValueDefault,
	}
}
