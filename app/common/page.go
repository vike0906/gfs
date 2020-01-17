package common

type Page struct {
	Total int32
	Data  interface{}
}

func NewPage(total int32, data interface{}) *Page {
	return &Page{Total: total, Data: data}
}
