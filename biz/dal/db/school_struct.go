package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
)

// GetCollegeList 获取学院列表
func GetCollegeList(ctx context.Context, pageNum, pageSize int, school string) ([]*College, int64, error) {
	var colleges []*College
	var total int64

	// 获取分页数据
	err := DB.WithContext(ctx).Table(constants.CollegeTableName).
		Where("school = ?", school).
		Order("college_id").
		Count(&total).
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&colleges).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询学院列表失败: "+err.Error())
	}

	return colleges, total, nil
}

// GetMajorListByCollege 根据学院ID获取专业列表
func GetMajorListByCollege(ctx context.Context, collegeID int64, pageNum, pageSize int) ([]*Major, int64, error) {
	var majors []*Major
	var total int64

	// 获取分页数据
	err := DB.WithContext(ctx).Table(constants.MajorTableName).
		Where("college_id = ?", collegeID).
		Order("major_id").
		Count(&total).
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&majors).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询专业列表失败: "+err.Error())
	}

	return majors, total, nil
}

// GetTeacherListByCollegeId 根据专业ID获取教师列表
func GetTeacherListByCollegeId(ctx context.Context, collegeId int64, pageNum, pageSize int) ([]*Teacher, int64, error) {
	var teachers []*Teacher
	var total int64

	// 获取分页数据
	err := DB.WithContext(ctx).Table(constants.TeacherTableName).
		Where("college_id = ?", collegeId).
		Order("teacher_id").
		Count(&total).
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&teachers).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询教师列表失败: "+err.Error())
	}

	return teachers, total, nil
}

func GetTeacherByID(ctx context.Context, teacherID int64) (*Teacher, error) {
	{
		var teacher Teacher

		err := DB.WithContext(ctx).Table(constants.TeacherTableName).
			Where("teacher_id = ?", teacherID).
			First(&teacher).Error
		if err != nil {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询教师失败: "+err.Error())
		}

		return &teacher, nil
	}
}

// AdminAddCollege 管理员添加学院
func AdminAddCollege(ctx context.Context, collegeName string) (int64, error) {
	college := &College{
		CollegeName: collegeName,
		School:      "福州大学",
	}

	err := DB.WithContext(ctx).Table(constants.CollegeTableName).Create(college).Error
	if err != nil {
		return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加学院失败: "+err.Error())
	}

	return college.CollegeID, nil
}

// AdminAddCollegeAsync 管理员添加学院（异步）
func AdminAddCollegeAsync(ctx context.Context, collegeName string) (chan int64, chan error) {
	idChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	pool := GetAsyncPool()
	pool.Submit(func() error {
		id, err := AdminAddCollege(ctx, collegeName)
		if err != nil {
			errChan <- err
			close(idChan)
			close(errChan)
			return err
		}
		idChan <- id
		close(idChan)
		close(errChan)
		return nil
	})

	return idChan, errChan
}

// AdminAddMajor 管理员添加专业
func AdminAddMajor(ctx context.Context, majorName string, collegeID int64) (int64, error) {
	major := &Major{
		MajorName: majorName,
		CollegeID: collegeID,
	}

	err := DB.WithContext(ctx).Table(constants.MajorTableName).Create(major).Error
	if err != nil {
		return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加专业失败: "+err.Error())
	}

	return major.MajorID, nil
}

// AdminAddMajorAsync 管理员添加专业（异步）
func AdminAddMajorAsync(ctx context.Context, majorName string, collegeID int64) (chan int64, chan error) {
	idChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	pool := GetAsyncPool()
	pool.Submit(func() error {
		id, err := AdminAddMajor(ctx, majorName, collegeID)
		if err != nil {
			errChan <- err
			close(idChan)
			close(errChan)
			return err
		}
		idChan <- id
		close(idChan)
		close(errChan)
		return nil
	})

	return idChan, errChan
}

// AdminAddTeacher 管理员添加教师
func AdminAddTeacher(ctx context.Context, teacherName string, collegeID int64, introduction, email string) (int64, error) {
	teacher := &Teacher{
		Name:         teacherName,
		CollegeID:    &collegeID,
		Introduction: &introduction,
		Email:        &email,
	}

	err := DB.WithContext(ctx).Table(constants.TeacherTableName).Create(teacher).Error
	if err != nil {
		return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加教师失败: "+err.Error())
	}

	return teacher.TeacherID, nil
}

// AdminAddTeacherAsync 管理员添加教师（异步）
func AdminAddTeacherAsync(ctx context.Context, teacherName string, collegeID int64, introduction, email string) (chan int64, chan error) {
	idChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	pool := GetAsyncPool()
	pool.Submit(func() error {
		id, err := AdminAddTeacher(ctx, teacherName, collegeID, introduction, email)
		if err != nil {
			errChan <- err
			close(idChan)
			close(errChan)
			return err
		}
		idChan <- id
		close(idChan)
		close(errChan)
		return nil
	})

	return idChan, errChan
}
