package top.wecoding.iam.server.service.impl;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.iam.common.enums.IamErrorCode;
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

    AssertUtils.isNotNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

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
    User user = userMapper.getByUsername(username);

    AssertUtils.isNotNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    UserInfo userInfo =
        UserInfo.builder()
            .userId(user.getUserId())
            .tenantId(user.getTenantId())
            .tenantName("wecoding")
            .username(user.getUsername())
            .nickName(user.getNickName())
            .password(user.getPassword())
            .avatar(user.getAvatar())
            .birthday(user.getBirthday())
            .gender(user.getGender())
            .email(user.getEmail())
            .phone(user.getPhone())
            .lastLoginIp(user.getLastLoginIp())
            .lastLoginTime(user.getLastLoginTime())
            .userType(user.getUserType())
            .userState(user.getUserState())
            .defaultPwd(user.getDefaultPwd())
            .infos(user.getInfos())
            .build();

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
