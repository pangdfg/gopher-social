package mailer

import "fmt"

type MockClient struct{}

func (m *MockClient) Send(to string, subject string, body string, data any, transactional bool) (int, error) {
    fmt.Printf("[MOCK EMAIL] To: %s, Subject: %s, Body: %s, Data: %+v, Transactional: %v\n",
        to, subject, body, data, transactional)
    return 202, nil
}
