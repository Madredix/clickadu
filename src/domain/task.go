package domain

type (
	Tasks []Task
	Task  []*TaskItem

	TaskItem struct {
		NumberOfRequests int16  `json:"number_of_requests,omitempty"`
		URL              string `json:"url,omitempty"`
	}
)
