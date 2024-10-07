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

	LogFieldKeyExtraError  = "error"
	LogFieldKeyExtraObject = "object"

	LogFieldKeyExtraActionType = "action_type"

	LogFieldKeyExtraRequest       = "request"
	LogFieldKeyExtraStatusCode    = "status_code"
	LogFieldKeyExtraContentLength = "content_length"
	LogFieldKeyExtraDataBytes     = "data_bytes"

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
		LogFieldKeyExtraError:  LogFieldValueDefault,
		LogFieldKeyExtraObject: LogFieldValueDefault,
	}
}

func GetLogExtraFieldsActionWorker() map[string]any {
	return map[string]any{
		LogFieldKeyExtraError:  LogFieldValueDefault,
		LogFieldKeyExtraObject: LogFieldValueDefault,
	}
}
