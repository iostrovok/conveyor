package faces

// File describes the trace interface.

const (
	// EmptySkipName is constant for skipped handlers.
	EmptySkipName Name = ""

	// SkipAll is constant to skip all rest handlers.
	SkipAll Name = "###SKIP_EVERYTHING_BY_THE_END"
)

type ITrace interface {
	// LazyPrintf evaluates its arguments with fmt.Sprintf each time the
	// /debug/requests page is rendered. Any memory referenced by a will be
	// pinned until the trace is finished and later discarded.
	LazyPrintf(format string, a ...interface{})

	// SetError declares that this trace resulted in an error.
	SetError()

	// Flush will call at the end on cycle.
	Flush()

	// Flush will call at the end on cycle.
	ForceFlush()
}
