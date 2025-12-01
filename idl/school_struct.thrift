namespace go school_struct
include "model.thrift"

struct GetCollegeListReq {
    required i64 page_num
    required i64 page_size
}
struct GetCollegeListResp {
     required i64 total
     required list<model.College> college_list
}

struct GetMajorListReq {
    required i64 college_id
    required i64 page_num
    required i64 page_size
}
struct GetMajorListResp {
     required i64 total
     required list<model.Major> major_list
}

struct GetTeacherListReq {
    required i64 college_id
    required i64 page_num
    required i64 page_size
}
struct GetTeacherListResp {
     required i64 total
     required list<model.Teacher> teacher_list
}

service SchoolStructService {
    GetCollegeListResp GetCollegeList(1: GetCollegeListReq req)(api.get="/api/school/college/list")
    GetMajorListResp GetMajorList(1: GetMajorListReq req)(api.get="/api/school/major/list")
    GetTeacherListResp GetTeacherList(1: GetTeacherListReq req)(api.get="/api/school/teacher/list")
}

struct AdminAddCollegeReq{
    required string college_name,
}
struct AdminAddCollegeResp{
    required model.BaseResp base_resp,
    required i64 college_id,
}


struct AdminAddMajorReq{
    required string major_name,
    required i64 college_id,
}
struct AdminAddMajorResp{
    required model.BaseResp base_resp,
    required i64 major_id,
}


struct AdminAddTeacherReq{
     required string teacher_name,
     required string college_id,
     required string introduction,
     required string email,
//     required binary avatar,
}
struct AdminAddTeacherResp{
     required model.BaseResp base_resp,
     required i64 teacher_id,
}

service AdminSchoolStructureService {
    AdminAddCollegeResp AdminAddCollege(1:AdminAddCollegeReq req)(api.post="/api/admin/colleges"),
    AdminAddMajorResp AdminAddMajor(1:AdminAddMajorReq req)(api.post="/api/admin/majors"),
    AdminAddTeacherResp AdminAddTeacher(1:AdminAddTeacherReq req)(api.post="/api/admin/teachers"),
}
