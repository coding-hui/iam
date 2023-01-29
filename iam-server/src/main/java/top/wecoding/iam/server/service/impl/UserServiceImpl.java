package top.wecoding.iam.server.service.impl;

import com.baomidou.mybatisplus.core.conditions.update.LambdaUpdateWrapper;
import com.baomidou.mybatisplus.core.toolkit.IdWorker;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import java.util.Date;
import java.util.Objects;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import top.wecoding.commons.core.model.PageInfo;
import top.wecoding.commons.core.util.ArgumentAssert;
import top.wecoding.iam.common.entity.UserInfo;
import top.wecoding.iam.common.enums.IamErrorCode;
import top.wecoding.iam.common.model.request.CreateUserRequest;
import top.wecoding.iam.common.model.request.DisableUserRequest;
import top.wecoding.iam.common.model.request.PasswordRequest;
import top.wecoding.iam.common.model.request.UpdateUserRequest;
import top.wecoding.iam.common.model.request.UserInfoPageRequest;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.userdetails.LoginUser;
import top.wecoding.iam.common.util.AuthUtil;
import top.wecoding.iam.common.util.PasswordUtil;
import top.wecoding.iam.server.convert.UserConvert;
import top.wecoding.iam.server.entity.User;
import top.wecoding.iam.server.enums.UserStateEnum;
import top.wecoding.iam.server.mapper.UserMapper;
import top.wecoding.iam.server.service.UserProfileService;
import top.wecoding.iam.server.service.UserService;
import top.wecoding.iam.server.util.PasswordEncoderUtil;
import top.wecoding.mybatis.base.BaseServiceImpl;
import top.wecoding.mybatis.helper.PageHelper;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class UserServiceImpl extends BaseServiceImpl<UserMapper, User> implements UserService {

  private final UserProfileService userProfileService;

  @Override
  public UserInfoResponse getInfo(LoginUser loginUser) {
    return UserInfoResponse.builder()
        .userInfo(loginUser.userInfo())
        .groups(loginUser.groups())
        .permissions(loginUser.permissions())
        .roles(loginUser.roles())
        .build();
  }

  @Override
  public UserInfoResponse getInfoById(String userId) {
    UserInfo userInfo = baseMapper.getInfoById(userId);
    ArgumentAssert.notNull(userInfo, IamErrorCode.USER_DOES_NOT_EXIST);

    return UserInfoResponse.builder().userInfo(userInfo).build();
  }

  @Override
  public UserInfoResponse getInfoByUsername(String username) {
    UserInfo userInfo = baseMapper.getInfoByUsername(username);
    ArgumentAssert.notNull(userInfo, IamErrorCode.USER_DOES_NOT_EXIST);

    return UserInfoResponse.builder().userInfo(userInfo).build();
  }

  @Override
  @Transactional(rollbackFor = Exception.class)
  public void create(CreateUserRequest createUserRequest) {
    String tenantId = AuthUtil.currentTenantId();
    String username = createUserRequest.getUsername();
    String phone = createUserRequest.getPhone();
    String email = createUserRequest.getEmail();
    String password = createUserRequest.getPassword();
    String userId = IdWorker.getIdStr();

    ArgumentAssert.isFalse(
        StringUtils.isAllBlank(username, phone, email), IamErrorCode.USER_ADD_FAILED);

    PasswordUtil.checkPwd(password);

    User user = UserConvert.INSTANCE.toUser(createUserRequest);
    user.setId(userId);
    user.setTenantId(tenantId);
    user.setUserState(UserStateEnum.DEFAULT.code());
    user.setPassword(PasswordEncoderUtil.encode(password));
    user.setDefaultPwd(true);

    ArgumentAssert.isTrue(
        (this.save(user) && userProfileService.create(userId, createUserRequest)),
        IamErrorCode.USER_ADD_FAILED);
  }

  @Override
  @Transactional(rollbackFor = Exception.class)
  public void update(String userId, UpdateUserRequest updateUserRequest) {
    User oldUser = baseMapper.getById(userId);
    ArgumentAssert.notNull(oldUser, IamErrorCode.USER_DOES_NOT_EXIST);

    User user = UserConvert.INSTANCE.toUser(updateUserRequest);
    user.setId(oldUser.getId());

    ArgumentAssert.isTrue(
        (this.updateById(user) && userProfileService.update(user.getId(), updateUserRequest)),
        IamErrorCode.USER_UPDATE_FAILED);
  }

  @Override
  @Transactional(rollbackFor = Exception.class)
  public void delete(String userId) {
    ArgumentAssert.isFalse(
        Objects.equals(AuthUtil.currentUserId(), userId), IamErrorCode.CANNOT_MODIFIED_YOURSELF);

    User user = baseMapper.getById(userId);

    ArgumentAssert.notNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    ArgumentAssert.isTrue(
        (this.removeById(userId) && userProfileService.delete(userId)),
        IamErrorCode.USER_DELETE_FAILED);
  }

  @Override
  public void disable(String userId, DisableUserRequest disableUserRequest) {
    Boolean disable = disableUserRequest.getDisable();

    UserStateEnum newState = disable ? UserStateEnum.DISABLE : UserStateEnum.DEFAULT;
    UserStateEnum oldState = disable ? UserStateEnum.DEFAULT : UserStateEnum.DISABLE;

    ArgumentAssert.isFalse(
        Objects.equals(AuthUtil.currentUserId(), userId), IamErrorCode.CANNOT_MODIFIED_YOURSELF);

    User user = baseMapper.getById(userId);

    ArgumentAssert.notNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    ArgumentAssert.isTrue(oldState.code() == user.getUserState(), IamErrorCode.USER_STATE_ABNORMAL);

    if (1
        != baseMapper.updateState(
            user.getId(), newState.code(), oldState.code(), AuthUtil.currentUsername())) {
      ArgumentAssert.error(
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
      ArgumentAssert.isFalse(
          Objects.equals(AuthUtil.currentUserId(), userId), IamErrorCode.CANNOT_MODIFIED_YOURSELF);
    } else {
      ArgumentAssert.isFalse(Objects.equals(oldPwd, newPwd), IamErrorCode.PASSWORD_SAME_AS_OLD);
    }

    User user = baseMapper.getById(userId);

    ArgumentAssert.notNull(user, IamErrorCode.USER_DOES_NOT_EXIST);

    // ArgumentAssert.isTrue(
    //     UserTypeEnum.LOCAL.code() == user.getUserType(),
    //     IamErrorCode.NON_LOCAL_USERS_DO_NOT_SUPPORT_CHANGING_PASSWORDS);

    if (!reset) {
      if (PasswordEncoderUtil.matches(newPwd, user.getPassword())) {
        ArgumentAssert.error(IamErrorCode.PASSWORD_SAME_AS_OLD);
      }
      if (!PasswordEncoderUtil.matches(oldPwd, user.getPassword())) {
        ArgumentAssert.error(IamErrorCode.OLD_PASSWORD_IS_WRONG);
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

    ArgumentAssert.isTrue(this.update(updateWrapper), IamErrorCode.PASSWORD_UPDATE_FAILED);
  }

  @Override
  public PageInfo<UserInfoResponse> infoPage(UserInfoPageRequest userInfoPageRequest) {
    Page<UserInfo> pageResult =
        baseMapper.page(PageHelper.startPage(userInfoPageRequest), userInfoPageRequest);

    return PageInfo.of(pageResult.getRecords(), userInfoPageRequest, pageResult.getTotal())
        .map(userInfo -> UserInfoResponse.builder().userInfo(userInfo).build());
  }
}
