package musly

type ErrorResponse struct {
    Message     string    `json:"message"`
    Pathname    string    `json:"pathname,omitempty"`
}
