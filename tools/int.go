package tools

// Intptr returns an int pointer from the given int.
func Intptr(i int) *int {
	return &i
}

// Int32ptr returns an int32 pointer from the given int32.
func Int32ptr(i int32) *int32 {
	return &i
}

// Int64ptr returns an int64 pointer from the given int64.
func Int64ptr(i int64) *int64 {
	return &i
}
