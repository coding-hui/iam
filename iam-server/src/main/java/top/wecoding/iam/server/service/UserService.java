package top.wecoding.iam.server.service;

import top.wecoding.commons.core.model.PageInfo;
import top.wecoding.iam.common.model.request.*;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.userdetails.LoginUser;
import top.wecoding.iam.server.entity.User;
import top.wecoding.mybatis.base.BaseService;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
public interface UserService extends BaseService<User> {

  UserInfoResponse getInfo(LoginUser loginUser);

  UserInfoResponse getInfoById(String userId);

  UserInfoResponse getInfoByUsername(String username);

  void create(CreateUserRequest createUserRequest);

  void update(String userId, UpdateUserRequest updateUserRequest);

  void delete(String userId);

  void disable(String userId, DisableUserRequest disableUserRequest);

  void password(String userId, PasswordRequest passwordRequest);

  PageInfo<UserInfoResponse> infoPage(UserInfoPageRequest userInfoPageRequest);
}
