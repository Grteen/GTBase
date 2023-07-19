package page

// PagePool caches all pages
// every read should read PagePool first
// if no cache in PagePool, it will read from disk and cache it
type PagePool struct {
	caches map[int32]*Page
}
