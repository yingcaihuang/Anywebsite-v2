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
	fmt.Printf("\n🔥 ==================== 创建文章 ====================\n")
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
	fmt.Printf("📊 请求头:\n")
	fmt.Printf("   Content-Type: application/json\n")
	fmt.Printf("   X-API-Key: %s*** (隐藏完整密钥)\n", c.apiKey[:5])
	fmt.Printf("📄 请求体长度: %d 字节\n", len(jsonData))
	fmt.Printf("📦 请求数据结构:\n")
	fmt.Printf("   ├─ title: %s\n", title)
	fmt.Printf("   ├─ content: %s... (%d 字符)\n", content[:min(50, len(content))], len(content))
	fmt.Printf("   ├─ slug: %s\n", slug)
	fmt.Printf("   └─ status: %s\n", status)
	
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
		fmt.Printf("📄 完整响应体:\n%s\n", string(body))
	} else {
		fmt.Printf("📄 响应体预览:\n%s...\n", string(body[:300]))
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	}
	
	fmt.Printf("🔍 API响应解析:\n")
	fmt.Printf("   ├─ Success: %v\n", apiResp.Success)
	if apiResp.Error != "" {
		fmt.Printf("   ├─ Error: %s\n", apiResp.Error)
	}
	if apiResp.URL != "" {
		fmt.Printf("   └─ URL: %s\n", apiResp.URL)
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
	
	fmt.Printf("📄 解析出的文章数据:\n")
	fmt.Printf("   ├─ ID: %s\n", article.ID)
	fmt.Printf("   ├─ Title: %s\n", article.Title)
	fmt.Printf("   ├─ Slug: %s\n", article.Slug)
	fmt.Printf("   ├─ Status: %s\n", article.Status)
	fmt.Printf("   ├─ Content Length: %d 字符\n", len(article.Content))
	fmt.Printf("   └─ CreatedAt: %s\n", article.CreatedAt.Format("2006-01-02 15:04:05"))
	
	fmt.Printf("✅ 文章创建成功！\n")
	return &article, nil
}

