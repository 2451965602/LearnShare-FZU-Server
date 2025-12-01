package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/module"
	"LearnShare/biz/model/school_struct"
	"LearnShare/pkg/errno"
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

type SchoolStructService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewSchoolStructService(ctx context.Context, c *app.RequestContext) *SchoolStructService {
	return &SchoolStructService{ctx: ctx, c: c}
}

// GetCollegeList 获取学院列表
func (s *SchoolStructService) GetCollegeList(req *school_struct.GetCollegeListReq) ([]*module.College, int64, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	colleges, total, err := db.GetCollegeList(s.ctx, int(req.PageNum), int(req.PageSize), "福州大学")
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取学院列表失败: "+err.Error())
	}

	var collegeModules []*module.College
	for _, c := range colleges {
		collegeModules = append(collegeModules, c.ToCollegeModule())
	}

	return collegeModules, total, nil
}

// GetMajorList 获取专业列表
func (s *SchoolStructService) GetMajorList(req *school_struct.GetMajorListReq) ([]*module.Major, int64, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	majors, total, err := db.GetMajorListByCollege(s.ctx, req.CollegeID, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取专业列表失败: "+err.Error())
	}

	var majorModules []*module.Major
	for _, m := range majors {
		majorModules = append(majorModules, m.ToMajorModule())
	}

	return majorModules, total, nil
}

// GetTeacherList 获取教师列表
func (s *SchoolStructService) GetTeacherList(req *school_struct.GetTeacherListReq) ([]*module.Teacher, int64, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	teachers, total, err := db.GetTeacherListByCollegeId(s.ctx, req.CollegeID, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取教师列表失败: "+err.Error())
	}

	var teacherModules []*module.Teacher
	for _, t := range teachers {
		teacherModules = append(teacherModules, t.ToTeacherModule())
	}

	return teacherModules, total, nil
}

// AdminAddCollege 管理员添加学院
func (s *SchoolStructService) AdminAddCollege(req *school_struct.AdminAddCollegeReq) (int64, error) {
	idChan, errChan := db.AdminAddCollegeAsync(s.ctx, req.CollegeName)

	select {
	case err := <-errChan:
		if err != nil {
			return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加学院失败: "+err.Error())
		}
	case id := <-idChan:
		return id, nil
	}

	return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加学院失败")
}

// AdminAddMajor 管理员添加专业
func (s *SchoolStructService) AdminAddMajor(req *school_struct.AdminAddMajorReq) (int64, error) {
	idChan, errChan := db.AdminAddMajorAsync(s.ctx, req.MajorName, req.CollegeID)

	select {
	case err := <-errChan:
		if err != nil {
			return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加专业失败: "+err.Error())
		}
	case id := <-idChan:
		return id, nil
	}

	return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加专业失败")
}

// AdminAddTeacher 管理员添加教师
func (s *SchoolStructService) AdminAddTeacher(req *school_struct.AdminAddTeacherReq) (int64, error) {
	// 将 college_id 从 string 转换为 int64
	collegeID, err := strconv.ParseInt(req.CollegeID, 10, 64)
	if err != nil {
		return 0, errno.NewErrNo(errno.ParamVerifyErrorCode, "无效的学院ID")
	}

	idChan, errChan := db.AdminAddTeacherAsync(s.ctx, req.TeacherName, collegeID, req.Introduction, req.Email)

	select {
	case err := <-errChan:
		if err != nil {
			return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加教师失败: "+err.Error())
		}
	case id := <-idChan:
		return id, nil
	}

	return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加教师失败")
}
