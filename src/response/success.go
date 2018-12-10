package response

import "fmt"

const (
	StatusOk Success = "OK"
)

type Success string

func (s *Success) String() string {
	return fmt.Sprintf("{\"status\":\"%v\"}", s)
}

func (s *Success) MarshalJSON() ([]byte, error) {
	return []byte(s.String()), nil
}