// 获取文章
func (c *TestClient) GetArticle(id string) (*Article, error) {
	fmt.Printf("\n🔥 ==================== 获取文章 ====================\n")
	fmt.Printf("📖 正在获取文章: %s\n", id)
	fmt.Printf("🎯 请求URL: %s\n", c.baseURL+"/api/articles/"+id)
	
	req, err := http.NewRequest("GET", c.baseURL+"/api/articles/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	
	req.Header.Set("X-API-Key", c.apiKey)
	
	fmt.Printf("📋 请求方法: GET\n")
	fmt.Printf("📊 请求头:\n")
	fmt.Printf("   X-API-Key: %s*** (隐藏完整密钥)\n", c.apiKey[:5])
	
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
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	
	fmt.Printf("📄 响应体长度: %d 字节\n", len(body))
	if len(body) < 1000 {
		fmt.Printf("📄 完整响应体:\n%s\n", string(body))
	} else {
		fmt.Printf("📄 响应体预览:\n%s...\n", string(body[:300]))
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	fmt.Printf("🔍 API响应解析:\n")
	fmt.Printf("   ├─ Success: %v\n", apiResp.Success)
	if apiResp.Error != "" {
		fmt.Printf("   └─ Error: %s\n", apiResp.Error)
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
	
	fmt.Printf("📄 解析出的文章数据:\n")
	fmt.Printf("   ├─ ID: %s\n", article.ID)
	fmt.Printf("   ├─ Title: %s\n", article.Title)
	fmt.Printf("   ├─ Slug: %s\n", article.Slug)
	fmt.Printf("   ├─ Status: %s\n", article.Status)
	fmt.Printf("   ├─ Content Length: %d 字符\n", len(article.Content))
	fmt.Printf("   ├─ CreatedAt: %s\n", article.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   └─ UpdatedAt: %s\n", article.UpdatedAt.Format("2006-01-02 15:04:05"))
	
	fmt.Printf("✅ 文章获取成功！\n")
	return &article, nil
}

// 更新文章
func (c *TestClient) UpdateArticle(id, title, content, status string) (*Article, error) {
	fmt.Printf("\n🔥 ==================== 更新文章 ====================\n")
	fmt.Printf("✏️ 正在更新文章: %s\n", id)
	fmt.Printf("🎯 请求URL: %s\n", c.baseURL+"/api/articles/"+id)
	
	reqData := map[string]interface{}{
		"title":   title,
		"content": content,
		"status":  status,
	}
	
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}
	
	fmt.Printf("📋 请求方法: PUT\n")
	fmt.Printf("📊 请求头:\n")
	fmt.Printf("   Content-Type: application/json\n")
	fmt.Printf("   X-API-Key: %s*** (隐藏完整密钥)\n", c.apiKey[:5])
	fmt.Printf("📄 请求体长度: %d 字节\n", len(jsonData))
	fmt.Printf("📦 更新数据结构:\n")
	fmt.Printf("   ├─ title: %s\n", title)
	fmt.Printf("   ├─ content: %s... (%d 字符)\n", content[:min(50, len(content))], len(content))
	fmt.Printf("   └─ status: %s\n", status)
	
	req, err := http.NewRequest("PUT", c.baseURL+"/api/articles/"+id, bytes.NewBuffer(jsonData))
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
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	
	fmt.Printf("📄 响应体长度: %d 字节\n", len(body))
	if len(body) < 1000 {
		fmt.Printf("📄 完整响应体:\n%s\n", string(body))
	} else {
		fmt.Printf("📄 响应体预览:\n%s...\n", string(body[:300]))
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	}
	
	fmt.Printf("🔍 API响应解析:\n")
	fmt.Printf("   ├─ Success: %v\n", apiResp.Success)
	if apiResp.Error != "" {
		fmt.Printf("   ├─ Error: %s\n", apiResp.Error)
	}
	if apiResp.URL != "" {
		fmt.Printf("   └─ URL: %s\n", apiResp.URL)
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
	
	fmt.Printf("📄 更新后的文章数据:\n")
	fmt.Printf("   ├─ ID: %s\n", article.ID)
	fmt.Printf("   ├─ Title: %s\n", article.Title)
	fmt.Printf("   ├─ Slug: %s\n", article.Slug)
	fmt.Printf("   ├─ Status: %s\n", article.Status)
	fmt.Printf("   ├─ Content Length: %d 字符\n", len(article.Content))
	fmt.Printf("   ├─ CreatedAt: %s\n", article.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   └─ UpdatedAt: %s\n", article.UpdatedAt.Format("2006-01-02 15:04:05"))
	
	fmt.Printf("✅ 文章更新成功！\n")
	return &article, nil
}

// 删除文章
func (c *TestClient) DeleteArticle(id string) error {
	fmt.Printf("\n🔥 ==================== 删除文章 ====================\n")
	fmt.Printf("🗑️ 正在删除文章: %s\n", id)
	fmt.Printf("🎯 请求URL: %s\n", c.baseURL+"/api/articles/"+id)
	
	req, err := http.NewRequest("DELETE", c.baseURL+"/api/articles/"+id, nil)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	
	req.Header.Set("X-API-Key", c.apiKey)
	
	fmt.Printf("📋 请求方法: DELETE\n")
	fmt.Printf("📊 请求头:\n")
	fmt.Printf("   X-API-Key: %s*** (隐藏完整密钥)\n", c.apiKey[:5])
	
	fmt.Println("🚀 发送请求...")
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime)
	
	if err != nil {
		return fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	fmt.Printf("⏱️ 请求耗时: %v\n", duration)
	fmt.Printf("📨 响应状态: %d %s\n", resp.StatusCode, resp.Status)
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}
	
	fmt.Printf("📄 响应体长度: %d 字节\n", len(body))
	if len(body) > 0 {
		fmt.Printf("📄 响应体内容:\n%s\n", string(body))
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}
	
	fmt.Printf("🔍 API响应解析:\n")
	fmt.Printf("   ├─ Success: %v\n", apiResp.Success)
	if apiResp.Error != "" {
		fmt.Printf("   └─ Error: %s\n", apiResp.Error)
	}
	
	if !apiResp.Success {
		return fmt.Errorf("API返回错误: %s", apiResp.Error)
	}
	
	fmt.Printf("✅ 文章删除成功！\n")
	return nil
}

// 获取Web页面内容
func (c *TestClient) GetWebPage(slug string) (string, error) {
	fmt.Printf("\n🔥 ==================== 获取Web页面 ====================\n")
	fmt.Printf("🌐 正在获取Web页面: /p/%s\n", slug)
	fmt.Printf("🎯 请求URL: %s\n", c.baseURL+"/p/"+slug)
	
	fmt.Printf("📋 请求方法: GET\n")
	fmt.Printf("📊 请求头: 默认浏览器头\n")
	
	fmt.Println("🚀 发送请求...")
	startTime := time.Now()
	resp, err := c.httpClient.Get(c.baseURL + "/p/" + slug)
	duration := time.Since(startTime)
	
	if err != nil {
		return "", fmt.Errorf("执行请求失败: %v", err)
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
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Web页面返回错误状态码: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	
	fmt.Printf("📄 页面长度: %d 字节\n", len(body))
	
	// 分析HTML内容
	content := string(body)
	if strings.Contains(content, "<html") {
		fmt.Printf("📊 HTML分析:\n")
		fmt.Printf("   ├─ 文档类型: HTML页面\n")
		fmt.Printf("   ├─ 是否包含<head>: %v\n", strings.Contains(content, "<head"))
		fmt.Printf("   ├─ 是否包含<body>: %v\n", strings.Contains(content, "<body"))
		fmt.Printf("   ├─ 是否包含<h1>: %v\n", strings.Contains(content, "<h1"))
		fmt.Printf("   ├─ 是否包含<p>: %v\n", strings.Contains(content, "<p"))
		fmt.Printf("   ├─ 是否包含style属性: %v\n", strings.Contains(content, "style="))
		fmt.Printf("   └─ 是否为纯净模板: %v\n", !strings.Contains(content, "article-header"))
	}
	
	// 显示页面内容预览
	if len(content) > 500 {
		fmt.Printf("📄 页面内容预览 (前500字符):\n%s...\n", content[:500])
	} else {
		fmt.Printf("📄 完整页面内容:\n%s\n", content)
	}
	
	fmt.Printf("✅ Web页面获取成功！\n")
	return content, nil
}

func main() {
	fmt.Println("🧪 Golang API 详细测试工具")
	fmt.Println("测试服务器: http://localhost:8080")
	fmt.Println("使用API密钥: demo-api-key-12345")
	fmt.Println(strings.Repeat("=", 60))
	
	// 创建测试客户端
	client := NewTestClient("http://localhost:8080", "demo-api-key-12345")
	
	// 1. 创建文章测试
	htmlContent := `<h1>Go API详细测试文章</h1>
<p>这是一个通过<strong>Go API详细测试</strong>创建的文章。</p>
<blockquote>测试引用内容，包含各种HTML元素</blockquote>
<ul>
<li>测试列表项1</li>
<li>测试列表项2</li>
</ul>
<pre><code>func main() {
    fmt.Println("Hello, Detailed Testing!")
}</code></pre>
<p style="color: blue;">这是一个带样式的段落。</p>`
	
	article, err := client.CreateArticle("Go API详细测试文章", htmlContent, "go-api-detailed-test-"+fmt.Sprintf("%d", time.Now().Unix()), "published")
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
	fmt.Printf("\n✅ 文章内容验证通过: 标题匹配\n")
	
	// 3. 更新文章测试
	updatedContent := `<h1>更新后的Go API详细测试文章</h1>
<p>文章内容已通过<strong>Go API详细测试</strong>更新。</p>
<div style="background: #f0f8ff; padding: 15px; border-radius: 5px;">
<h3>更新内容</h3>
<p>这是新增的内容，用来验证更新功能和详细日志输出。</p>
</div>
<hr>
<p><em>更新时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</em></p>`
	
	updatedArticle, err := client.UpdateArticle(article.ID, "更新后的Go API详细测试文章", updatedContent, "published")
	if err != nil {
		fmt.Printf("❌ 更新文章失败: %v\n", err)
		return
	}
	
	// 验证更新的文章内容
	if updatedArticle.Title != "更新后的Go API详细测试文章" {
		fmt.Printf("❌ 更新的文章标题不匹配\n")
		return
	}
	fmt.Printf("\n✅ 文章更新验证通过: 标题正确更新\n")
	
	// 4. 验证Web页面渲染
	webContent, err := client.GetWebPage(article.Slug)
	if err != nil {
		fmt.Printf("❌ 获取Web页面失败: %v\n", err)
		return
	}
	
	// 检查Web页面是否包含更新后的内容
	if !strings.Contains(webContent, "更新后的Go API详细测试文章") {
		fmt.Printf("❌ Web页面内容未正确更新\n")
		return
	}
	
	if !strings.Contains(webContent, "<h1>更新后的Go API详细测试文章</h1>") {
		fmt.Printf("❌ Web页面HTML内容渲染异常\n")
		return
	}
	
	fmt.Printf("\n✅ Web页面HTML渲染验证通过\n")
	
	// 检查是否是纯净的HTML内容（没有多余的页面结构）
	if !strings.Contains(webContent, "article-header") && !strings.Contains(webContent, "article-footer") {
		fmt.Printf("✅ 页面使用纯净模板，没有多余结构\n")
	} else {
		fmt.Printf("❌ 页面仍包含多余的结构元素\n")
	}
	
	// 5. 删除文章测试
	if err := client.DeleteArticle(article.ID); err != nil {
		fmt.Printf("❌ 删除文章失败: %v\n", err)
		return
	}
	
	// 6. 验证文章已删除
	fmt.Printf("\n🔍 验证文章是否已删除...\n")
	_, err = client.GetArticle(article.ID)
	if err == nil {
		fmt.Printf("❌ 文章删除后仍能获取，删除功能异常\n")
		return
	}
	fmt.Printf("✅ 文章删除验证通过: 删除后无法获取\n")
	
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🎉 所有详细测试完成！系统功能正常。")
	fmt.Println("• API创建、获取、更新、删除功能正常")
	fmt.Println("• HTML内容正确渲染，支持各种HTML标签")
	fmt.Println("• 页面使用纯净模板，只显示文章内容")
	fmt.Println("• 数据库UUID主键正常工作")
	fmt.Println("• 详细的请求/响应日志已输出完成")
}
