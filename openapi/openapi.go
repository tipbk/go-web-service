package openapi

type OpenApiError struct {
	ErrorCode   string `json:"errorCode"`
	ErrorReason string `json:"errorReason"`
}
