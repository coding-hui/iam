package top.wecoding.iam.server.service.impl;

import com.baomidou.mybatisplus.core.toolkit.Wrappers;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import top.wecoding.core.enums.rest.CommonErrorCodeEnum;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.iam.common.model.UserInfo;
import top.wecoding.iam.common.model.request.*;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.server.mapper.UserMapper;
import top.wecoding.iam.server.pojo.User;
import top.wecoding.iam.server.service.UserService;
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
  public UserInfoResponse getInfoById(String userId) {
    User user = this.getById(userId);

    AssertUtils.isNotNull(user, CommonErrorCodeEnum.COMMON_ERROR, "user does not exist");

    UserInfo userInfo =
        UserInfo.builder()
            .tenantId(user.getTenantId())
            .userId(userId)
            .username(user.getUsername())
            .nickName(user.getNickName())
            .avatar(user.getAvatar())
            .build();

    return UserInfoResponse.builder().userInfo(userInfo).build();
  }

  @Override
  public UserInfoResponse getInfoByUsername(String username) {
    User user = userMapper.selectOne(Wrappers.<User>lambdaQuery().eq(User::getUsername, username));

    AssertUtils.isNotNull(user, CommonErrorCodeEnum.COMMON_ERROR, "user does not exist");

    UserInfo userInfo =
        UserInfo.builder().userId(user.getUserId()).username(user.getUsername()).build();

    return UserInfoResponse.builder().userInfo(userInfo).build();
  }

  @Override
  public UserInfoResponse getInfoByUsernameAndTenantId(String username, String tenantId) {
    return null;
  }

  @Override
  public void create(CreateUserRequest createUserRequest) {
    User user =
        User.builder()
            .username(createUserRequest.getUsername())
            .password(createUserRequest.getPassword())
            .build();

    userMapper.insert(user);
  }

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
