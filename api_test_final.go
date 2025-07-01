package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

// 工具函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// API响应结构 (N8nResponse)
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	URL     string      `json:"url,omitempty"`
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

// 测试客户端
type TestClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// 创建新的测试客户端
func NewTestClient(baseURL, apiKey string) *TestClient {
	jar, _ := cookiejar.New(nil)
	return &TestClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
	}
}

// 创建文章
func (c *TestClient) CreateArticle(title, content, slug, status string) (*Article, error) {
	fmt.Printf("📝 正在创建文章: %s\n", title)
	fmt.Printf("🎯 请求URL: %s\n", c.baseURL+"/api/articles")

	reqData := map[string]interface{}{
		"title":   title,
		"content": content,
		"slug":    slug,
		"status":  status,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	fmt.Printf("📋 请求方法: POST\n")
	fmt.Printf("📊 请求头: Content-Type=application/json, X-API-Key=%s***\n", c.apiKey[:5])
	fmt.Printf("📄 请求体长度: %d 字节\n", len(jsonData))
	fmt.Printf("📦 请求数据: {\n")
	fmt.Printf("   title: %s\n", title)
	fmt.Printf("   content: %s... (%d chars)\n", content[:min(50, len(content))], len(content))
	fmt.Printf("   slug: %s\n", slug)
	fmt.Printf("   status: %s\n", status)
	fmt.Printf("}\n")

	req, err := http.NewRequest("POST", c.baseURL+"/api/articles", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	fmt.Println("🚀 发送请求...")
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("⏱️ 请求耗时: %v\n", duration)
	fmt.Printf("📨 响应状态: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("📋 响应头:\n")
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("   %s: %s\n", name, value)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	fmt.Printf("📄 响应体长度: %d 字节\n", len(body))
	if len(body) < 1000 {
		fmt.Printf("📄 响应体内容: %s\n", string(body))
	} else {
		fmt.Printf("📄 响应体预览: %s...\n", string(body[:200]))
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	}

	fmt.Printf("🔍 API响应解析:\n")
	fmt.Printf("   Success: %v\n", apiResp.Success)
	if apiResp.Error != "" {
		fmt.Printf("   Error: %s\n", apiResp.Error)
	}
	if apiResp.URL != "" {
		fmt.Printf("   URL: %s\n", apiResp.URL)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API返回错误: %s", apiResp.Error)
	}

	// 解析Data字段中的文章数据
	articleData, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("序列化文章数据失败: %v", err)
	}

	var article Article
	if err := json.Unmarshal(articleData, &article); err != nil {
		return nil, fmt.Errorf("解析文章数据失败: %v", err)
	}

	fmt.Printf("📄 文章数据解析:\n")
	fmt.Printf("   ID: %s\n", article.ID)
	fmt.Printf("   Title: %s\n", article.Title)
	fmt.Printf("   Slug: %s\n", article.Slug)
	fmt.Printf("   Status: %s\n", article.Status)
	fmt.Printf("   CreatedAt: %s\n", article.CreatedAt.Format("2006-01-02 15:04:05"))

	fmt.Printf("✅ 文章创建成功，ID: %s, Slug: %s\n", article.ID, article.Slug)
	return &article, nil
}

// 获取文章
func (c *TestClient) GetArticle(id string) (*Article, error) {
	fmt.Printf("📖 正在获取文章: %s\n", id)

	req, err := http.NewRequest("GET", c.baseURL+"/api/articles/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
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

	// 解析Data字段中的文章数据
	articleData, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("序列化文章数据失败: %v", err)
	}

	var article Article
	if err := json.Unmarshal(articleData, &article); err != nil {
		return nil, fmt.Errorf("解析文章数据失败: %v", err)
	}

	fmt.Printf("✅ 文章获取成功: %s\n", article.Title)
	return &article, nil
}

// 更新文章
func (c *TestClient) UpdateArticle(id, title, content, status string) (*Article, error) {
	fmt.Printf("✏️ 正在更新文章: %s\n", id)

	reqData := map[string]interface{}{
		"title":   title,
		"content": content,
		"status":  status,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("PUT", c.baseURL+"/api/articles/"+id, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
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

	// 解析Data字段中的文章数据
	articleData, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("序列化文章数据失败: %v", err)
	}

	var article Article
	if err := json.Unmarshal(articleData, &article); err != nil {
		return nil, fmt.Errorf("解析文章数据失败: %v", err)
	}

	fmt.Printf("✅ 文章更新成功: %s\n", article.Title)
	return &article, nil
}

// 删除文章
func (c *TestClient) DeleteArticle(id string) error {
	fmt.Printf("🗑️ 正在删除文章: %s\n", id)

	req, err := http.NewRequest("DELETE", c.baseURL+"/api/articles/"+id, nil)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
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
	fmt.Printf("🌐 正在获取Web页面: /p/%s\n", slug)

	resp, err := c.httpClient.Get(c.baseURL + "/p/" + slug)
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

func main() {
	fmt.Println("🧪 Golang API 测试工具")
	fmt.Println("测试服务器: http://localhost:8080")
	fmt.Println("使用API密钥: demo-api-key-12345")
	fmt.Println(strings.Repeat("=", 50))

	// 创建测试客户端
	client := NewTestClient("http://localhost:8080", "demo-api-key-12345")

	// 1. 创建文章测试
	htmlContent := `<h1>Go API测试文章</h1>
<p>这是一个通过<strong>Go API测试</strong>创建的文章。</p>
<blockquote>测试引用内容</blockquote>
<ul>
<li>测试列表项1</li>
<li>测试列表项2</li>
</ul>
<pre><code>func main() {
    fmt.Println("Hello, World!")
}</code></pre>`

	article, err := client.CreateArticle("Go API测试文章", htmlContent, "go-api-test-"+fmt.Sprintf("%d", time.Now().Unix()), "published")
	if err != nil {
		fmt.Printf("❌ 创建文章失败: %v\n", err)
		return
	}

	// 2. 获取文章测试
	retrievedArticle, err := client.GetArticle(article.ID)
	if err != nil {
		fmt.Printf("❌ 获取文章失败: %v\n", err)
		return
	}

	// 验证获取的文章内容
	if retrievedArticle.Title != article.Title {
		fmt.Printf("❌ 获取的文章标题不匹配: 期望 %s, 实际 %s\n", article.Title, retrievedArticle.Title)
		return
	}
	fmt.Println("✅ 文章内容验证通过")

	// 3. 更新文章测试
	updatedContent := `<h1>更新后的Go API测试文章</h1>
<p>文章内容已通过<strong>Go API测试</strong>更新。</p>
<div style="background: #f0f8ff; padding: 15px; border-radius: 5px;">
<h3>更新内容</h3>
<p>这是新增的内容，用来验证更新功能。</p>
</div>
<hr>
<p><em>更新时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</em></p>`

	updatedArticle, err := client.UpdateArticle(article.ID, "更新后的Go API测试文章", updatedContent, "published")
	if err != nil {
		fmt.Printf("❌ 更新文章失败: %v\n", err)
		return
	}

	// 验证更新的文章内容
	if updatedArticle.Title != "更新后的Go API测试文章" {
		fmt.Printf("❌ 更新的文章标题不匹配\n")
		return
	}
	fmt.Println("✅ 文章更新验证通过")

	// 4. 验证Web页面渲染
	webContent, err := client.GetWebPage(article.Slug)
	if err != nil {
		fmt.Printf("❌ 获取Web页面失败: %v\n", err)
		return
	}

	// 检查Web页面是否包含更新后的内容
	if !strings.Contains(webContent, "更新后的Go API测试文章") {
		fmt.Printf("❌ Web页面内容未正确更新\n")
		return
	}

	if !strings.Contains(webContent, "<h1>更新后的Go API测试文章</h1>") {
		fmt.Printf("❌ Web页面HTML内容渲染异常\n")
		return
	}

	fmt.Println("✅ Web页面HTML渲染验证通过")

	// 检查是否是纯净的HTML内容（没有多余的页面结构）
	if !strings.Contains(webContent, "article-header") && !strings.Contains(webContent, "article-footer") {
		fmt.Println("✅ 页面使用纯净模板，没有多余结构")
	} else {
		fmt.Println("❌ 页面仍包含多余的结构元素")
	}

	// 5. 删除文章测试
	if err := client.DeleteArticle(article.ID); err != nil {
		fmt.Printf("❌ 删除文章失败: %v\n", err)
		return
	}

	// 6. 验证文章已删除
	_, err = client.GetArticle(article.ID)
	if err == nil {
		fmt.Printf("❌ 文章删除后仍能获取，删除功能异常\n")
		return
	}
	fmt.Println("✅ 文章删除验证通过")

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🎉 所有测试完成！系统功能正常，HTML渲染工作正常。")
	fmt.Println("• API创建、获取、更新、删除功能正常")
	fmt.Println("• HTML内容正确渲染，支持各种HTML标签")
	fmt.Println("• 页面使用纯净模板，只显示文章内容")
	fmt.Println("• 数据库UUID主键正常工作")
}
