package top.wecoding.iam.server.service;

import java.util.List;
import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.common.model.request.*;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.server.pojo.User;
import top.wecoding.mybatis.base.BaseService;

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

  void update(UpdateUserRequest updateUserRequest);

  void delete(String userId);

  void disable(String userId, DisableUserRequest disableUserRequest);

  void password(String userId, PasswordRequest passwordRequest);

  PageInfo<UserInfoResponse> infoPage(UserInfoPageRequest userInfoPageRequest);

  List<UserInfoResponse> infoList(UserInfoListRequest userInfoListRequest);
}
