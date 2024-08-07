package messaging

import "time"

type ConsumerConfig struct {
	GroupName      string          `json:"durable_name,omitempty"`
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description,omitempty"`
	AckWait        time.Duration   `json:"ack_wait,omitempty"`
	MaxDeliver     int             `json:"max_deliver,omitempty"`
	FilterSubject  string          `json:"filter_subject,omitempty"`
	FilterSubjects []string        `json:"filter_subjects,omitempty"`
	RateLimit      uint64          `json:"rate_limit_bps,omitempty"` // Bits per sec
	MaxWaiting     int             `json:"max_waiting,omitempty"`
	MaxAckPending  int             `json:"max_ack_pending,omitempty"`
	HeadersOnly    bool            `json:"headers_only,omitempty"`
	BackOff        []time.Duration `json:"backoff,omitempty"`

	// Pull based options.
	MaxRequestBatch    int           `json:"max_batch,omitempty"`
	MaxRequestExpires  time.Duration `json:"max_expires,omitempty"`
	MaxRequestMaxBytes int           `json:"max_bytes,omitempty"`

	// Push based consumers.
	DeliverSubject string `json:"deliver_subject,omitempty"`
	DeliverGroup   string `json:"deliver_group,omitempty"`

	// Inactivity threshold.
	InactiveThreshold time.Duration `json:"inactive_threshold,omitempty"`

	// Generally inherited by parent stream and other markers, now can be configured directly.
	Replicas int `json:"num_replicas"`
}
