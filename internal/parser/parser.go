package parser

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type RankingTableParser struct {
	root *html.Node
}

func NewRankingTableParser(r io.Reader) (p RankingTableParser, err error) {
	root, err := html.Parse(r)
	if err != nil {
		return
	}
	p.root = root

	return p, err
}

func (p RankingTableParser) Parse() (table []map[string]string, err error) {
	t := p.findRankingTable()
	if t == nil {
		return
	}
	header, rows := p.findRanking(t)
	if len(rows) == 0 {
		return
	}

	table = make([]map[string]string, 0, len(rows))
	for _, r := range rows {
		s, e := p.transformRowIntoStock(header, r)
		if e != nil {
			log.Printf("failed to parse row due to err: %v", e)
			continue
		}
		table = append(table, s)
	}

	return
}

func (p RankingTableParser) findRankingTable() *html.Node {
	return p.findByID(p.root, "table-ranking")
}

func (p RankingTableParser) findRanking(table *html.Node) (*html.Node, []*html.Node) {
	header := p.findFirstByTag(p.findFirstByTag(table, "thead"), "tr")
	rows := p.findAllByTag(p.findFirstByTag(table, "tbody"), "tr")
	return header, rows
}

func (p RankingTableParser) getNodeText(n *html.Node) (s string) {
	if n == nil {
		return
	}

	b := strings.Builder{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.TextNode {
			continue
		}

		if b.Len() > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(c.Data)
	}

	if b.Len() == 0 {
		return p.getNodeText(n.FirstChild)
	}

	return b.String()
}

func (p RankingTableParser) transformRowIntoStock(h, r *html.Node) (m map[string]string, err error) {
	titles := p.findAllByTag(h, "th")
	cells := p.findAllByTag(r, "td")
	m = make(map[string]string, len(titles))
	for i := 0; i < len(titles); i++ {
		title := p.getNodeText(titles[i])
		data := p.getNodeText(cells[i])
		m[title] = data
	}
	return
}

func (p RankingTableParser) getAttribute(n *html.Node, key string) (val string, ok bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}

func (p RankingTableParser) isSameID(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := p.getAttribute(n, "id")
		if ok && s == id {
			return true
		}
	}

	return false
}

func (p RankingTableParser) findByID(n *html.Node, id string) (target *html.Node) {
	if p.isSameID(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		target = p.findByID(c, id)
		if target != nil {
			return target
		}
	}

	return target
}

func (p RankingTableParser) findFirstByTag(n *html.Node, tag string) (target *html.Node) {
	if n == nil {
		return nil
	}
	if p.isSameTag(n, tag) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		target = p.findFirstByTag(c, tag)
		if target != nil {
			return target
		}
	}

	return target
}

func (p RankingTableParser) findAllByTag(n *html.Node, tag string) (targets []*html.Node) {
	first := p.findFirstByTag(n, tag)
	if first == nil {
		return
	}

	targets = append(targets, first)
	for s := first.NextSibling; s != nil; s = s.NextSibling {
		if p.isSameTag(s, tag) {
			targets = append(targets, s)
		}
	}

	return targets
}

func (p RankingTableParser) isSameTag(n *html.Node, tag string) bool {
	return n != nil && n.Data == tag
}
