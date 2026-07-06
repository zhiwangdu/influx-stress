package statement

// A Template contains all information to fill in templated variables in inset and query statements
type Template struct {
	Tags     []string
	Function *Function
}

// Templates are a collection of Template
type Templates []*Template

// Init makes Stringers out of the Templates for quick point creation
func (t Templates) Init(seriesCount int) Stringers {
	arr := make([]Stringer, len(t))
	for i, tmp := range t {
		arr[i] = stringerFromAppender(tmp.NewAppender(seriesCount))
	}
	return arr
}

// InitAppenders makes append-based value generators out of the templates.
func (t Templates) InitAppenders(seriesCount int) []ValueAppender {
	arr := make([]ValueAppender, len(t))
	for i, tmp := range t {
		arr[i] = tmp.NewAppender(seriesCount)
	}
	return arr
}

// NewAppender returns an append-based generator for a template.
func (t *Template) NewAppender(seriesCount int) ValueAppender {
	if len(t.Tags) == 0 {
		return t.Function.NewAppender(seriesCount)
	}
	return t.NewTagAppender()
}

// Calculates the number of series implied by a template
func (t *Template) numSeries() int {
	// If !t.Tags then tag cardinality is t.Function.Count
	if len(t.Tags) == 0 {
		return t.Function.Count
	}
	// Else tag cardinality is len(t.Tags)
	return len(t.Tags)
}

// NewTagFunc returns a Stringer that loops through the given tags
func (t *Template) NewTagFunc() Stringer {
	return stringerFromAppender(t.NewTagAppender())
}

// NewTagAppender returns a ValueAppender that loops through the given tags.
func (t *Template) NewTagAppender() ValueAppender {
	if len(t.Tags) == 0 {
		return errorAppender("EMPTY TAGS")
	}

	i := 0
	return func(dst []byte) []byte {
		s := t.Tags[i]
		i = (i + 1) % len(t.Tags)
		return append(dst, s...)
	}
}
