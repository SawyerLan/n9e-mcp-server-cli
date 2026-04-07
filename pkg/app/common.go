package app

// SlicePage performs client-side pagination on a slice.
// It returns the paginated slice and the total count of original items.
// Defaults: page=1, limit=20.
func SlicePage[T any](items []T, page, limit int) ([]T, int64) {
	total := int64(len(items))
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	start := (page - 1) * limit
	if int64(start) >= total {
		return []T{}, total
	}

	end := start + limit
	if int64(end) > total {
		end = int(total)
	}

	return items[start:end], total
}
