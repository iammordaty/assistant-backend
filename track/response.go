package track

type SuccessResponse struct {
    InitialKey  string    `json:"initial_key"`
    Bpm         float64   `json:"bpm"`
}

type ErrorResponse  struct {
    Message     string    `json:"message"`
    Pathname    string    `json:"pathname"`
}

type CommandErrorResponse  struct {
    ErrorResponse
    Command     string    `json:"command"`
}
