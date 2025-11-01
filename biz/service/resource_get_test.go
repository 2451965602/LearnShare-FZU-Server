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

// ResourceGetTestSuite 获取单个资源信息测试套件
type ResourceGetTestSuite struct {
	suite.Suite
	db      *gorm.DB
	service *GetResourceService
	ctx     context.Context
}

// SetupSuite 测试套件初始化
func (suite *ResourceGetTestSuite) SetupSuite() {
	// 初始化测试数据库连接
	suite.ctx = context.Background()

	// 使用测试数据库配置
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
	suite.service = NewGetResourceService(suite.ctx)
}

// TearDownSuite 测试套件清理
func (suite *ResourceGetTestSuite) TearDownSuite() {
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
func (suite *ResourceGetTestSuite) SetupTest() {
	// 清理测试数据 - 使用单数表名
	suite.db.Exec("DELETE FROM resource_tag_mapping")
	suite.db.Exec("DELETE FROM resource")
	suite.db.Exec("DELETE FROM resource_tag")

	// 插入测试数据
	suite.insertTestData()
}

// insertTestData 插入测试数据
func (suite *ResourceGetTestSuite) insertTestData() {
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
			ResourceID:    1001,
			Title:         "算法导论第四版",
			Description:   "经典算法教材，涵盖各种算法和数据结构",
			FilePath:      "/files/algorithm_intro.pdf",
			FileType:      "pdf",
			FileSize:      10485760, // 10MB
			UploaderID:    2001,
			CourseID:      1001,
			DownloadCount: 150,
			AverageRating: 4.8,
			RatingCount:   25,
			Status:        1,                                   // 已发布
			CreatedAt:     time.Now().Add(-7 * 24 * time.Hour), // 7天前
		},
		{
			ResourceID:    1002,
			Title:         "Python机器学习实战",
			Description:   "Python机器学习入门教程，包含大量实例",
			FilePath:      "/files/python_ml.pdf",
			FileType:      "pdf",
			FileSize:      8388608, // 8MB
			UploaderID:    2002,
			CourseID:      1002,
			DownloadCount: 89,
			AverageRating: 4.2,
			RatingCount:   18,
			Status:        1,
			CreatedAt:     time.Now().Add(-3 * 24 * time.Hour), // 3天前
		},
		{
			ResourceID:    1003,
			Title:         "深度学习基础",
			Description:   "深度学习理论基础和实践指南",
			FilePath:      "/files/deep_learning.docx",
			FileType:      "docx",
			FileSize:      5242880, // 5MB
			UploaderID:    2001,
			CourseID:      1003,
			DownloadCount: 200,
			AverageRating: 4.9,
			RatingCount:   35,
			Status:        1,
			CreatedAt:     time.Now().Add(-1 * 24 * time.Hour), // 1天前
		},
	}

	for _, res := range resources {
		suite.db.Create(&res)
	}

	// 插入资源标签关联数据
	mappings := []db.ResourceTagMapping{
		{ResourceID: 1001, TagID: 1}, // 算法导论 - 算法
		{ResourceID: 1001, TagID: 2}, // 算法导论 - 数据结构
		{ResourceID: 1002, TagID: 3}, // Python机器学习 - 机器学习
		{ResourceID: 1002, TagID: 5}, // Python机器学习 - Python
		{ResourceID: 1003, TagID: 4}, // 深度学习基础 - 深度学习
		{ResourceID: 1003, TagID: 3}, // 深度学习基础 - 机器学习
	}

	for _, mapping := range mappings {
		suite.db.Create(&mapping)
	}
}

// TestGetExistingResource 测试获取存在的资源
func (suite *ResourceGetTestSuite) TestGetExistingResource() {
	// 测试获取资源ID为1001的资源
	req := &resource.GetResourceReq{
		ResourceId: 1001,
	}

	result, err := suite.service.GetResource(req)

	suite.Assert().NoError(err)
	suite.Assert().NotNil(result)

	// 验证资源基本信息
	suite.Assert().Equal(int64(1001), result.ResourceId)
	suite.Assert().Equal("算法导论第四版", result.Title)
	suite.Assert().Equal("经典算法教材，涵盖各种算法和数据结构", *result.Description)
	suite.Assert().Equal("/files/algorithm_intro.pdf", result.FilePath)
	suite.Assert().Equal("pdf", result.FileType)
	suite.Assert().Equal(int64(10485760), result.FileSize)
	suite.Assert().Equal(int64(2001), result.UploaderId)
	suite.Assert().Equal(int64(1001), result.CourseId)
	suite.Assert().Equal(int64(150), result.DownloadCount)
	suite.Assert().Equal(4.8, result.AverageRating)
	suite.Assert().Equal(int64(25), result.RatingCount)
	suite.Assert().Equal(int32(1), result.Status)

	// 验证标签信息
	suite.Assert().NotNil(result.Tags)
	suite.Assert().Len(result.Tags, 2)

	// 验证标签内容
	tagNames := make([]string, len(result.Tags))
	for i, tag := range result.Tags {
		tagNames[i] = tag.TagName
	}
	suite.Assert().Contains(tagNames, "算法")
	suite.Assert().Contains(tagNames, "数据结构")
}

