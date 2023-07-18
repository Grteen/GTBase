package page

// Page is the basic unit store in disk and in xxx.pf file
// It is always 16KB
type Page struct {
	pageHeader PageHeader
	src        []byte
}

// PageHeader is the header info of a Page
type PageHeader struct {
	pageIndex int32
}
