package site

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func ToJquerySelector(text string) string {
	if strings.HasPrefix(text, "<") && strings.HasSuffix(text, ">") {
		if strings.Contains(text, " ") {
			return fmt.Sprintf("selector:%s", getNodeInfo(text))
		} else {
			text = strings.ReplaceAll(text, "<", "")
			text = strings.ReplaceAll(text, ">", "")
			return fmt.Sprintf("selector:%s", text)
		}
	} else if strings.HasPrefix(text, "[['") {
		text = strings.ReplaceAll(text, "[['", "")
		text = strings.ReplaceAll(text, "']]", "")
		return fmt.Sprintf("text:%s", text)
	} else if strings.HasPrefix(text, "[[/src=") {
		return srcToSelector(text)
	} else if strings.HasPrefix(text, "[[{$") {
		return pipeline(text)
	}
	return text
}

func srcToSelector(text string) string {
	match := text
	match = strings.ReplaceAll(match, "[[/", "")
	match = strings.ReplaceAll(match, "'/]]", "")
	match = strings.ReplaceAll(match, "\\S+", ".*")
	match = strings.ReplaceAll(match, "(", "\\((")
	match = strings.ReplaceAll(match, ")", ")\\)")
	reg, _ := regexp.Compile(match)
	if reg.MatchString(text) {
		group := reg.FindAllStringSubmatch(text, -1)
		return fmt.Sprintf("selector:[src*=%s]", group[0][1])
	} else {
		return text
	}
}

func getNodeInfo(text string) string {
	var buff bytes.Buffer
	buff.WriteString(text)
	node, err := goquery.NewDocumentFromReader(&buff)
	if err != nil {
		return ""
	}
	sel := node.Find("body").Children()
	selector := goquery.NodeName(sel)
	if class, ok := sel.Attr("class"); ok {
		classList := strings.Split(class, " ")
		for _, s := range classList {
			selector += fmt.Sprintf(".%s", s)
		}
	} else {
		id, _ := sel.Attr("id")
		selector += fmt.Sprintf("#%s", id)
	}

	return selector
}

func pipeline(text string) string {
	// [[{$('.title h2').text()}]]
	reg1, _ := regexp.Compile(`^\[\[\{\$\('([#\w\d\s\.\-\:\[\]\"\=]*)'\)\.(text|html)\(\)\}\]\]$`)
	if reg1.MatchString(text) {
		group := reg1.FindAllStringSubmatch(text, -1)
		return fmt.Sprintf("pipeline:selector:%s|%s", group[0][1], group[0][2])
	}
	// [[{$('.main_editor ').find('.title').text()}]]
	reg2, _ := regexp.Compile(`^\[\[\{\$\('([#\w\d\s\.\-\:\[\]\"\=]*)'\)\.find\('([#\w\d\s\.\-\:\[\]\"\=]*)'\)\.(text|html)\(\)\}\]\]$`)
	if reg2.MatchString(text) {
		group := reg2.FindAllStringSubmatch(text, -1)
		return fmt.Sprintf("pipeline:selector:%s|%s|%s", group[0][1], group[0][2], group[0][3])
	}
	// [[{$('.content___2CL42').find('.summary___3oqrM').parent().html()}]]
	reg3, _ := regexp.Compile(`^\[\[\{\$\('([#\w\d\s\.\-\:\[\]\"\=]*)'\)\.find\('([#\w\d\s\.\-\:\[\]\"\=]*)'\)\.parent\(\)\.(text|html)\(\)\}\]\]$`)
	if reg3.MatchString(text) {
		group := reg3.FindAllStringSubmatch(text, -1)
		return fmt.Sprintf("pipeline:selector:%s|selector:%s|parent|%s", group[0][1], group[0][2], group[0][3])
	}
	// [[{$($('section.content')[1]).html()}]]
	reg4, _ := regexp.Compile(`^\[\[\{\$\(\$\('([#\w\d\s\.\-\:\[\]\"\=]*)'\)\[(\d)\]\)\.(text|html)\(\)\}\]\]$`)
	if reg4.MatchString(text) {
		group := reg4.FindAllStringSubmatch(text, -1)
		return fmt.Sprintf("pipeline:selector:%s|index:%s|%s", group[0][1], group[0][2], group[0][3])
	}
	// [[{$('meta[name=description]').attr('content')}]]
	reg5, _ := regexp.Compile(`^\[\[\{\$\('([#\w\d\s\.\-\:\[\]\"\=]*)'\)\.attr\('([\w]+)'\)\}\]\]$`)
	if reg5.MatchString(text) {
		group := reg5.FindAllStringSubmatch(text, -1)
		return fmt.Sprintf("pipeline:selector:%s|attr:%s", group[0][1], group[0][2])
	}
	// [[{$('title').text().replace( ' - 少数派', '' )}]]
	//reg6, _ := regexp.Compile(`\[\[\{\$\('([\w\d\s\.-]*)'\)\.(text|html)\(\).replace\(\s+'(.*)',.*\)\}\]\]`)
	//if reg6.MatchString(text) {
	//	group := reg6.FindAllStringSubmatch(text, -1)
	//	return fmt.Sprintf("pipeline:selector:%s|%s|replace:%s", group[0][1], group[0][2], group[0][3])
	//}
	// [[{$('.RichContent-inner')}]]
	reg7, _ := regexp.Compile(`^\[\[\{\$\('([#\w\d\s\.\-\:\[\]\"\=]*)'\)\}\]\]$`)
	if reg7.MatchString(text) {
		group := reg7.FindAllStringSubmatch(text, -1)
		return fmt.Sprintf("pipeline:selector:%s|html", group[0][1])
	}
	// [[{$('meta[name=Description]').attr('content')||$('meta[name=description]').attr('content')}]]
	reg8, _ := regexp.Compile(`\$\('([#\w\d\s\.\-\:\[\]\"\=]*)'\)\.attr\('([\w]+)'\)+`)
	if reg8.MatchString(text) {
		groups := reg8.FindAllStringSubmatch(text, -1)
		match := "pipeline:"
		for i, group := range groups {
			if i > 0 {
				match += "||"
			}
			match += fmt.Sprintf("selector:%s|attr:%s", group[1], group[2])
		}
		return match
	}
	// [[{$('.view_content').html()||$('#view_content').html()||$('#articleContent').html()||$('.forum-viewthread-article-box').html()}]]
	reg9, _ := regexp.Compile(`\$\('([#\w\d\s\.\-\:\[\]]*)'\)\.(text|html)\(\)`)
	if reg9.MatchString(text) {
		groups := reg9.FindAllStringSubmatch(text, -1)
		match := "pipeline:"
		for i, group := range groups {
			if i > 0 {
				match += "||"
			}
			match += fmt.Sprintf("selector:%s|%s", group[1], group[2])
		}
		return match
	}
	// [[{$($('.article-con')[0]).html()}]]
	reg10, _ := regexp.Compile(`^\[\[\{\$\(\$\('([#\w\d\s\.\-\:\[\]]*)'\)\[(\d)\]\)\.(text|html)\(\)\}\]\]$`)
	if reg10.MatchString(text) {
		group := reg10.FindAllStringSubmatch(text, -1)
		return fmt.Sprintf("pipeline:selector:%s|index:%s|%s", group[0][1], group[0][2], group[0][3])
	}
	// [[{$('h3[itemprop="description"]').text()}]]
	return text
}
