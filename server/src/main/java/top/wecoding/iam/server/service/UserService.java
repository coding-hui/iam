package top.wecoding.iam.server.service;

import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.common.model.request.*;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.server.pojo.User;
import top.wecoding.mybatis.base.BaseService;

import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface UserService extends BaseService<User> {

  UserInfoResponse getInfoById(String userId);

  UserInfoResponse getInfoByUsername(String username);

  UserInfoResponse getInfoByUsernameAndTenantId(String username, String tenantId);

  void create(CreateUserRequest createUserRequest);

  void delete(DeleteUserRequest deleteUserRequest);

  void disable(String userId, DisableUserRequest disableUserRequest);

  void password(String userId, PasswordRequest passwordRequest);

  PageInfo<UserInfoResponse> infoPage(UserInfoPageRequest userInfoPageRequest);

  List<UserInfoResponse> infoList(UserInfoPageRequest userInfoPageRequest);
}