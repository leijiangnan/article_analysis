package service

import (
	"article-analysis/internal/model"
	"article-analysis/internal/repository"
	"article-analysis/pkg/logger"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type CrawlerService struct {
	articleRepo *repository.ArticleRepository
	log         *logger.Logger
	client      *http.Client
	maxArticles int
}

func NewCrawlerService(repo *repository.ArticleRepository, log *logger.Logger) *CrawlerService {
	return &CrawlerService{
		articleRepo: repo,
		log:         log,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxArticles: 2,
	}
}

type CrawledArticle struct {
	Title   string
	Author  string
	Date    string
	Content string
	URL     string
}

type CrawlResult struct {
	TotalFound   int
	CrawledCount int
	Articles     []CrawledArticle
	Errors       []string
}

func (s *CrawlerService) CrawlArticles(startURL string, count int) (*CrawlResult, error) {
	if count <= 0 {
		return nil, errors.New("爬取数量必须大于0")
	}

	if count > s.maxArticles {
		s.log.Info("限制爬取数量", zap.Int("requested", count), zap.Int("max", s.maxArticles))
		count = s.maxArticles
	}

	parsedURL, err := url.Parse(startURL)
	if err != nil {
		return nil, errors.New("无效的URL")
	}
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	s.log.Info("开始爬取文章",
		zap.String("url", startURL),
		zap.Int("count", count))

	doc, err := s.fetchPage(startURL)
	if err != nil {
		return nil, err
	}

	links := s.extractArticleLinksWithDates(doc, baseURL)
	s.log.Info("找到文章链接", zap.Int("count", len(links)))

	result := &CrawlResult{
		TotalFound: len(links),
		Articles:   make([]CrawledArticle, 0),
		Errors:     make([]string, 0),
	}

	crawledCount := 0
	for i, linkInfo := range links {
		if crawledCount >= count {
			break
		}

		s.log.Info("正在爬取文章", zap.Int("index", i+1), zap.String("url", linkInfo.URL))

		article, err := s.crawlArticle(linkInfo.URL, linkInfo.Title)
		if err != nil {
			s.log.Warn("爬取文章失败", zap.String("url", linkInfo.URL), zap.Error(err))
			result.Errors = append(result.Errors, linkInfo.URL+": "+err.Error())
			continue
		}

		if article.Date == "" && linkInfo.Date != "" {
			article.Date = linkInfo.Date
		}

		result.Articles = append(result.Articles, *article)
		crawledCount++
		s.log.Info("文章爬取成功", zap.String("title", article.Title))
	}

	result.CrawledCount = crawledCount
	return result, nil
}

func (s *CrawlerService) SaveCrawledArticles(articles []CrawledArticle) ([]model.Article, []error) {
	savedArticles := make([]model.Article, 0)
	errs := make([]error, 0)

	for _, article := range articles {
		title := strings.TrimSpace(article.Title)
		if title == "" {
			title = "无标题"
		}

		exists, err := s.articleRepo.ExistsByTitle(title)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if exists {
			s.log.Warn("文章已存在，跳过", zap.String("title", title))
			continue
		}

		author := strings.TrimSpace(article.Author)
		if author == "" {
			author = "未知作者"
		}

		content := strings.TrimSpace(article.Content)
		if content == "" {
			s.log.Warn("文章内容为空，跳过", zap.String("title", title))
			continue
		}

		dbArticle := &model.Article{
			Title:    title,
			Author:   author,
			Content:  content,
			FilePath: "crawled:" + article.URL,
			FileSize: int64(len(content)),
		}

		if err := s.articleRepo.Create(dbArticle); err != nil {
			errs = append(errs, err)
			continue
		}

		savedArticles = append(savedArticles, *dbArticle)
		s.log.Info("文章保存成功", zap.String("title", title), zap.Uint64("id", dbArticle.ID))
	}

	return savedArticles, errs
}

func (s *CrawlerService) fetchPage(urlStr string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.7,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Sec-Ch-Ua", `"Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	if strings.Contains(urlStr, "mp.weixin.qq.com") {
		req.Header.Set("Referer", "https://mp.weixin.qq.com/")
		req.Header.Set("Cookie", "wxtokenkey=777; pac_uid=0_Z9cww77Z6p1HQ; _qimei_uuid42=19c050e1c08100cec8f11ef416dd8038b94bdb108b; _qimei_fingerprint=4955c71558659de38ad284d3a93a002b; _qimei_h38=f2d212b1c8f11ef416dd803803000004119c05; poc_sid=HC0zoGmjUe4GVSn9V68EnjaB4xYxyqEWWGtSKp8S; __root_domain_v=.weixin.qq.com; _qddaz=QD.713570345653712")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP请求失败: " + resp.Status)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func (s *CrawlerService) extractArticleLinks(doc *goquery.Document, baseURL string) []string {
	links := make([]string, 0)
	seen := make(map[string]bool)

	doc.Find("a[href]").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists || href == "" {
			return
		}

		href = strings.TrimSpace(href)
		if href == "#" || href == "javascript:void(0)" {
			return
		}

		if strings.HasPrefix(href, "/") {
			href = baseURL + href
		} else if !strings.HasPrefix(href, "http://") && !strings.HasPrefix(href, "https://") {
			return
		}

		parsedHref, err := url.Parse(href)
		if err != nil {
			return
		}

		parsedBase, _ := url.Parse(baseURL)
		if parsedHref.Host != parsedBase.Host {
			return
		}

		ext := strings.ToLower(href)
		if strings.HasSuffix(ext, ".jpg") || strings.HasSuffix(ext, ".png") ||
			strings.HasSuffix(ext, ".gif") || strings.HasSuffix(ext, ".pdf") ||
			strings.HasSuffix(ext, ".zip") || strings.HasSuffix(ext, ".doc") ||
			strings.HasSuffix(ext, ".docx") {
			return
		}

		if !seen[href] {
			seen[href] = true
			links = append(links, href)
		}
	})

	return links
}

