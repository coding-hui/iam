package top.wecoding.iam.server.service;

import top.wecoding.iam.common.model.request.CreateUserRequest;
import top.wecoding.iam.common.model.request.UpdateUserRequest;
import top.wecoding.iam.server.entity.UserProfile;
import top.wecoding.mybatis.base.BaseService;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface UserProfileService extends BaseService<UserProfile> {

  boolean create(String userId, CreateUserRequest createUserRequest);

  boolean update(String userId, UpdateUserRequest updateUserRequest);

  boolean delete(String userId);
}
