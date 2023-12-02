package types

const (
	// '0' is sort articles by pushlish time
	SortPublishTime = iota
	// '1' is sort articles ny like count
	SortLikeCount
)

const (
	DefaultPageSize       = 20
	DefaultLimit          = 200
	DefaultSortLikeCursor = 1 << 30
)
