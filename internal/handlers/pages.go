package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"le-grimoire/internal/middleware"
	"le-grimoire/internal/render"
	"le-grimoire/internal/state"
)

func (h *Handler) CurrentPageHandler(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := middleware.GetDeviceID(r.Context())
	if !ok {
		h.RespondWithStatusError(w, http.StatusUnauthorized, fmt.Errorf("device ID not found in headers"))
		return
	}

	ds, err := h.App.DeviceRepository.GetDeviceState(deviceID)
	if err != nil || ds == nil {
		ds = state.NewDeviceState(deviceID)
	}

	// Render page based on current top of stack
	frame, err := h.renderCurrentState(r.Context(), ds)
	if err != nil {
		// Log error and fallback to pure stdlib 30,000-byte bitmap (never 500 to ESP32)
		frame = render.RenderFallbackErrorBitmap(err.Error())
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(frame)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(frame)
}

func (h *Handler) renderCurrentState(ctx context.Context, ds *state.DeviceState) ([]byte, error) {
	top := ds.Top()
	cursor := top.State["cursor"]

	switch top.Type {
	case state.PageLibrary:
		libraries, err := h.App.KavitaRepository.GetLibraries(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetch libraries: %w", err)
		}
		html, err := render.BuildLibraryHTML(libraries, cursor)
		if err != nil {
			return nil, err
		}
		return h.App.Renderer.RenderListPage(ctx, html)

	case state.PageSeries:
		libID, _ := strconv.Atoi(top.Params["library_id"])
		seriesList, err := h.App.KavitaRepository.GetSeries(ctx, libID)
		if err != nil {
			return nil, fmt.Errorf("fetch series: %w", err)
		}
		html, err := render.BuildSeriesHTML(seriesList, cursor)
		if err != nil {
			return nil, err
		}
		return h.App.Renderer.RenderListPage(ctx, html)

	case state.PageBookList:
		seriesID, _ := strconv.Atoi(top.Params["series_id"])
		chapters, err := h.App.KavitaRepository.GetFlattenedChapters(ctx, seriesID)
		if err != nil {
			return nil, fmt.Errorf("fetch chapters: %w", err)
		}
		html, err := render.BuildBookListHTML(chapters, cursor)
		if err != nil {
			return nil, err
		}
		return h.App.Renderer.RenderListPage(ctx, html)

	case state.PageReader:
		return h.renderReaderPage(ctx, top)

	default:
		html := render.BuildPlaceholderHTML(*top)
		return h.App.Renderer.RenderListPage(ctx, html)
	}
}

func (h *Handler) renderReaderPage(ctx context.Context, p *state.Page) ([]byte, error) {
	chapterID, _ := strconv.Atoi(p.Params["chapter_id"])
	format, _ := strconv.Atoi(p.Params["format"])
	subPageIndex := p.State["sub_page"]

	// Format 0: Manga / Comic
	if format == 0 {
		imgBytes, err := h.App.KavitaRepository.GetChapterPageImage(ctx, chapterID, subPageIndex)
		if err != nil {
			return nil, fmt.Errorf("fetch manga page %d: %w", subPageIndex, err)
		}
		return render.ProcessMangaImage(imgBytes)
	}

	// Format 1+: Book / EPUB
	// 1. Check if all frames for this chapter are already cached
	var frames [][]byte
	if cachedFrames, exists := h.App.FrameCache.GetAllFrames(chapterID); exists {
		frames = cachedFrames
	} else {
		// 2. Cache miss: Fetch and render full chapter into frames
		chapterHTML, err := h.App.KavitaRepository.GetBookPage(ctx, chapterID, 0)
		if err != nil {
			return nil, fmt.Errorf("fetch book content: %w", err)
		}

		cleanHTML := render.SanitizeEPUBHTML(chapterHTML, h.App.Config.GetKavitaAPIURI())
		renderedHTML, err := render.BuildReaderHTML(cleanHTML)
		if err != nil {
			return nil, fmt.Errorf("build reader html: %w", err)
		}

		frames, err = h.App.Renderer.RenderBookFrames(ctx, renderedHTML, 24)
		if err != nil {
			return nil, fmt.Errorf("render book frames: %w", err)
		}

		h.App.FrameCache.Set(chapterID, frames)
	}

	if len(frames) == 0 {
		return nil, fmt.Errorf("no rendered frames produced for chapter %d", chapterID)
	}

	// 3. Clamp index
	if subPageIndex >= len(frames) {
		subPageIndex = len(frames) - 1
		p.State["sub_page"] = subPageIndex
	}
	if subPageIndex < 0 {
		subPageIndex = 0
		p.State["sub_page"] = 0
	}

	return frames[subPageIndex], nil
}
