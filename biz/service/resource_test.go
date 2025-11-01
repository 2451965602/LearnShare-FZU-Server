package service

import (
	"context"
	"testing"
	"time"

	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/resource"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ResourceSearchTestSuite 资源搜索测试套件
type ResourceSearchTestSuite struct {
	suite.Suite
	db      *gorm.DB
	service *SearchResourcesService
	ctx     context.Context
}

// SetupSuite 测试套件初始化
func (suite *ResourceSearchTestSuite) SetupSuite() {
	// 初始化测试数据库连接
	suite.ctx = context.Background()

	// 使用测试数据库配置 - 修正密码和数据库名
	dsn := "root:123456@tcp(localhost:3306)/learnshare_test?charset=utf8mb4&parseTime=True&loc=Local"
	testDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，与主配置保持一致
		},
	})
	if err != nil {
		suite.T().Skipf("跳过测试：无法连接到测试数据库: %v", err)
		return
	}

	suite.db = testDB
	db.DB = testDB // 设置全局数据库连接

	// 自动迁移表结构
	err = suite.db.AutoMigrate(
		&db.Resource{},
		&db.ResourceTag{},
		&db.ResourceTagMapping{},
	)
	suite.Require().NoError(err)

	// 初始化服务
	suite.service = NewSearchResourcesService(suite.ctx)
}