type ArticleLinkInfo struct {
	URL   string
	Date  string
	Title string
}

func (s *CrawlerService) extractArticleLinksWithDates(doc *goquery.Document, baseURL string) []ArticleLinkInfo {
	links := make([]ArticleLinkInfo, 0)
	seen := make(map[string]bool)

	// 获取整个页面的文本，用于提取日期
	// 微信文章列表页的格式：序号 标题作者 日期
	fullText := doc.Find("#js_content").Text()
	if fullText == "" {
		fullText = doc.Text()
	}

	doc.Find("a[href]").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists || href == "" {
			return
		}

		href = strings.TrimSpace(href)
		if href == "#" || href == "javascript:void(0)" {
			return
		}

		if strings.HasPrefix(href, "/") {
			href = baseURL + href
		} else if !strings.HasPrefix(href, "http://") && !strings.HasPrefix(href, "https://") {
			return
		}

		parsedHref, err := url.Parse(href)
		if err != nil {
			return
		}

		parsedBase, _ := url.Parse(baseURL)
		if parsedHref.Host != parsedBase.Host {
			return
		}

		ext := strings.ToLower(href)
		if strings.HasSuffix(ext, ".jpg") || strings.HasSuffix(ext, ".png") ||
			strings.HasSuffix(ext, ".gif") || strings.HasSuffix(ext, ".pdf") ||
			strings.HasSuffix(ext, ".zip") || strings.HasSuffix(ext, ".doc") ||
			strings.HasSuffix(ext, ".docx") {
			return
		}

		if !seen[href] {
			seen[href] = true
			title := strings.TrimSpace(sel.Text())
			date := s.extractDateForTitle(fullText, title)
			links = append(links, ArticleLinkInfo{
				URL:   href,
				Date:  date,
				Title: title,
			})
		}
	})

	return links
}

// extractDateForTitle 从页面文本中提取指定标题对应的日期
func (s *CrawlerService) extractDateForTitle(fullText, title string) string {
	if title == "" || fullText == "" {
		return ""
	}

	// 在文本中查找标题位置
	idx := strings.Index(fullText, title)
	if idx == -1 {
		return ""
	}

	// 从标题位置开始，向后查找100个字符内的日期
	endIdx := idx + len(title) + 100
	if endIdx > len(fullText) {
		endIdx = len(fullText)
	}
	subText := fullText[idx:endIdx]

	// 提取日期
	dateMatch := regexp.MustCompile(`(\d{4}-\d{1,2}-\d{1,2})`).FindString(subText)
	return dateMatch
}

