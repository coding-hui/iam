package top.wecoding.iam.service.impl;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.mapper.UserMapper;
import top.wecoding.iam.model.request.*;
import top.wecoding.iam.model.response.UserInfoResponse;
import top.wecoding.iam.pojo.User;
import top.wecoding.iam.service.UserService;
import top.wecoding.mybatis.base.BaseServiceImpl;

import javax.annotation.Resource;
import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Slf4j
@Service
public class UserServiceImpl extends BaseServiceImpl<UserMapper, User> implements UserService {

  @Resource private UserMapper userMapper;

  @Override
  public UserInfoResponse getInfo(String userId) {
    return null;
  }

  @Override
  public void create(CreateUserRequest createUserRequest) {}

  @Override
  public void delete(DeleteUserRequest deleteUserRequest) {}

  @Override
  public void disable(String userId, DisableUserRequest disableUserRequest) {}

  @Override
  public void password(String userId, PasswordRequest passwordRequest) {}

  @Override
  public PageInfo<UserInfoResponse> infoPage(UserInfoPageRequest userInfoPageRequest) {
    return null;
  }

  @Override
  public List<UserInfoResponse> infoList(UserInfoPageRequest userInfoPageRequest) {
    return null;
  }
}
