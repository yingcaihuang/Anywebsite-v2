package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Article Article     `json:"article,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// 文章结构
type Article struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Slug      string     `json:"slug"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// 创建文章请求结构
type CreateArticleRequest struct {
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Slug      string     `json:"slug,omitempty"`
	Status    string     `json:"status"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// 更新文章请求结构
type UpdateArticleRequest struct {
	Title     string     `json:"title,omitempty"`
	Content   string     `json:"content,omitempty"`
	Status    string     `json:"status,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// 测试客户端
type TestClient struct {
	baseURL    string
	httpClient *http.Client
	sessionID  string
}

// 创建新的测试客户端
func NewTestClient(baseURL string) *TestClient {
	jar, _ := cookiejar.New(nil)
	return &TestClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
	}
}

// 管理员登录
func (c *TestClient) Login(username, password string) error {
	fmt.Println("🔐 正在进行管理员登录...")

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", c.baseURL+"/admin/login", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("创建登录请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("执行登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("登录失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	fmt.Println("✅ 管理员登录成功")
	return nil
}

// 创建文章
func (c *TestClient) CreateArticle(req CreateArticleRequest) (*Article, error) {
	fmt.Printf("📝 正在创建文章: %s\n", req.Title)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/api/articles", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API返回错误: %s", apiResp.Error)
	}

	fmt.Printf("✅ 文章创建成功，ID: %s, Slug: %s\n", apiResp.Article.ID, apiResp.Article.Slug)
	return &apiResp.Article, nil
}

// 获取文章
func (c *TestClient) GetArticle(id string) (*Article, error) {
	fmt.Printf("📖 正在获取文章: %s\n", id)

	resp, err := c.httpClient.Get(c.baseURL + "/api/articles/" + id)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API返回错误: %s", apiResp.Error)
	}

	fmt.Printf("✅ 文章获取成功: %s\n", apiResp.Article.Title)
	return &apiResp.Article, nil
}

// 更新文章
func (c *TestClient) UpdateArticle(id string, req UpdateArticleRequest) (*Article, error) {
	fmt.Printf("✏️ 正在更新文章: %s\n", id)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequest("PUT", c.baseURL+"/api/articles/"+id, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API返回错误: %s", apiResp.Error)
	}

	fmt.Printf("✅ 文章更新成功: %s\n", apiResp.Article.Title)
	return &apiResp.Article, nil
}

// 删除文章
func (c *TestClient) DeleteArticle(id string) error {
	fmt.Printf("🗑️ 正在删除文章: %s\n", id)

	httpReq, err := http.NewRequest("DELETE", c.baseURL+"/api/articles/"+id, nil)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("API返回错误: %s", apiResp.Error)
	}

	fmt.Printf("✅ 文章删除成功\n")
	return nil
}