// TestGetResourceWithMultipleTags 测试获取有多个标签的资源
func (suite *ResourceGetTestSuite) TestGetResourceWithMultipleTags() {
	// 测试获取资源ID为1003的资源（深度学习基础）
	req := &resource.GetResourceReq{
		ResourceId: 1003,
	}

	result, err := suite.service.GetResource(req)

	suite.Assert().NoError(err)
	suite.Assert().NotNil(result)

	// 验证资源基本信息
	suite.Assert().Equal(int64(1003), result.ResourceId)
	suite.Assert().Equal("深度学习基础", result.Title)
	suite.Assert().Equal(int64(200), result.DownloadCount)
	suite.Assert().Equal(4.9, result.AverageRating)

	// 验证标签信息
	suite.Assert().NotNil(result.Tags)
	suite.Assert().Len(result.Tags, 2)

	// 验证标签内容
	tagNames := make([]string, len(result.Tags))
	for i, tag := range result.Tags {
		tagNames[i] = tag.TagName
	}
	suite.Assert().Contains(tagNames, "深度学习")
	suite.Assert().Contains(tagNames, "机器学习")
}

// TestGetNonExistingResource 测试获取不存在的资源
func (suite *ResourceGetTestSuite) TestGetNonExistingResource() {
	// 测试获取不存在的资源ID
	req := &resource.GetResourceReq{
		ResourceId: 9999, // 不存在的资源ID
	}

	result, err := suite.service.GetResource(req)

	// 应该返回错误
	suite.Assert().Error(err)
	suite.Assert().Nil(result)
}

// TestGetResourceWithNoTags 测试获取没有标签的资源
func (suite *ResourceGetTestSuite) TestGetResourceWithNoTags() {
	// 创建一个没有标签的资源
	resourceWithoutTags := db.Resource{
		ResourceID:    1004,
		Title:         "无标签测试资源",
		Description:   "这是一个没有标签的测试资源",
		FilePath:      "/files/no_tags.pdf",
		FileType:      "pdf",
		FileSize:      5242880, // 5MB
		UploaderID:    2001,
		CourseID:      1001,
		DownloadCount: 10,
		AverageRating: 4.0,
		RatingCount:   5,
		Status:        1,
		CreatedAt:     time.Now().Add(-1 * 24 * time.Hour),
	}
	suite.db.Create(&resourceWithoutTags)

	// 测试获取这个资源
	req := &resource.GetResourceReq{
		ResourceId: 1004,
	}

	result, err := suite.service.GetResource(req)

	suite.Assert().NoError(err)
	suite.Assert().NotNil(result)

	// 验证资源基本信息
	suite.Assert().Equal(int64(1004), result.ResourceId)
	suite.Assert().Equal("无标签测试资源", result.Title)

	// 验证标签信息 - 应该为空切片而不是nil
	suite.Assert().NotNil(result.Tags)
	suite.Assert().Len(result.Tags, 0)
}

// TestGetResourceWithInvalidID 测试使用无效的资源ID
func (suite *ResourceGetTestSuite) TestGetResourceWithInvalidID() {
	// 测试使用0作为资源ID
	req := &resource.GetResourceReq{
		ResourceId: 0, // 无效的资源ID
	}

	result, err := suite.service.GetResource(req)

	// 应该返回错误
	suite.Assert().Error(err)
	suite.Assert().Nil(result)
}

// TestGetResourceFieldsConsistency 测试资源字段一致性
func (suite *ResourceGetTestSuite) TestGetResourceFieldsConsistency() {
	// 测试获取Python机器学习资源
	req := &resource.GetResourceReq{
		ResourceId: 1002,
	}

	result, err := suite.service.GetResource(req)

	suite.Assert().NoError(err)
	suite.Assert().NotNil(result)

	// 验证所有字段都正确填充
	suite.Assert().Equal(int64(1002), result.ResourceId)
	suite.Assert().Equal("Python机器学习实战", result.Title)
	suite.Assert().Equal("Python机器学习入门教程，包含大量实例", *result.Description)
	suite.Assert().Equal("/files/python_ml.pdf", result.FilePath)
	suite.Assert().Equal("pdf", result.FileType)
	suite.Assert().Equal(int64(8388608), result.FileSize)
	suite.Assert().Equal(int64(2002), result.UploaderId)
	suite.Assert().Equal(int64(1002), result.CourseId)
	suite.Assert().Equal(int64(89), result.DownloadCount)
	suite.Assert().Equal(4.2, result.AverageRating)
	suite.Assert().Equal(int64(18), result.RatingCount)
	suite.Assert().Equal(int32(1), result.Status)
	suite.Assert().True(result.CreatedAt > 0) // 创建时间应该有效

	// 验证标签
	suite.Assert().Len(result.Tags, 2)
	tagNames := make([]string, len(result.Tags))
	for i, tag := range result.Tags {
		tagNames[i] = tag.TagName
	}
	suite.Assert().Contains(tagNames, "机器学习")
	suite.Assert().Contains(tagNames, "Python")
}

// 运行测试套件
func TestResourceGetSuite(t *testing.T) {
	suite.Run(t, new(ResourceGetTestSuite))
}