// TearDownSuite 测试套件清理
func (suite *ResourceSearchTestSuite) TearDownSuite() {
	if suite.db != nil {
		// 清理测试数据 - 按正确顺序删除表（先删除有外键的表）
		suite.db.Exec("DROP TABLE IF EXISTS resource_tag_mapping")
		suite.db.Exec("DROP TABLE IF EXISTS resource_tag")
		suite.db.Exec("DROP TABLE IF EXISTS resource")

		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

// SetupTest 每个测试前的准备
func (suite *ResourceSearchTestSuite) SetupTest() {
	// 清理测试数据 - 使用单数表名
	suite.db.Exec("DELETE FROM resource_tag_mapping")
	suite.db.Exec("DELETE FROM resource")
	suite.db.Exec("DELETE FROM resource_tag")

	// 插入测试数据
	suite.insertTestData()
}

// insertTestData 插入测试数据
func (suite *ResourceSearchTestSuite) insertTestData() {
	// 插入标签数据
	tags := []db.ResourceTag{
		{TagID: 1, TagName: "算法"},
		{TagID: 2, TagName: "数据结构"},
		{TagID: 3, TagName: "机器学习"},
		{TagID: 4, TagName: "深度学习"},
		{TagID: 5, TagName: "Python"},
	}
	for _, tag := range tags {
		suite.db.Create(&tag)
	}

	// 插入资源数据
	resources := []db.Resource{
		{
			ResourceID:    1,
			Title:         "算法导论第四版",
			Description:   "经典算法教材，涵盖各种算法和数据结构",
			FilePath:      "/files/algorithm_intro.pdf",
			FileType:      "pdf",
			FileSize:      10485760, // 10MB
			UploaderID:    1,
			CourseID:      1,
			DownloadCount: 150,
			AverageRating: 4.8,
			RatingCount:   25,
			Status:        1,                                   // 已发布
			CreatedAt:     time.Now().Add(-7 * 24 * time.Hour), // 7天前
		},
		{
			ResourceID:    2,
			Title:         "Python机器学习实战",
			Description:   "Python机器学习入门教程，包含大量实例",
			FilePath:      "/files/python_ml.pdf",
			FileType:      "pdf",
			FileSize:      8388608, // 8MB
			UploaderID:    2,
			CourseID:      2,
			DownloadCount: 89,
			AverageRating: 4.2,
			RatingCount:   18,
			Status:        1,
			CreatedAt:     time.Now().Add(-3 * 24 * time.Hour), // 3天前
		},
		{
			ResourceID:    3,
			Title:         "深度学习基础",
			Description:   "深度学习理论基础和实践指南",
			FilePath:      "/files/deep_learning.docx",
			FileType:      "docx",
			FileSize:      5242880, // 5MB
			UploaderID:    1,
			CourseID:      3,
			DownloadCount: 200,
			AverageRating: 4.9,
			RatingCount:   35,
			Status:        1,
			CreatedAt:     time.Now().Add(-1 * 24 * time.Hour), // 1天前
		},
		{
			ResourceID:    4,
			Title:         "数据结构与算法分析",
			Description:   "C++语言描述的数据结构教材",
			FilePath:      "/files/data_structure.pdf",
			FileType:      "pdf",
			FileSize:      12582912, // 12MB
			UploaderID:    3,
			CourseID:      1,
			DownloadCount: 75,
			AverageRating: 4.0,
			RatingCount:   12,
			Status:        1,
			CreatedAt:     time.Now().Add(-5 * 24 * time.Hour), // 5天前
		},
		{
			ResourceID:    5,
			Title:         "Python编程从入门到实践",
			Description:   "Python基础编程教程",
			FilePath:      "/files/python_basic.pdf",
			FileType:      "pdf",
			FileSize:      6291456, // 6MB
			UploaderID:    2,
			CourseID:      2,
			DownloadCount: 120,
			AverageRating: 4.5,
			RatingCount:   20,
			Status:        1,
			CreatedAt:     time.Now().Add(-2 * 24 * time.Hour), // 2天前
		},
	}

	for _, res := range resources {
		suite.db.Create(&res)
	}

	// 插入资源标签关联数据
	mappings := []db.ResourceTagMapping{
		{ResourceID: 1, TagID: 1}, // 算法导论 - 算法
		{ResourceID: 1, TagID: 2}, // 算法导论 - 数据结构
		{ResourceID: 2, TagID: 3}, // Python机器学习 - 机器学习
		{ResourceID: 2, TagID: 5}, // Python机器学习 - Python
		{ResourceID: 3, TagID: 4}, // 深度学习基础 - 深度学习
		{ResourceID: 3, TagID: 3}, // 深度学习基础 - 机器学习
		{ResourceID: 4, TagID: 2}, // 数据结构与算法 - 数据结构
		{ResourceID: 4, TagID: 1}, // 数据结构与算法 - 算法
		{ResourceID: 5, TagID: 5}, // Python编程 - Python
	}

	for _, mapping := range mappings {
		suite.db.Create(&mapping)
	}
}

// TestSearchByKeyword 测试关键词搜索
func (suite *ResourceSearchTestSuite) TestSearchByKeyword() {
	// 测试搜索"算法"
	keyword := "算法"
	req := &resource.SearchResourceReq{
		Keyword:  &keyword,
		PageSize: 10,
		PageNum:  1,
	}

	results, total, err := suite.service.SearchResources(req)

	suite.Assert().NoError(err)
	suite.Assert().Equal(int64(2), total) // 应该找到2个包含"算法"的资源
	suite.Assert().Len(results, 2)

	// 验证结果包含正确的资源
	titles := make([]string, len(results))
	for i, res := range results {
		titles[i] = res.Title
	}
	suite.Assert().Contains(titles, "算法导论第四版")
	suite.Assert().Contains(titles, "数据结构与算法分析")
}

// TestSearchByTag 测试标签搜索
func (suite *ResourceSearchTestSuite) TestSearchByTag() {
	// 测试搜索标签ID为5的资源（Python）
	tagID := int64(5)
	req := &resource.SearchResourceReq{
		TagId:    &tagID,
		PageSize: 10,
		PageNum:  1,
	}

	results, total, err := suite.service.SearchResources(req)

	suite.Assert().NoError(err)
	suite.Assert().Equal(int64(2), total) // 应该找到2个Python相关的资源
	suite.Assert().Len(results, 2)

	// 验证结果包含正确的资源
	titles := make([]string, len(results))
	for i, res := range results {
		titles[i] = res.Title
	}
	suite.Assert().Contains(titles, "Python机器学习实战")
	suite.Assert().Contains(titles, "Python编程从入门到实践")
}

// TestSearchByCourse 测试课程筛选
func (suite *ResourceSearchTestSuite) TestSearchByCourse() {
	// 测试搜索课程ID为1的资源
	courseID := int64(1)
	req := &resource.SearchResourceReq{
		CourseId: &courseID,
		PageSize: 10,
		PageNum:  1,
	}

	results, total, err := suite.service.SearchResources(req)

	suite.Assert().NoError(err)
	suite.Assert().Equal(int64(2), total) // 应该找到2个课程1的资源
	suite.Assert().Len(results, 2)

	// 验证所有结果都属于课程1
	for _, res := range results {
		suite.Assert().Equal(int64(1), res.CourseId)
	}
}

// TestSortByDownloadCount 测试按下载量排序
func (suite *ResourceSearchTestSuite) TestSortByDownloadCount() {
	sortBy := "hot"
	req := &resource.SearchResourceReq{
		SortBy:   &sortBy,
		PageSize: 10,
		PageNum:  1,
	}

	results, total, err := suite.service.SearchResources(req)

	suite.Assert().NoError(err)
	suite.Assert().Equal(int64(5), total)
	suite.Assert().Len(results, 5)

	// 验证按下载量降序排列
	suite.Assert().True(results[0].DownloadCount >= results[1].DownloadCount)
	suite.Assert().True(results[1].DownloadCount >= results[2].DownloadCount)

	// 第一个应该是下载量最高的"深度学习基础"(200次)
	suite.Assert().Equal("深度学习基础", results[0].Title)
	suite.Assert().Equal(int64(200), results[0].DownloadCount)
}

// TestSortByRating 测试按评分排序
func (suite *ResourceSearchTestSuite) TestSortByRating() {
	sortBy := "rating"
	req := &resource.SearchResourceReq{
		SortBy:   &sortBy,
		PageSize: 10,
		PageNum:  1,
	}

	results, total, err := suite.service.SearchResources(req)

	suite.Assert().NoError(err)
	suite.Assert().Equal(int64(5), total)
	suite.Assert().Len(results, 5)

	// 验证按评分降序排列
	suite.Assert().True(results[0].AverageRating >= results[1].AverageRating)
	suite.Assert().True(results[1].AverageRating >= results[2].AverageRating)

	// 第一个应该是评分最高的"深度学习基础"(4.9分)
	suite.Assert().Equal("深度学习基础", results[0].Title)
	suite.Assert().Equal(4.9, results[0].AverageRating)
}

// TestSortByLatest 测试按时间排序（默认）
func (suite *ResourceSearchTestSuite) TestSortByLatest() {
	req := &resource.SearchResourceReq{
		PageSize: 10,
		PageNum:  1,
	}

	results, total, err := suite.service.SearchResources(req)

	suite.Assert().NoError(err)
	suite.Assert().Equal(int64(5), total)
	suite.Assert().Len(results, 5)

	// 验证按创建时间降序排列（最新的在前）
	suite.Assert().True(results[0].CreatedAt >= results[1].CreatedAt)
	suite.Assert().True(results[1].CreatedAt >= results[2].CreatedAt)

	// 第一个应该是最新的"深度学习基础"
	suite.Assert().Equal("深度学习基础", results[0].Title)
}

// TestPagination 测试分页功能
func (suite *ResourceSearchTestSuite) TestPagination() {
	// 测试第一页
	req1 := &resource.SearchResourceReq{
		PageSize: 2,
		PageNum:  1,
	}

	results1, total1, err1 := suite.service.SearchResources(req1)

	suite.Assert().NoError(err1)
	suite.Assert().Equal(int64(5), total1)
	suite.Assert().Len(results1, 2)

	// 测试第二页
	req2 := &resource.SearchResourceReq{
		PageSize: 2,
		PageNum:  2,
	}

	results2, total2, err2 := suite.service.SearchResources(req2)

	suite.Assert().NoError(err2)
	suite.Assert().Equal(int64(5), total2)
	suite.Assert().Len(results2, 2)

	// 验证两页的结果不同
	suite.Assert().NotEqual(results1[0].ResourceId, results2[0].ResourceId)
	suite.Assert().NotEqual(results1[1].ResourceId, results2[1].ResourceId)
}

// TestCombinedSearch 测试组合搜索
func (suite *ResourceSearchTestSuite) TestCombinedSearch() {
	// 测试关键词+标签+排序的组合搜索
	keyword := "Python"
	tagID := int64(5)
	sortBy := "rating"

	req := &resource.SearchResourceReq{
		Keyword:  &keyword,
		TagId:    &tagID,
		SortBy:   &sortBy,
		PageSize: 10,
		PageNum:  1,
	}

	results, total, err := suite.service.SearchResources(req)

	suite.Assert().NoError(err)
	suite.Assert().Equal(int64(2), total) // 应该找到2个同时包含Python关键词和Python标签的资源
	suite.Assert().Len(results, 2)

	// 验证按评分排序
	suite.Assert().True(results[0].AverageRating >= results[1].AverageRating)
}

// TestEmptyResult 测试空结果
func (suite *ResourceSearchTestSuite) TestEmptyResult() {
	keyword := "不存在的关键词"
	req := &resource.SearchResourceReq{
		Keyword:  &keyword,
		PageSize: 10,
		PageNum:  1,
	}

	results, total, err := suite.service.SearchResources(req)

	suite.Assert().NoError(err)
	suite.Assert().Equal(int64(0), total)
	suite.Assert().Len(results, 0)
}

// TestResourceTags 测试资源标签关联
func (suite *ResourceSearchTestSuite) TestResourceTags() {
	req := &resource.SearchResourceReq{
		PageSize: 1,
		PageNum:  1,
	}

	results, _, err := suite.service.SearchResources(req)

	suite.Assert().NoError(err)
	suite.Assert().Len(results, 1)

	// 验证第一个资源有标签
	if len(results) > 0 && results[0].Tags != nil {
		suite.Assert().True(len(results[0].Tags) > 0)

		// 验证标签结构
		for _, tag := range results[0].Tags {
			suite.Assert().NotZero(tag.TagId)
			suite.Assert().NotEmpty(tag.TagName)
		}
	}
}

// TestInvalidParameters 测试无效参数
func (suite *ResourceSearchTestSuite) TestInvalidParameters() {
	// 测试无效的页码
	req := &resource.SearchResourceReq{
		PageSize: 10,
		PageNum:  0, // 无效页码
	}

	results, total, err := suite.service.SearchResources(req)

	// 应该返回错误或空结果
	suite.Assert().True(err != nil || (total == 0 && len(results) == 0))
}

// 运行测试套件
func TestResourceSearchSuite(t *testing.T) {
	suite.Run(t, new(ResourceSearchTestSuite))
}

// 基准测试
func BenchmarkSearchResources(b *testing.B) {
	// 跳过基准测试，因为需要数据库连接
	b.Skip("跳过基准测试：需要数据库连接")

	// 初始化测试环境
	ctx := context.Background()
	service := NewSearchResourcesService(ctx)

	req := &resource.SearchResourceReq{
		PageSize: 10,
		PageNum:  1,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = service.SearchResources(req)
	}
}
