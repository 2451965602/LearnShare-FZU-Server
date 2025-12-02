namespace go user
include "model.thrift"

//用户注册
struct RegisterReq {
  required string username;
  required string password;
  required string email;
}
struct RegisterResp {
  required model.BaseResp baseResponse;
}

//用户登录
struct LoginInReq {
  required string email;
  required string password;
}
struct LoginInResp {
  required model.BaseResp baseResponse;
  optional model.User user;
}

//用户登出
struct LoginOutReq {
}
struct LoginOutResp {
  required model.BaseResp baseResponse;
}

//获取邮箱验证码
struct SendVerifyEmailReq {
  required string email;
}
struct SendVerifyEmailResp {
  required model.BaseResp baseResponse;
}

//验证邮箱验证码
struct VerifyEmailReq {
  required string email;
  required string code;
}
struct VerifyEmailResp {
  required model.BaseResp baseResponse;
}

//修改邮箱
struct updateEmailReq {
  required string new_email;
  required string code;
}
struct updateEmailResp {
  required model.BaseResp baseResponse;
}

//修改密码
struct UpdatePasswordReq {
  required string old_password;
  required string new_password;
}
struct UpdatePasswordResp {
  required model.BaseResp baseResponse;
}

//修改专业
struct updateMajorReq {
  required i64 new_majorId;
}
struct updateMajorResp {
  required model.BaseResp baseResponse;
}


//上传头像
struct uploadAvatarReq {
//  required binary avatar;
}
struct uploadAvatarResp {
  required model.BaseResp baseResponse;
}


//重置密码
struct ResetPasswordReq {
  required string email;
  required string newPassword;
  required string code;
}
struct ResetPasswordResp {
  required model.BaseResp baseResponse;
}


//刷新Token
struct RefreshTokenReq {
}
struct RefreshTokenResp {
  required model.BaseResp baseResponse;
}

//获取用户信息
struct GetUserInfoReq {
  required i64 user_id (api.path="user_id");
}
struct GetUserInfoResp {
  required model.BaseResp baseResponse;
  optional model.User user;
}



service UserService {
  RegisterResp register(1: RegisterReq req)(api.post="/api/auth/register"),
  LoginInResp loginIn(1: LoginInReq req)(api.post="/api/auth/login"),
  LoginOutResp loginOut(1: LoginOutReq req)(api.post="/api/auth/logout"),
  SendVerifyEmailResp sendVerifyEmail(1: SendVerifyEmailReq req)(api.post="/api/users/me/email/get"),
  VerifyEmailResp verifyEmail(1: VerifyEmailReq req)(api.post="/api/users/me/email/verify"),
  updateEmailResp updateEmail(1: updateEmailReq req)(api.put="/api/users/me/email"),
  UpdatePasswordResp updatePassword(1: UpdatePasswordReq req)(api.put="/api/users/me/password"),
  updateMajorResp updateMajor(1: updateMajorReq req)(api.put="/api/users/me/major"),
  uploadAvatarResp uploadAvatar(1: uploadAvatarReq req)(api.put="/api/users/avatar"),
  ResetPasswordResp resetPassword(1: ResetPasswordReq req)(api.post="/api/users/me/password/reset"),
  RefreshTokenResp refreshToken(1: RefreshTokenReq req)(api.post="/api/auth/refresh"),
  GetUserInfoResp getUserInfo(1: GetUserInfoReq req)(api.get="/api/users/:user_id"),
}

struct AdminAddUserReq{
    required string username,
    required string password,
    required string email,
    required i64 role_id,
    required string status,
}
struct AdminAddUserResp{
    required model.BaseResp base_resp,
    required i64 user_id,
}

struct AdminUpdateUserReq{
    required i64 user_id,
    optional string username,
    optional string password,
    optional string email,
    optional string college_id,
    optional string major_id,
    optional binary avatar,
    optional i64 reputation_score,
    optional i64 role_id,
    optional string status,
}
struct AdminUpdateUserResp{
    required model.BaseResp base_resp,
}

struct AdminGetUserListReq{
    required i64 page_size
    required i64 page_num
}

struct AdminGetUserListResp{
    required model.BaseResp base_resp,
    required list<model.User> user_list,
    required i64 total,
}

service UserAdminService {
    AdminAddUserResp AdminAddUser(1:AdminAddUserReq req)(api.post="/api/admin/users"),
    AdminUpdateUserResp AdminUpdateUser(1:AdminUpdateUserReq req)(api.put="/api/admin/users"),
    AdminGetUserListResp AdminGetUserList(1:AdminGetUserListReq req)(api.get="/api/admin/users"),
}

struct GetPermissionListReq{
}
struct GetPermissionListResp{
    required model.BaseResp base_resp,
    required list<model.Permission> permission_list,
}

struct GetRoleListReq{
}
struct GetRoleListResp{
    required model.BaseResp base_resp,
    required list<model.Role> role_list,
}

struct AddRoleReq{
    required string role_name,
    required list<i64> permission_ids,
}
struct AddRoleResp{
    required model.BaseResp base_resp,
    required i64 role_id,
}

service RoleAdminService {
    GetPermissionListResp GetPermissionList(1:GetPermissionListReq req)(api.get="/api/admin/permissions"),
    GetRoleListResp GetRoleList(1:GetRoleListReq req)(api.get="/api/admin/roles"),
    AddRoleResp AddRole(1:AddRoleReq req)(api.post="/api/admin/roles"),
}