// 获取Web页面内容
func (c *TestClient) GetWebPage(slug string) (string, error) {
	fmt.Printf("🌐 正在获取Web页面: /articles/%s\n", slug)

	resp, err := c.httpClient.Get(c.baseURL + "/articles/" + slug)
	if err != nil {
		return "", fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Web页面返回错误状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	fmt.Printf("✅ Web页面获取成功，内容长度: %d 字节\n", len(body))
	return string(body), nil
}

// 运行完整测试
func runFullTest() error {
	fmt.Println("🚀 开始运行完整API测试")
	fmt.Println("=" * 50)

	// 创建测试客户端
	client := NewTestClient("http://localhost:8080")

	// 1. 管理员登录
	if err := client.Login("admin", "password"); err != nil {
		return fmt.Errorf("登录测试失败: %v", err)
	}

	// 2. 创建文章测试
	createReq := CreateArticleRequest{
		Title: "API测试文章",
		Content: `<h1>这是API测试文章</h1>
<p>这是一个通过<strong>Go API测试</strong>创建的文章。</p>
<blockquote>测试引用内容</blockquote>
<ul>
<li>测试列表项1</li>
<li>测试列表项2</li>
</ul>
<pre><code>func main() {
    fmt.Println("Hello, World!")
}</code></pre>`,
		Status: "published",
		Slug:   "api-test-article",
	}

	article, err := client.CreateArticle(createReq)
	if err != nil {
		return fmt.Errorf("创建文章测试失败: %v", err)
	}

	// 3. 获取文章测试
	retrievedArticle, err := client.GetArticle(article.ID)
	if err != nil {
		return fmt.Errorf("获取文章测试失败: %v", err)
	}

	// 验证获取的文章内容
	if retrievedArticle.Title != article.Title {
		return fmt.Errorf("获取的文章标题不匹配: 期望 %s, 实际 %s", article.Title, retrievedArticle.Title)
	}

	// 4. 更新文章测试
	updateReq := UpdateArticleRequest{
		Title: "更新后的API测试文章",
		Content: `<h1>这是更新后的API测试文章</h1>
<p>文章内容已通过<strong>Go API测试</strong>更新。</p>
<div style="background: #f0f8ff; padding: 15px; border-radius: 5px;">
<h3>更新内容</h3>
<p>这是新增的内容，用来验证更新功能。</p>
</div>
<hr>
<p><em>更新时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</em></p>`,
		Status: "published",
	}

	updatedArticle, err := client.UpdateArticle(article.ID, updateReq)
	if err != nil {
		return fmt.Errorf("更新文章测试失败: %v", err)
	}

	// 验证更新的文章内容
	if updatedArticle.Title != updateReq.Title {
		return fmt.Errorf("更新的文章标题不匹配: 期望 %s, 实际 %s", updateReq.Title, updatedArticle.Title)
	}

	// 5. 验证Web页面渲染
	webContent, err := client.GetWebPage(article.Slug)
	if err != nil {
		return fmt.Errorf("获取Web页面测试失败: %v", err)
	}

	// 检查Web页面是否包含更新后的内容
	if !strings.Contains(webContent, "更新后的API测试文章") {
		return fmt.Errorf("Web页面内容未正确更新")
	}

	if !strings.Contains(webContent, "<h1>这是更新后的API测试文章</h1>") {
		return fmt.Errorf("Web页面HTML内容渲染异常")
	}

	fmt.Println("✅ Web页面HTML渲染验证通过")

	// 6. 删除文章测试
	if err := client.DeleteArticle(article.ID); err != nil {
		return fmt.Errorf("删除文章测试失败: %v", err)
	}

	// 7. 验证文章已删除
	_, err = client.GetArticle(article.ID)
	if err == nil {
		return fmt.Errorf("文章删除后仍能获取，删除功能异常")
	}
	fmt.Println("✅ 文章删除验证通过")

	fmt.Println("=" * 50)
	fmt.Println("🎉 所有API测试通过！")
	return nil
}

// 运行HTML渲染专项测试
func runHTMLRenderingTest() error {
	fmt.Println("🎨 开始运行HTML渲染专项测试")
	fmt.Println("=" * 50)

	client := NewTestClient("http://localhost:8080")

	// 登录
	if err := client.Login("admin", "password"); err != nil {
		return fmt.Errorf("登录失败: %v", err)
	}

	// 创建包含丰富HTML内容的文章
	htmlContent := `<h1>HTML渲染测试</h1>
<h2>各种HTML元素测试</h2>

<h3>文本格式</h3>
<p>这是<strong>粗体</strong>、<em>斜体</em>、<u>下划线</u>和<code>行内代码</code>的测试。</p>

<h3>引用和列表</h3>
<blockquote>
这是一个引用块，用来测试样式是否正确应用。
</blockquote>

<ul>
<li>无序列表项1</li>
<li>无序列表项2 <a href="#">带链接</a></li>
</ul>

<ol>
<li>有序列表项1</li>
<li>有序列表项2</li>
</ol>

<h3>代码展示</h3>
<pre><code>package main

import "fmt"

func main() {
    fmt.Println("Hello, HTML!")
    
    // 测试代码高亮
    var message = "渲染测试"
    fmt.Println(message)
}</code></pre>

<h3>样式测试</h3>
<p style="color: red; font-size: 18px;">红色大字体文本</p>
<p style="background: yellow; padding: 10px;">黄色背景文本</p>

<div style="border: 2px solid blue; padding: 15px; margin: 10px 0;">
<h4>自定义容器</h4>
<p>这是一个蓝色边框的容器。</p>
</div>

<h3>表格测试</h3>
<table>
<tr><th>项目</th><th>状态</th><th>说明</th></tr>
<tr><td>HTML标签</td><td>✅</td><td>正常渲染</td></tr>
<tr><td>CSS样式</td><td>✅</td><td>正常应用</td></tr>
<tr><td>JavaScript</td><td>❌</td><td>已过滤</td></tr>
</table>

<hr>
<p><small>测试时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</small></p>`

	createReq := CreateArticleRequest{
		Title:   "HTML渲染专项测试",
		Content: htmlContent,
		Status:  "published",
		Slug:    "html-rendering-test",
	}

	article, err := client.CreateArticle(createReq)
	if err != nil {
		return fmt.Errorf("创建测试文章失败: %v", err)
	}

	// 获取Web页面并验证HTML渲染
	webContent, err := client.GetWebPage(article.Slug)
	if err != nil {
		return fmt.Errorf("获取Web页面失败: %v", err)
	}

	// 检查各种HTML元素是否正确渲染
	testCases := []struct {
		name    string
		pattern string
	}{
		{"H1标题", "<h1>HTML渲染测试</h1>"},
		{"粗体文字", "<strong>粗体</strong>"},
		{"斜体文字", "<em>斜体</em>"},
		{"行内代码", "<code>行内代码</code>"},
		{"引用块", "<blockquote>"},
		{"无序列表", "<ul>"},
		{"有序列表", "<ol>"},
		{"代码块", "<pre><code>"},
		{"内联样式", `style="color: red`},
		{"表格", "<table>"},
		{"分割线", "<hr>"},
	}

	for _, tc := range testCases {
		if strings.Contains(webContent, tc.pattern) {
			fmt.Printf("✅ %s 渲染正确\n", tc.name)
		} else {
			fmt.Printf("❌ %s 渲染异常\n", tc.name)
		}
	}

	// 检查是否是纯净的HTML内容（没有多余的页面结构）
	if !strings.Contains(webContent, "article-header") && !strings.Contains(webContent, "article-footer") {
		fmt.Println("✅ 页面使用纯净模板，没有多余结构")
	} else {
		fmt.Println("❌ 页面仍包含多余的结构元素")
	}

	// 清理测试数据
	if err := client.DeleteArticle(article.ID); err != nil {
		fmt.Printf("⚠️ 清理测试数据失败: %v\n", err)
	}

	fmt.Println("=" * 50)
	fmt.Println("🎨 HTML渲染专项测试完成！")
	return nil
}

func main() {
	fmt.Println("🧪 Golang API 测试工具")
	fmt.Println("测试服务器: http://localhost:8080")
	fmt.Println()

	// 运行完整API测试
	if err := runFullTest(); err != nil {
		fmt.Printf("❌ 完整API测试失败: %v\n", err)
		return
	}

	fmt.Println()

	// 运行HTML渲染专项测试
	if err := runHTMLRenderingTest(); err != nil {
		fmt.Printf("❌ HTML渲染测试失败: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("🎉 所有测试完成！系统功能正常。")
}
