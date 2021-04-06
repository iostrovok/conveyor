package faces

// IWorkBench is interface for support the storage for IItem.
type IWorkBench interface {
	// Set puts new IItem by number in WorkBench
	Add(item IItem) int
	// Get returns item by number in WorkBench
	Get(i int) (IItem, error)
	// Len returns the total length of WorkBench
	Len() int
	// Count returns the number of active IItem in WorkBench
	Count() int
	// Clean removes IItem from WorkBench (makes no-active)
	Clean(i int)
	// GetPriority returns the priority for item by number. If item is not fund, return 0.
	GetPriority(i int) int
}
