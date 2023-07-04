package utils

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// Paginator 分页器
type Paginator struct {
	Request     *http.Request //请求
	PerPageNums int           // 一页显示多少条记录
	MaxPages    int           //一共有多少页
	nums        int64         // 一共有多少条记录
	pageRange   []int
	pageNums    int // pageNums总页数
	page        int
}

// PageNums 计算总页数
func (p *Paginator) PageNums() int {
	if p.pageNums != 0 {
		return p.pageNums
	}
	pageNums := math.Ceil(float64(p.nums) / float64(p.PerPageNums))
	if p.MaxPages > 0 {
		pageNums = math.Min(pageNums, float64(p.MaxPages))
	}
	p.pageNums = int(pageNums)
	return p.pageNums
}

// Nums 得到总记录数
func (p *Paginator) Nums() int64 {
	return p.nums
}

// SetNums 赋值总记录数
func (p *Paginator) SetNums(nums interface{}) {
	p.nums = AnyToInt64(nums)
}

// Page 得到当前传入的页数
func (p *Paginator) Page() int {
	if p.page != 0 {
		return p.page
	}
	if p.Request.Form == nil {
		_ = p.Request.ParseForm()
	}
	p.page, _ = strconv.Atoi(p.Request.Form.Get("index"))
	if p.page > p.PageNums() {
		// p.page = p.PageNums()
	}
	if p.page <= 0 {
		p.page = 1
	}
	return p.page
}

// Pages 分页
func (p *Paginator) Pages() []int {
	if p.pageRange == nil && p.nums > 0 {
		var pages []int
		pageNums := p.PageNums()
		page := p.Page()
		switch {
		case page >= pageNums-4 && pageNums > 9:
			start := pageNums - 9 + 1
			pages = make([]int, 9)
			for i := range pages {
				pages[i] = start + i
			}
		case page >= 5 && pageNums > 9:
			start := page - 5 + 1
			pages = make([]int, int(math.Min(9, float64(page+4+1))))
			for i := range pages {
				pages[i] = start + i
			}
		default:
			pages = make([]int, int(math.Min(9, float64(pageNums))))
			for i := range pages {
				pages[i] = i + 1
			}
		}
		p.pageRange = pages
	}
	return p.pageRange
}

// PageLink 得到页面的链接
func (p *Paginator) PageLink(page int) string {
	link, _ := url.ParseRequestURI(p.Request.RequestURI)
	values := link.Query()
	if page == 1 {
		values.Del("index")
	} else {
		values.Set("index", strconv.Itoa(page))
	}
	link.RawQuery = values.Encode()
	return link.String()
}

// PageLinkPrev 上一页链接
func (p *Paginator) PageLinkPrev() string {
	link := ""
	if p.HasPrev() {
		link = p.PageLink(p.Page() - 1)
	}
	return link
}

// PageLinkNext 下一页的链接
func (p *Paginator) PageLinkNext() string {
	link := ""
	if p.HasNext() {
		link = p.PageLink(p.Page() + 1)
	}
	return link
}

// PageLinkFirst 第一页的链接
func (p *Paginator) PageLinkFirst() string {
	return p.PageLink(1)
}

// PageLinkLast 最后一页的链接
func (p *Paginator) PageLinkLast() string {
	return p.PageLink(p.PageNums())
}

// HasPrev 是否存在上一页
func (p *Paginator) HasPrev() bool {
	return p.Page() > 1
}

// HasNext 是否存在下一页
func (p *Paginator) HasNext() bool {
	return p.Page() < p.PageNums()
}

// IsActive 是否点击
func (p *Paginator) IsActive(page int) bool {
	return p.Page() == page
}

// Offset 偏移量
func (p *Paginator) Offset() int {
	return (p.Page() - 1) * p.PerPageNums
}

// HasPages 是否存在页
func (p *Paginator) HasPages() bool {
	return p.PageNums() > 1
}

// NewPaginator 创建分页
func NewPaginator(req *http.Request, per int, nums interface{}) *Paginator {
	p := Paginator{}
	p.Request = req
	if per <= 0 {
		per = 10
	}
	p.PerPageNums = per
	p.SetNums(nums)
	return &p
}
