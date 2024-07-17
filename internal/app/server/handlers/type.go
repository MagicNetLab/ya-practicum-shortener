package handlers

type APIRequest struct {
	URL string `json:"url"`
}

type APIBatchRequest []APIBatchRequestEntity

type APIBatchRequestEntity struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type APIResponse struct {
	Result string `json:"result"`
}

type APIBatchResponse []APIBatchResponseEntity

type APIBatchResponseEntity struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
