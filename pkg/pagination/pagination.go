package pagination

// Calculate returns the limit and offset from the given page and pageSize arguments.
func Calculate(page, pageSize *int) (limit, offset int) {
	limit = 10
	offset = 0
	if pageSize != nil {
		limit = *pageSize
		if page != nil {
			offset = *page * *pageSize
		}
	}
	return
}
