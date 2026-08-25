package state

// ApplyButton mutates ds.Stack in place per the button table in the
// architecture doc and returns a human-readable description of what
// happened. Content fetch/render is intentionally out of scope here —
// Params on pushed pages are stubbed with placeholder IDs.

func ApplyButton(ds *DeviceState, buttonID, pressType string) string {
	p := ds.Top()

	switch p.Type {
	case PageLibrary, PageSeries, PageBookList:
		return applyListButton(ds, p, buttonID, pressType)
	case PageReader:
		return applyReaderButton(ds, p, buttonID, pressType)
	default:
		return "unknown page type, no-op"
	}
}

func applyListButton(ds *DeviceState, p *Page, buttonID, pressType string) string {
	if pressType == "long_press" {
		switch buttonID {
		case "E":
			return "force refresh (cache bypass) — no state change"
		default:
			return buttonID + " long-press: unassigned, no-op"
		}
	}

	switch buttonID {
	case "A": // Up
		if p.State["cursor"] > 0 {
			p.State["cursor"]--
		}
		return "cursor up"
	case "B": // Down
		p.State["cursor"]++ // NOTE: real impl clamps against fetched list length
		return "cursor down"
	case "C": // Select
		return selectFromList(ds, p)
	case "D": // Back
		if len(ds.Stack) == 1 {
			return "back at root — no-op"
		}
		ds.Pop()
		return "back — popped page"
	default:
		return buttonID + " short-press: unassigned, no-op"
	}
}

func selectFromList(ds *DeviceState, p *Page) string {
	switch p.Type {
	case PageLibrary:
		libID := stubID("lib", p.State["cursor"])
		ds.Push(Page{
			Type:   PageSeries,
			Params: map[string]string{"library_id": libID},
			State:  map[string]int{"cursor": 0, "scroll": 0},
		})
		return "selected library " + libID + " -> pushed Series"
	case PageSeries:
		seriesID := stubID("series", p.State["cursor"])
		ds.Push(Page{
			Type:   PageBookList,
			Params: map[string]string{"series_id": seriesID},
			State:  map[string]int{"cursor": 0, "scroll": 0},
		})
		return "selected series " + seriesID + " -> pushed BookList"
	case PageBookList:
		chapterID := stubID("chapter", p.State["cursor"])
		ds.Push(Page{
			Type:   PageReader,
			Params: map[string]string{"chapter_id": chapterID},
			State:  map[string]int{"sub_page": 0},
		})
		return "selected chapter " + chapterID + " -> pushed Reader"
	}
	return "select: nothing to do"
}

func applyReaderButton(ds *DeviceState, p *Page, buttonID, pressType string) string {
	if pressType == "long_press" {
		switch buttonID {
		case "E":
			return "force refresh (cache bypass) — no state change"
		default:
			return buttonID + " long-press: unassigned, no-op"
		}
	}

	switch buttonID {
	case "A": // Previous sub-page
		if p.State["sub_page"] > 0 {
			p.State["sub_page"]--
		}
		return "previous sub-page"
	case "B": // Next sub-page
		p.State["sub_page"]++ // NOTE: real impl clamps against cached frameCount
		return "next sub-page"
	case "C": // Back to BookList
		ds.Pop()
		return "back — popped to BookList"
	case "D": // Previous chapter
		return "previous chapter: stub — no-op at first chapter" // open item: chapter-boundary
	case "E": // Next chapter (short)
		return "next chapter: stub — no-op at last chapter" // open item: chapter-boundary
	default:
		return buttonID + " short-press: unassigned, no-op"
	}
}

func stubID(kind string, cursor int) string {
	return kind + "-" + itoa(cursor)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