func (s *CrawlerService) isCaptchaPage(doc *goquery.Document) bool {
	html, _ := doc.Html()
	text := doc.Text()
	return strings.Contains(text, "验证") || strings.Contains(text, "验证码") ||
		strings.Contains(text, "环境异常") || strings.Contains(text, "尝试太多了") ||
		len(strings.TrimSpace(text)) < 100 || !strings.Contains(html, "js_content")
}

func (s *CrawlerService) crawlArticle(urlStr string, preExtractedTitle string) (*CrawledArticle, error) {
	doc, err := s.fetchPage(urlStr)
	if err != nil {
		return nil, err
	}

	article := &CrawledArticle{
		URL: urlStr,
	}

	if s.isCaptchaPage(doc) {
		s.log.Warn("检测到验证码页面或异常页面", zap.String("url", urlStr))

		// 1. 优先使用预提取的标题（从起始页面的链接获取）
		if preExtractedTitle != "" {
			article.Title = preExtractedTitle
		} else {
			// 如果没有预提取标题，尝试从页面获取
			article.Title = s.extractTitle(doc)
			if article.Title == "" {
				article.Title = doc.Find("title").Text()
			}
		}

		// 2. 尝试从页面提取日期
		article.Date = s.extractDate(doc)

		// 3. 仅打印异常页面显示的文字
		// 为了提取纯文本，先移除不可见元素和脚本
		// 注意：extractTitle 可能依赖 meta 标签，所以要在提取标题后执行移除
		doc.Find("script, style, noscript, link, meta, iframe").Remove()
		// 移除内联样式为 display:none 的元素
		doc.Find("[style*='display:none']").Remove()
		doc.Find("[style*='display: none']").Remove()

		text := doc.Find("body").Text()

		// 清洗多余空白字符
		text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
		article.Content = strings.TrimSpace(text)

		// 即使是验证码页面也返回成功，以便保存查看
		return article, nil
	}

	article.Title = s.extractTitle(doc)
	article.Author = s.extractAuthor(doc)
	article.Date = s.extractDate(doc)
	article.Content = s.extractContent(doc)
	return article, nil
}

func (s *CrawlerService) extractTitle(doc *goquery.Document) string {
	selectors := []string{
		"#activity-name",
		".rich_media_title",
		"h1",
		"meta[property='og:title']",
		"meta[name='title']",
		"title",
	}

	for _, sel := range selectors {
		if sel == "h1" || sel == "#activity-name" || sel == ".rich_media_title" {
			text := strings.TrimSpace(doc.Find(sel).First().Text())
			if text != "" && len(text) < 200 {
				return text
			}
		} else if strings.HasPrefix(sel, "meta") {
			if content, exists := doc.Find(sel).Attr("content"); exists && content != "" {
				return strings.TrimSpace(content)
			}
		} else {
			text := strings.TrimSpace(doc.Find(sel).Text())
			if text != "" {
				if len(text) > 200 {
					text = text[:200]
				}
				return text
			}
		}
	}

	return ""
}

func (s *CrawlerService) extractAuthor(doc *goquery.Document) string {
	selectors := []string{
		"#js_name",
		".rich_media_meta_nickname",
		"meta[name='author']",
		"meta[property='article:author']",
		"[class*='author']",
		"[class*='Author']",
		"[class*='writer']",
		"[class*='source']",
	}

	for _, sel := range selectors {
		if strings.HasPrefix(sel, "meta") {
			if content, exists := doc.Find(sel).Attr("content"); exists && content != "" {
				return strings.TrimSpace(content)
			}
		} else {
			text := strings.TrimSpace(doc.Find(sel).First().Text())
			if text != "" && len(text) < 100 {
				text = regexp.MustCompile(`作者[：:]\s*`).ReplaceAllString(text, "")
				text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
				return strings.TrimSpace(text)
			}
		}
	}

	return ""
}

