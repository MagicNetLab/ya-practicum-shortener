package handlers

// APIRequest структура запроса для создания короткой ссылки
type APIRequest struct {
	URL string `json:"url"`
}

// APIBatchRequest структура запроса для пакетного сокращения ссылок
type APIBatchRequest []APIBatchRequestEntity

// APIBatchRequestEntity структура данных ссылки для пакетного сокращения
type APIBatchRequestEntity struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// APIResponse структура ответа после сокращения ссылки
type APIResponse struct {
	Result string `json:"result"`
}

// APIBatchResponse структура ответа при пакетном сокращении ссылок
type APIBatchResponse []APIBatchResponseEntity

// APIBatchResponseEntity структура ссылки в ответе при пакетном сокращении
type APIBatchResponseEntity struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserLinkEntity структура ссылки в ответе со списком всех ссылок пользователя
type UserLinkEntity struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// UserLinksResponse структура ответа при запросе всех ссылок пользователя
type UserLinksResponse []UserLinkEntity

// APIDeleteRequest структура запроса при пакетном удалении ссылок пользователя
type APIDeleteRequest []string
