package send_model

// ExcelUploadRequest represents the multipart form data for uploading spreadsheet files
type ExcelUploadRequest struct {
	Text     string `form:"text"`
	MediaUrl string `form:"mediaUrl"`
	Delay    int32  `form:"delay"`
}

// BulkSendRowResult represents the result of sending a message to a single row in the spreadsheet
type BulkSendRowResult struct {
	Row    int    `json:"row"`
	Number string `json:"number"`
	Name   string `json:"name,omitempty"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// BulkSendSummary represents the complete report returned after processing all rows in the spreadsheet
type BulkSendSummary struct {
	Status  string              `json:"status"`
	Total   int                 `json:"total"`
	Sent    int                 `json:"sent"`
	Failed  int                 `json:"failed"`
	Results []BulkSendRowResult `json:"results"`
}
