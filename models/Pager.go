package models

import (
	"bytes"
	"fmt"
)

type Pager struct {
	Page int //第几页
	Pagesize int //每页大小
	Totalnum int//总页数
	urlpath string //每页所对应的url
}

//创建Pager对象
func NewPager(page, pagesize, totalnum int, urlpath string) *Pager {
	pager := new(Pager)
	pager.Page = page
	pager.Pagesize = pagesize
	pager.Totalnum = totalnum
	pager.urlpath = urlpath
	return pager
}

//修改page
func (this *Pager) SetPage(page int) {
	this.Page = page
}

//修改pagesize
func (this *Pager) SetPagesize(pagesize int) {
	this.Pagesize = pagesize
}

//设置总数量
func (this *Pager) SetTotalnum(totalnum int) {
	this.Totalnum = totalnum
}

//设置rootpath
func (this *Pager) SetUrlpath(urlpath string) {
	this.urlpath = urlpath
}

func (this *Pager) url(page int) string {
	return fmt.Sprintf(this.urlpath, page)
}

func (this *Pager) ToString() string{//str2html

	if this.Totalnum <= this.Pagesize {
		return ""
	}
	offset := 5
	linknum := 10
	var totalpage int
	var from int//从哪一页开始显示
	var to int //显示到哪一页
	if this.Totalnum % this.Pagesize != 0 {
		totalpage = this.Totalnum / this.Pagesize + 1
	}else {
		totalpage = this.Totalnum / this.Pagesize
	}
	//fmt.Println("总页面 = ", totalpage)
	if totalpage < linknum {
		from = 1
		to = totalpage
	}else {
		from = this.Page - offset
		to = from + linknum
		if from < 1 {
			from = 1
			to = from + linknum - 1
		}else if to > totalpage {
			to = totalpage
			from = to - linknum + 1//20 - 10 + 1 = 11(11-20) 21 - 10 + 1 = 12(12-21)
		}
	}
	//开辟空间
	var buf bytes.Buffer
	buf.WriteString("<div class='page'>")
	//上一页
	if this.Page > 1 {
		buf.WriteString(fmt.Sprintf("<a href='%s'>&laquo;</a>", this.url(this.Page-1)))//<<
	}

	for i := from; i <= to; i++ {
		if i == this.Page {
			buf.WriteString(fmt.Sprintf("<b>%d</b>", i))
		}else {
			buf.WriteString(fmt.Sprintf("<a href='%s'>%d</a>", this.url(i), i))
		}
	}

	//设置下一页标签
	if this.Page < totalpage {
		buf.WriteString(fmt.Sprintf("<a href='%s'>&raquo;</a>", this.url(this.Page+1)))//>>
	}
    buf.WriteString("</div>")
    str := buf.String()
    return str
}



