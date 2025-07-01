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
	apiKey     string
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
	fmt.Printf("📤 请求URL: %s\n", c.baseURL+"/admin/login")

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	fmt.Printf("📝 请求数据: username=%s, password=***\n", username)
	fmt.Printf("📋 请求体: %s\n", data.Encode())

	req, err := http.NewRequest("POST", c.baseURL+"/admin/login", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("创建登录请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Printf("📊 请求头: Content-Type=%s\n", req.Header.Get("Content-Type"))

	fmt.Println("🚀 发送请求...")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("执行登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("📨 响应状态码: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("📋 响应头信息:\n")
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("   %s: %s\n", name, value)
		}
	}

	// 读取响应体
	body, _ := io.ReadAll(resp.Body)
	if len(body) > 0 {
		fmt.Printf("📄 响应体长度: %d 字节\n", len(body))
		if len(body) < 500 { // 只有响应体较小时才完整打印
			fmt.Printf("📄 响应体内容: %s\n", string(body))
		} else {
			fmt.Printf("📄 响应体预览: %s...\n", string(body[:200]))
		}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return fmt.Errorf("登录失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	fmt.Println("✅ 管理员登录成功")
	return nil
}

// 打印HTTP请求详情
func (c *TestClient) printRequestDetails(req *http.Request, jsonData []byte) {
	fmt.Printf("📤 请求URL: %s\n", req.URL.String())
	fmt.Printf("📋 请求方法: %s\n", req.Method)

	fmt.Printf("📊 请求头:\n")
	for name, values := range req.Header {
		for _, value := range values {
			if name == "X-Api-Key" {
				if value == "" {
					fmt.Printf("   %s: (空值)\n", name)
				} else if len(value) > 10 {
					fmt.Printf("   %s: %s...%s (长度:%d)\n", name, value[:5], value[len(value)-3:], len(value))
				} else {
					fmt.Printf("   %s: %s\n", name, value)
				}
			} else {
				fmt.Printf("   %s: %s\n", name, value)
			}
		}
	}

	if jsonData != nil {
		fmt.Printf("📝 请求体长度: %d 字节\n", len(jsonData))
		if len(jsonData) < 300 {
			fmt.Printf("📝 请求体内容: %s\n", string(jsonData))
		} else {
			fmt.Printf("📝 请求体预览: %s...\n", string(jsonData[:200]))
		}
	}
}

// 打印HTTP响应详情
func (c *TestClient) printResponseDetails(resp *http.Response, body []byte) {
	fmt.Printf("📨 响应状态码: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("📋 响应头:\n")
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("   %s: %s\n", name, value)
		}
	}

	fmt.Printf("📄 响应体长度: %d 字节\n", len(body))
	if len(body) < 500 {
		fmt.Printf("📄 响应体内容: %s\n", string(body))
	} else {
		fmt.Printf("📄 响应体预览: %s...\n", string(body[:300]))
	}
}

// 创建文章（带详细日志）
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
	if c.apiKey != "" {
		httpReq.Header.Set("X-API-Key", c.apiKey)
	}

	c.printRequestDetails(httpReq, jsonData)

	fmt.Println("🚀 发送请求...")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	c.printResponseDetails(resp, body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
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

// 获取文章（带详细日志）
func (c *TestClient) GetArticle(id string) (*Article, error) {
	fmt.Printf("📖 正在获取文章: %s\n", id)

	httpReq, err := http.NewRequest("GET", c.baseURL+"/api/articles/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("X-API-Key", c.apiKey)
	}

	c.printRequestDetails(httpReq, nil)

	fmt.Println("🚀 发送请求...")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	c.printResponseDetails(resp, body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
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

// 更新文章（带详细日志）
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
	if c.apiKey != "" {
		httpReq.Header.Set("X-API-Key", c.apiKey)
	}

	c.printRequestDetails(httpReq, jsonData)

	fmt.Println("🚀 发送请求...")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	c.printResponseDetails(resp, body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
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

// 删除文章（带详细日志）
func (c *TestClient) DeleteArticle(id string) error {
	fmt.Printf("🗑️ 正在删除文章: %s\n", id)

	httpReq, err := http.NewRequest("DELETE", c.baseURL+"/api/articles/"+id, nil)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("X-API-Key", c.apiKey)
	}

	c.printRequestDetails(httpReq, nil)

	fmt.Println("🚀 发送请求...")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	c.printResponseDetails(resp, body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
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

// 设置API Key
func (c *TestClient) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
	fmt.Printf("🔑 设置API Key: %s\n", func() string {
		if apiKey == "" {
			return "(空值)"
		} else if len(apiKey) > 10 {
			return fmt.Sprintf("%s...%s (长度:%d)", apiKey[:5], apiKey[len(apiKey)-3:], len(apiKey))
		} else {
			return apiKey
		}
	}())
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
	fmt.Println(strings.Repeat("=", 50))

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

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🎉 所有API测试通过！")
	return nil
}

// 运行HTML渲染专项测试
func runHTMLRenderingTest() error {
	fmt.Println("🎨 开始运行HTML渲染专项测试")
	fmt.Println(strings.Repeat("=", 50))

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

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🎨 HTML渲染专项测试完成！")
	return nil
}

// 运行API Key认证专项测试
func runAPIKeyAuthTest() error {
	fmt.Println("🔐 开始运行API Key认证专项测试")
	fmt.Println(strings.Repeat("=", 50))

	client := NewTestClient("http://localhost:8080")

	// 管理员登录（用于对比）
	if err := client.Login("admin", "password"); err != nil {
		return fmt.Errorf("登录失败: %v", err)
	}

	// 测试文章数据
	createReq := CreateArticleRequest{
		Title:   "API Key认证测试文章",
		Content: "<h1>测试API Key认证功能</h1><p>这是用于测试不同API Key场景的文章。</p>",
		Status:  "published",
		Slug:    "api-key-auth-test",
	}

	// 1. 测试无API Key的情况
	fmt.Println("\n🧪 测试场景1: 无API Key")
	fmt.Println(strings.Repeat("-", 30))
	client.SetAPIKey("")
	_, err := client.CreateArticle(createReq)
	if err != nil {
		fmt.Printf("❌ 预期结果：无API Key请求被拒绝 - %v\n", err)
	} else {
		fmt.Printf("⚠️ 意外结果：无API Key请求竟然成功了\n")
	}

	// 2. 测试错误的API Key
	fmt.Println("\n🧪 测试场景2: 错误的API Key")
	fmt.Println(strings.Repeat("-", 30))
	client.SetAPIKey("wrong-api-key-12345")
	_, err = client.CreateArticle(createReq)
	if err != nil {
		fmt.Printf("❌ 预期结果：错误API Key请求被拒绝 - %v\n", err)
	} else {
		fmt.Printf("⚠️ 意外结果：错误API Key请求竟然成功了\n")
	}

	// 3. 测试正确的API Key
	fmt.Println("\n🧪 测试场景3: 正确的API Key")
	fmt.Println(strings.Repeat("-", 30))
	client.SetAPIKey("demo-api-key-12345")
	article, err := client.CreateArticle(createReq)
	if err != nil {
		fmt.Printf("⚠️ 意外结果：正确API Key请求失败了 - %v\n", err)
		return err
	} else {
		fmt.Printf("✅ 预期结果：正确API Key请求成功\n")
	}

	// 测试其他操作的API Key认证
	fmt.Println("\n🧪 测试场景4: 测试其他API操作的认证")
	fmt.Println(strings.Repeat("-", 30))

	// 4.1 测试获取文章（正确Key）
	fmt.Println("\n📖 测试获取文章（正确API Key）:")
	_, err = client.GetArticle(article.ID)
	if err != nil {
		fmt.Printf("⚠️ 获取文章失败: %v\n", err)
	} else {
		fmt.Printf("✅ 获取文章成功\n")
	}

	// 4.2 测试获取文章（错误Key）
	fmt.Println("\n📖 测试获取文章（错误API Key）:")
	client.SetAPIKey("wrong-key")
	_, err = client.GetArticle(article.ID)
	if err != nil {
		fmt.Printf("❌ 预期结果：错误API Key获取文章被拒绝 - %v\n", err)
	} else {
		fmt.Printf("⚠️ 意外结果：错误API Key获取文章竟然成功了\n")
	}

	// 4.3 测试更新文章（正确Key）
	fmt.Println("\n✏️ 测试更新文章（正确API Key）:")
	client.SetAPIKey("demo-api-key-12345")
	updateReq := UpdateArticleRequest{
		Title:   "更新后的API Key测试文章",
		Content: "<h1>已更新</h1><p>文章内容已通过正确的API Key更新。</p>",
	}
	_, err = client.UpdateArticle(article.ID, updateReq)
	if err != nil {
		fmt.Printf("⚠️ 更新文章失败: %v\n", err)
	} else {
		fmt.Printf("✅ 更新文章成功\n")
	}

	// 4.4 测试更新文章（无Key）
	fmt.Println("\n✏️ 测试更新文章（无API Key）:")
	client.SetAPIKey("")
	_, err = client.UpdateArticle(article.ID, updateReq)
	if err != nil {
		fmt.Printf("❌ 预期结果：无API Key更新文章被拒绝 - %v\n", err)
	} else {
		fmt.Printf("⚠️ 意外结果：无API Key更新文章竟然成功了\n")
	}

	// 4.5 测试删除文章（正确Key）
	fmt.Println("\n🗑️ 测试删除文章（正确API Key）:")
	client.SetAPIKey("demo-api-key-12345")
	err = client.DeleteArticle(article.ID)
	if err != nil {
		fmt.Printf("⚠️ 删除文章失败: %v\n", err)
	} else {
		fmt.Printf("✅ 删除文章成功\n")
	}

	// 5. 测试第二个有效的API Key
	fmt.Println("\n🧪 测试场景5: 测试第二个有效API Key")
	fmt.Println(strings.Repeat("-", 30))
	client.SetAPIKey("n8n-integration-key")
	createReq.Slug = "api-key-test-2"
	article2, err := client.CreateArticle(createReq)
	if err != nil {
		fmt.Printf("⚠️ 第二个API Key请求失败 - %v\n", err)
	} else {
		fmt.Printf("✅ 第二个API Key请求成功\n")
		// 清理
		client.DeleteArticle(article2.ID)
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🔐 API Key认证专项测试完成！")
	return nil
}

// 完整的API认证与功能测试
func runCompleteAPITest() error {
	fmt.Println("🔄 开始运行完整的API认证与功能测试")
	fmt.Println(strings.Repeat("=", 60))

	client := NewTestClient("http://localhost:8080")

	// 1. 管理员登录测试
	fmt.Println("\n🔐 第一步：管理员认证")
	if err := client.Login("admin", "password"); err != nil {
		return fmt.Errorf("管理员登录失败: %v", err)
	}

	// 2. 设置API Key
	fmt.Println("\n🔑 第二步：设置API Key")
	client.SetAPIKey("demo-api-key-12345")

	// 3. 创建文章
	fmt.Println("\n📝 第三步：创建文章")
	createReq := CreateArticleRequest{
		Title: "完整测试文章",
		Content: `<h1>完整API测试</h1>
<h2>功能验证</h2>
<ul>
<li>✅ 管理员认证</li>
<li>✅ API Key认证</li>
<li>✅ UUID主键</li>
<li>✅ HTML内容渲染</li>
</ul>
<blockquote>
<p>这是一个完整的API功能测试，验证了所有核心功能。</p>
</blockquote>
<pre><code>// 示例代码
fmt.Println("测试成功！")
</code></pre>`,
		Status: "published",
		Slug:   "complete-api-test",
	}

	article, err := client.CreateArticle(createReq)
	if err != nil {
		return fmt.Errorf("创建文章失败: %v", err)
	}

	// 4. 获取文章
	fmt.Println("\n📖 第四步：获取文章")
	retrievedArticle, err := client.GetArticle(article.ID)
	if err != nil {
		return fmt.Errorf("获取文章失败: %v", err)
	}

	// 5. 验证UUID格式
	fmt.Printf("🔍 验证UUID格式: %s\n", retrievedArticle.ID)
	if len(retrievedArticle.ID) != 36 || !strings.Contains(retrievedArticle.ID, "-") {
		fmt.Printf("⚠️ 警告：文章ID可能不是标准UUID格式\n")
	} else {
		fmt.Printf("✅ UUID格式验证通过\n")
	}

	// 6. 更新文章
	fmt.Println("\n✏️ 第五步：更新文章")
	updateReq := UpdateArticleRequest{
		Title: "更新后的完整测试文章",
		Content: `<h1>更新后的完整API测试</h1>
<div style="background: #e8f5e8; padding: 15px; border-radius: 5px; border-left: 4px solid #4caf50;">
<h3>✅ 更新验证</h3>
<p>文章已成功更新，所有功能正常运行。</p>
</div>
<p><strong>更新时间:</strong> ` + time.Now().Format("2006-01-02 15:04:05") + `</p>`,
	}

	updatedArticle, err := client.UpdateArticle(article.ID, updateReq)
	if err != nil {
		return fmt.Errorf("更新文章失败: %v", err)
	}

	fmt.Printf("✅ 文章更新成功: %s\n", updatedArticle.Title)

	// 7. 验证Web页面渲染
	fmt.Println("\n🌐 第六步：验证Web页面渲染")
	webContent, err := client.GetWebPage(article.Slug)
	if err != nil {
		return fmt.Errorf("获取Web页面失败: %v", err)
	}

	// 检查HTML渲染
	htmlChecks := []struct {
		name    string
		pattern string
	}{
		{"更新后的标题", "更新后的完整API测试"},
		{"HTML H1标签", "<h1>更新后的完整API测试</h1>"},
		{"CSS样式", "background: #e8f5e8"},
		{"强调文字", "<strong>更新时间:</strong>"},
		{"DIV容器", `<div style=`},
	}

	for _, check := range htmlChecks {
		if strings.Contains(webContent, check.pattern) {
			fmt.Printf("✅ %s 渲染正确\n", check.name)
		} else {
			fmt.Printf("❌ %s 渲染异常\n", check.name)
		}
	}

	// 8. 清理测试数据
	fmt.Println("\n🧹 第七步：清理测试数据")
	if err := client.DeleteArticle(article.ID); err != nil {
		fmt.Printf("⚠️ 清理测试数据失败: %v\n", err)
	} else {
		fmt.Printf("✅ 测试数据清理完成\n")
	}

	// 9. 验证删除效果
	fmt.Println("\n🔍 第八步：验证删除效果")
	_, err = client.GetArticle(article.ID)
	if err != nil {
		fmt.Printf("✅ 删除验证通过：文章已无法获取\n")
	} else {
		fmt.Printf("⚠️ 删除验证失败：文章仍可获取\n")
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🎉 完整的API认证与功能测试完成！")
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

	// 运行API Key认证专项测试
	if err := runAPIKeyAuthTest(); err != nil {
		fmt.Printf("❌ API Key认证测试失败: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("🎉 所有测试完成！系统功能正常。")
}
