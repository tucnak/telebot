package telebot

// Filter is some thing that does filtering for
// incoming updates.
//
// Return false if you wish to sieve the update out.
type Filter interface {
	Filter(*Update) bool
}

// FilterFunc is basically a lightweight version of Filter.
type FilterFunc func(*Update) bool

func NewChain(parent Poller) *Chain {
	c := &Chain{}
	c.Poller = parent
	c.Filter = func(upd *Update) bool {
		for _, filter := range c.Filters {
			switch f := filter.(type) {
			case Filter:
				if !f.Filter(upd) {
					return false
				}

			case FilterFunc:
				if !f(upd) {
					return false
				}

			case func(*Update) bool:
				if !f(upd) {
					return false
				}
			}

		}

		return true
	}

	return c
}

// Chain is a chain of middle
type Chain struct {
	MiddlewarePoller

	// (Filter | FilterFunc | func(*Update) bool)
	Filters []interface{}
}

// Add accepts either Filter interface or FilterFunc
func (c *Chain) Add(filter interface{}) {
	switch filter.(type) {
	case Filter:
		break
	case FilterFunc:
		break
	case func(*Update) bool:
		break
	default:
		panic("telebot: unsupported filter type")
	}

	c.Filters = append(c.Filters, filter)
}
