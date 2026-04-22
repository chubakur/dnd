package transport

type queueMessage struct {
	Attributes             map[string]string `json:"attributes"`
	Body                   string            `json:"body"`
	Md5OfBody              string            `json:"md5_of_body"`
	Md5OfMessageAttributes string            `json:"md5_of_message_attributes"`
	// MessageAttributes      map[string]string `json:"message_attributes"`
	MessageId string `json:"message_id"`
}

type queueEventDetails struct {
	Message queueMessage `json:"message"`
	QueueId string       `json:"queue_id"`
}

type queueEventMetadata struct {
	CreatedAt string `json:"created_at"`
	EventId   string `json:"event_id"`
	EventType string `json:"event_type"`
}

type queueEvent struct {
	Details       queueEventDetails  `json:"details"`
	EventMetadata queueEventMetadata `json:"event_metadata"`
}

type queueRequest struct {
	Messages []queueEvent `json:"messages"`
}