func (s *CrawlerService) extractDate(doc *goquery.Document) string {
	dateSelectors := []string{
		"#publish_time",
		"meta[property='article:published_time']",
		"meta[name='publishdate']",
		"meta[name='date']",
		"time",
		"[class*='date']",
		"[class*='time']",
	}

	for _, sel := range dateSelectors {
		if sel == "#publish_time" {
			text := strings.TrimSpace(doc.Find(sel).Text())
			if text != "" && len(text) < 50 {
				text = strings.ReplaceAll(text, "年", "-")
				text = strings.ReplaceAll(text, "月", "-")
				text = strings.ReplaceAll(text, "日", "")
				text = regexp.MustCompile(`\s*-\s*`).ReplaceAllString(text, "-")
				text = regexp.MustCompile(`^[^0-9]*([0-9-]+).*`).ReplaceAllString(text, "$1")
				if regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}$`).MatchString(text) {
					return text
				}
				return strings.TrimSpace(text)
			}
		} else if strings.HasPrefix(sel, "meta") {
			if content, exists := doc.Find(sel).Attr("content"); exists && content != "" {
				dateStr := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})`).FindString(content)
				if dateStr != "" {
					return dateStr
				}
				return strings.TrimSpace(content)
			}
		} else {
			text := strings.TrimSpace(doc.Find(sel).First().Text())
			if text != "" && len(text) < 50 && !strings.Contains(text, "作者") {
				return text
			}
		}
	}

	// 尝试从文章内容开头匹配日期（微信文章日期通常是动态加载的）
	content := doc.Find("#js_content").Text()
	if content == "" {
		content = doc.Find("body").Text()
	}
	// 从内容前500字符中匹配日期
	if len(content) > 500 {
		content = content[:500]
	}
	datePatterns := []string{
		`(\d{4})年(\d{1,2})月(\d{1,2})日`,
		`(\d{4})-(\d{1,2})-(\d{1,2})`,
		`(\d{4})/(\d{1,2})/(\d{1,2})`,
	}
	for _, pattern := range datePatterns {
		matches := regexp.MustCompile(pattern).FindStringSubmatch(content)
		if len(matches) == 4 {
			year := matches[1]
			month := matches[2]
			day := matches[3]
			if len(month) == 1 {
				month = "0" + month
			}
			if len(day) == 1 {
				day = "0" + day
			}
			return year + "-" + month + "-" + day
		}
	}

	return ""
}

func (s *CrawlerService) extractContent(doc *goquery.Document) string {
	contentSelectors := []string{
		"#js_content",
		".rich_media_content",
		".js_name",
		"article",
		"[class*='content']",
		"[class*='article-body']",
		"[class*='post-content']",
		"[class*='entry-content']",
		"[class*='main-content']",
		".article",
		"#article",
	}

	var contentNode *goquery.Selection

	for _, sel := range contentSelectors {
		node := doc.Find(sel).First()
		if node.Length() > 0 {
			contentNode = node
			break
		}
	}

	if contentNode == nil {
		contentNode = doc.Find("body")
	}

	contentNode.Find("script, style, nav, header, footer, aside, .sidebar, .comment, .advertisement, .ad").Remove()

	var paragraphs []string
	contentNode.Find("p, br, div, section").Each(func(i int, sel *goquery.Selection) {
		sel.Find("a").Each(func(j int, aSel *goquery.Selection) {
			text := aSel.Text()
			aSel.ReplaceWithHtml(text)
		})

		text := sel.Text()
		text = strings.TrimSpace(text)
		if text != "" {
			paragraphs = append(paragraphs, text)
		}
	})

	if len(paragraphs) == 0 {
		text := contentNode.Text()
		text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
		return strings.TrimSpace(text)
	}

	var result []string
	for _, p := range paragraphs {
		p = regexp.MustCompile(`\s+`).ReplaceAllString(p, " ")
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return strings.Join(result, "\n")
}

func (s *CrawlerService) SetMaxArticles(max int) {
	if max > 0 {
		s.maxArticles = max
	}
}
