package top.wecoding.iam.server.service.impl;

import com.baomidou.mybatisplus.core.conditions.update.LambdaUpdateWrapper;
import com.baomidou.mybatisplus.core.toolkit.IdWorker;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import java.util.Date;
import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import top.wecoding.core.enums.iam.UserTypeEnum;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.core.util.PageUtil;
import top.wecoding.iam.common.enums.IamErrorCode;
import top.wecoding.iam.common.model.request.CreateUserRequest;
import top.wecoding.iam.common.model.request.DisableUserRequest;
import top.wecoding.iam.common.model.request.PasswordRequest;
import top.wecoding.iam.common.model.request.UpdateUserRequest;
import top.wecoding.iam.common.model.request.UserInfoListRequest;
import top.wecoding.iam.common.model.request.UserInfoPageRequest;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.util.AuthUtil;
import top.wecoding.iam.common.util.PasswordUtil;
import top.wecoding.iam.server.convert.UserConvert;
import top.wecoding.iam.server.enums.UserStateEnum;
import top.wecoding.iam.server.mapper.UserMapper;
import top.wecoding.iam.server.pojo.User;
import top.wecoding.iam.server.service.UserService;
import top.wecoding.iam.server.util.PasswordEncoderUtil;
import top.wecoding.mybatis.base.BaseServiceImpl;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class UserServiceImpl extends BaseServiceImpl<UserMapper, User> implements UserService {

  private final UserMapper userMapper;

  @Override
  public UserInfoResponse getInfoById(String userId) {
    User user = this.getById(userId);

    AssertUtils.isNotNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    return UserConvert.INSTANCE.toUserInfoResponse(user);
  }

  @Override
  public UserInfoResponse getInfoByUsername(String username) {
    User user = userMapper.getByUsername(username);

    AssertUtils.isNotNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    return UserConvert.INSTANCE.toUserInfoResponse(user);
  }

  @Override
  public UserInfoResponse getInfoByUsernameAndTenantId(String username, String tenantId) {

    User user = userMapper.getByTenantIdAndUsername(tenantId, username);

    AssertUtils.isNotNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    return UserConvert.INSTANCE.toUserInfoResponse(user);
  }

  @Override
  @Transactional(rollbackFor = Exception.class)
  public void create(CreateUserRequest createUserRequest) {
    String tenantId = AuthUtil.currentTenantId();
    String username = createUserRequest.getUsername();
    String password = createUserRequest.getPassword();
    String userId = IdWorker.getIdStr();

    AssertUtils.isNull(userMapper.getByUsername(username), IamErrorCode.USER_ALREADY_EXIST);

    PasswordUtil.checkPwd(password);

    User user = UserConvert.INSTANCE.toUser(createUserRequest);
    user.setTenantId(tenantId);
    user.setUserId(userId);
    user.setUserState(UserStateEnum.DEFAULT.code());
    user.setUserType(UserTypeEnum.LOCAL.code());
    user.setPassword(PasswordEncoderUtil.encode(password));
    user.setDefaultPwd(true);

    AssertUtils.isFalse(1 != userMapper.insert(user), IamErrorCode.USER_ADD_FAILED);
  }

  @Override
  @Transactional(rollbackFor = Exception.class)
  public void update(UpdateUserRequest updateUserRequest) {
    String userId = updateUserRequest.getUserId();

    User oldUser = userMapper.getByUserId(userId);
    AssertUtils.isNotNull(oldUser, IamErrorCode.USER_DOES_NOT_EXIST);

    User user = UserConvert.INSTANCE.toUser(updateUserRequest);
    user.setId(oldUser.getId());

    AssertUtils.isFalse(1 != userMapper.updateById(user), IamErrorCode.USER_UPDATE_FAILED);
  }

  @Override
  @Transactional(rollbackFor = Exception.class)
  public void delete(String userId) {
    AssertUtils.isFalse(
        Objects.equals(AuthUtil.currentUserId(), userId), IamErrorCode.CANNOT_MODIFIED_YOURSELF);

    User user = userMapper.getByUserId(userId);

    AssertUtils.isNotNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    AssertUtils.isTrue(1 == userMapper.deleteById(user.getId()), IamErrorCode.USER_DELETE_FAILED);
  }

  @Override
  public void disable(String userId, DisableUserRequest disableUserRequest) {
    Boolean disable = disableUserRequest.getDisable();

    UserStateEnum newState = disable ? UserStateEnum.DISABLE : UserStateEnum.DEFAULT;
    UserStateEnum oldState = disable ? UserStateEnum.DEFAULT : UserStateEnum.DISABLE;

    AssertUtils.isFalse(
        Objects.equals(AuthUtil.currentUserId(), userId), IamErrorCode.CANNOT_MODIFIED_YOURSELF);

    User user = userMapper.getByUserId(userId);

    AssertUtils.isNotNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    AssertUtils.isTrue(oldState.code() == user.getUserState(), IamErrorCode.USER_STATE_ABNORMAL);

    if (1
        != userMapper.updateState(
            user.getId(), newState.code(), oldState.code(), AuthUtil.currentUsername())) {
      AssertUtils.error(
          disable ? IamErrorCode.USER_DISABLE_FAILED : IamErrorCode.USER_ENABLE_FAILED);
    }
  }

  @Override
  public void password(String userId, PasswordRequest passwordRequest) {
    String newPwd = passwordRequest.getNewPwd();
    String oldPwd = passwordRequest.getOldPwd();
    boolean reset = passwordRequest.reset();

    PasswordUtil.checkPwd(newPwd);

    if (reset) {
      AssertUtils.isFalse(
          Objects.equals(AuthUtil.currentUserId(), userId), IamErrorCode.CANNOT_MODIFIED_YOURSELF);
    } else {
      AssertUtils.isFalse(Objects.equals(oldPwd, newPwd), IamErrorCode.PASSWORD_SAME_AS_OLD);
    }

    User user = userMapper.getByUserId(userId);

    AssertUtils.isNotNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    AssertUtils.isTrue(
        UserTypeEnum.LOCAL.code() == user.getUserType(),
        IamErrorCode.NON_LOCAL_USERS_DO_NOT_SUPPORT_CHANGING_PASSWORDS);

    if (!reset) {
      if (PasswordEncoderUtil.matches(newPwd, user.getPassword())) {
        AssertUtils.error(IamErrorCode.PASSWORD_SAME_AS_OLD);
      }
      if (!PasswordEncoderUtil.matches(oldPwd, user.getPassword())) {
        AssertUtils.error(IamErrorCode.OLD_PASSWORD_IS_WRONG);
      }
    }

    String encodePwd = PasswordEncoderUtil.encode(newPwd);

    LambdaUpdateWrapper<User> updateWrapper = new LambdaUpdateWrapper<>();
    updateWrapper
        .eq(User::getId, user.getId())
        .set(User::getUpdatedBy, AuthUtil.currentLoginUser().getUsername())
        .set(User::getUpdateTime, new Date())
        .set(User::getPassword, encodePwd)
        .set(User::getDefaultPwd, reset);

    AssertUtils.isTrue(this.update(updateWrapper), IamErrorCode.PASSWORD_UPDATE_FAILED);
  }

  @Override
  public PageInfo<UserInfoResponse> infoPage(UserInfoPageRequest userInfoPageRequest) {
    Page<User> pageResult =
        userMapper.page(PageUtil.getPageFromRequest(userInfoPageRequest), userInfoPageRequest);

    return PageInfo.of(pageResult.getRecords(), userInfoPageRequest, pageResult.getTotal())
        .map((UserConvert.INSTANCE::toUserInfoResponse));
  }

  @Override
  public List<UserInfoResponse> infoList(UserInfoListRequest userInfoListRequest) {
    List<User> list = userMapper.list(userInfoListRequest);
    return list.stream()
        .map((UserConvert.INSTANCE::toUserInfoResponse))
        .collect(Collectors.toList());
  }
}
