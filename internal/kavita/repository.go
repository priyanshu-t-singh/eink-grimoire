package kavita

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

type Repository interface {
	Authenticate(ctx context.Context) error
	GetLibraries(ctx context.Context) ([]Library, error)
	GetSeries(ctx context.Context, libraryID int) ([]Series, error)
	GetSeriesDetail(ctx context.Context, seriesID int) (*Series, error)
	GetFlattenedChapters(ctx context.Context, seriesID int) ([]Chapter, error)
	GetChapterMetadata(ctx context.Context, chapterID int) (*ChapterInfo, error)
	GetBookPage(ctx context.Context, chapterID int, page int) (string, error)
	GetChapterPageImage(ctx context.Context, chapterID int, page int) ([]byte, error)
}

type repository struct {
	client *Client
}

func NewRepository(client *Client) Repository {
	return &repository{client}
}

func (r *repository) Authenticate(ctx context.Context) error {
	return r.client.Authenticate(ctx)
}

func (r *repository) GetLibraries(ctx context.Context) ([]Library, error) {
	var libraries []Library
	if err := r.client.getJSON(ctx, "/api/Library/libraries", &libraries); err != nil {
		return nil, fmt.Errorf("get libraries: %w", err)
	}
	return libraries, nil
}

func (r *repository) GetSeries(ctx context.Context, libraryID int) ([]Series, error) {
	reqBody := SeriesFilterRequest{
		Statements: []FilterStatement{
			{
				Field:      19, // LibraryId
				Value:      strconv.Itoa(libraryID),
				Comparison: 0, // Equal
			},
		},
		Combination: 1,
		LimitTo:     0,
		SortOptions: SortOptions{
			IsAscending: true,
			SortField:   1, // Name
		},
	}

	var series []Series
	if err := r.client.postJSON(ctx, "/api/Series/v2", reqBody, &series); err != nil {
		return nil, fmt.Errorf("get series for library %d: %w", libraryID, err)
	}
	return series, nil
}

func (r *repository) GetSeriesDetail(ctx context.Context, seriesID int) (*Series, error) {
	var s Series
	path := fmt.Sprintf("/api/Series/%d", seriesID)
	if err := r.client.getJSON(ctx, path, &s); err != nil {
		return nil, fmt.Errorf("get series detail %d: %w", seriesID, err)
	}
	return &s, nil
}

func (r *repository) GetFlattenedChapters(ctx context.Context, seriesID int) ([]Chapter, error) {
	var volumes []Volume
	path := fmt.Sprintf("/api/Series/volumes?seriesId=%d", seriesID)
	if err := r.client.getJSON(ctx, path, &volumes); err != nil {
		return nil, fmt.Errorf("get volumes for series %d: %w", seriesID, err)
	}

	var chapters []Chapter
	for _, vol := range volumes {
		for _, ch := range vol.Chapters {
			// Ensure VolumeID is populated on the chapter object
			if ch.VolumeID == 0 {
				ch.VolumeID = vol.ID
			}
			chapters = append(chapters, ch)
		}
	}

	// Keep natural chapter ordering
	sort.SliceStable(chapters, func(i, j int) bool {
		numI, errI := strconv.ParseFloat(chapters[i].Number, 64)
		numJ, errJ := strconv.ParseFloat(chapters[j].Number, 64)
		if errI == nil && errJ == nil {
			return numI < numJ
		}
		return chapters[i].Number < chapters[j].Number
	})

	return chapters, nil
}

func (r *repository) GetChapterMetadata(ctx context.Context, chapterID int) (*ChapterInfo, error) {
	var info ChapterInfo
	path := fmt.Sprintf("/api/Reader/chapter-info?chapterId=%d", chapterID)
	if err := r.client.getJSON(ctx, path, &info); err != nil {
		return nil, fmt.Errorf("get chapter info for %d: %w", chapterID, err)
	}
	return &info, nil
}

// fetches the HTML/text content for EPUB/book readers.
func (r *repository) GetBookPage(ctx context.Context, chapterID int, page int) (string, error) {
	path := fmt.Sprintf("/api/book/%d/book-page?page=%d", chapterID, page)
	text, err := r.client.getText(ctx, path)
	if err != nil {
		return "", fmt.Errorf("get book page (chapter %d, page %d): %w", chapterID, page, err)
	}
	return text, nil
}

// downloads the raw page bitmap/image for image-based manga/comic readers.
func (r *repository) GetChapterPageImage(ctx context.Context, chapterID int, page int) ([]byte, error) {
	apiKey := r.client.GetUserAPIKey()
	path := fmt.Sprintf(
		"/api/Reader/image?chapterId=%d&page=%d&apiKey=%s&extractPdf=false",
		chapterID,
		page,
		url.QueryEscape(apiKey),
	)

	imgBytes, err := r.client.getBytes(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("download page image (chapter %d, page %d): %w", chapterID, page, err)
	}
	return imgBytes, nil
}
