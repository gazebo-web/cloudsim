package queues

// SwapResponse represents a response from the Swap method.
type SwapResponse struct {
	a QueueItemResponse
	b QueueItemResponse
}
