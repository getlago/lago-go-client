package lago

type EventRequest struct {
	client *Client
}

type EventParams struct {
	Event *EventInput `json:"event"`
}

type EventInput struct {
	TransactionID string            `json:"transaction_id,omitempty"`
	ExternalCustomerID    string    `json:"external_customer_id,omitempty"`
	Code          string            `json:"code,omitempty"`
	Timestamp     int64             `json:"timestamp,omitempty"`
	Properties    map[string]string `json:"properties,omitempty"`
}

func (c *Client) Event() *EventRequest {
	return &EventRequest{
		client: c,
	}
}

func (er *EventRequest) Create(eventInput *EventInput) *Error {
	eventParams := &EventParams{
		Event: eventInput,
	}

	clientRequest := &ClientRequest{
		Path: "events",
		Body: eventParams,
	}

	err := er.client.PostWithoutResult(clientRequest)
	if err != nil {
		return err
	}

	return nil
}
