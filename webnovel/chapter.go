package webnovel

type Chapter struct {
	ID      int
	Title   string
	Content string
}

type Chapters struct {
	List []*Chapter
}

func (ch Chapters) Len() int {
	return len(ch.List)
}

func (ch Chapters) Less(i, j int) bool {
	return ch.List[j].ID > ch.List[i].ID
}

func (ch Chapters) Swap(i, j int) {
	ch.List[i], ch.List[j] = ch.List[j], ch.List[i]
}

func (ch *Chapters) AddChapter(chapter *Chapter) {
	ch.List = append(ch.List, chapter)
}