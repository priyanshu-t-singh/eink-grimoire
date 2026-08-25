package kavita

// JSON payload from /api/Plugin/authenticate.
type AuthResponse struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	ApiKey        string `json:"apiKey"`
	KavitaVersion string `json:"kavitaVersion"`
}

// Library represents a Kavita library.
type Library struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  int    `json:"type"` // 0: Manga, 1: Comic, 2: Book, etc.
	Cover string `json:"coverImage"`
}

type Series struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"originalName"`
	LibraryID    int    `json:"libraryId"`
	Format       int    `json:"format"` // 0: Image/Manga, 1: EPUB/Book, 2: PDF
	Pages        int    `json:"pages"`
}

// the payload for POST /api/Series/v2.
type SeriesFilterRequest struct {
	Statements  []FilterStatement `json:"statements"`
	Combination int               `json:"combination"`
	LimitTo     int               `json:"limitTo"`
	SortOptions SortOptions       `json:"sortOptions"`
}

// FilterStatement represents a condition in SeriesFilterRequest.
type FilterStatement struct {
	Field      int    `json:"field"`      // 19 = LibraryId
	Value      string `json:"value"`      // library ID as string
	Comparison int    `json:"comparison"` // 0 = Equal
}

// SortOptions represents ordering options in SeriesFilterRequest.
type SortOptions struct {
	IsAscending bool `json:"isAscending"`
	SortField   int  `json:"sortField"` // 1 = Name/Title
}

type Volume struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Pages    int       `json:"pages"`
	SeriesID int       `json:"seriesId"`
	Number   int       `json:"number"`
	Chapters []Chapter `json:"chapters"`
}

// Chapter represents a single chapter or loose issue.
type Chapter struct {
	ID          int    `json:"id"`
	VolumeID    int    `json:"volumeId"`
	Title       string `json:"title"`
	Number      string `json:"number"`
	Pages       int    `json:"pages"`
	IsSpecial   bool   `json:"isSpecial"`
	ReleaseDate string `json:"releaseDate"`
}

// /api/Reader/chapter-info.
type ChapterInfo struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Pages      int    `json:"pages"`
	TotalPages int    `json:"totalPages"`
	ChapterNum string `json:"chapterNumber"`
	VolumeID   int    `json:"volumeId"`
}
