package top.wecoding.iam.sdk.util;

import cn.hutool.core.util.ObjectUtil;
import cn.hutool.core.util.StrUtil;
import lombok.experimental.UtilityClass;
import lombok.extern.slf4j.Slf4j;
import top.wecoding.core.constant.StrPool;
import top.wecoding.core.constant.TokenConstant;
import top.wecoding.core.enums.UserTypeEnum;
import top.wecoding.core.util.WebUtil;
import top.wecoding.iam.sdk.context.AuthContextHolder;
import top.wecoding.iam.common.model.LoginUser;

import javax.servlet.http.HttpServletRequest;

/**
 * 安全服务工具类
 *
 * @author liuyuhui
 * @date 2022/5/13
 * @qq 1515418211
 */
@Slf4j
@UtilityClass
public class AuthUtil {

  /**
   * 获取客户端 ID
   *
   * @return clientId
   */
  public String getClientId() {
    return AuthContextHolder.getContext().getClientId();
  }

  /**
   * 获取登录用户账户
   *
   * @return 登录用户账户
   */
  public String getAccount() {
    return AuthContextHolder.getContext().getAccount();
  }

  /**
   * 获取登录用户 ID
   *
   * @return UserID
   */
  public String getUserId() {
    return AuthContextHolder.getContext().getUserId();
  }

  /**
   * 获取登录用户，获取不到抛出异常
   *
   * @return LoginUser
   */
  public LoginUser getLoginUser() {
    return new LoginUser();
  }

  /**
   * 判断当前用户是否是管理员
   *
   * @return 结果
   */
  public boolean isAdmin() {
    return isUserType(UserTypeEnum.SUPER_ADMIN);
  }

  /**
   * 判断当前登录用户是否是指定用户类型
   *
   * @param userTypeEnum 用户类型
   * @return 结果
   */
  public boolean isUserType(UserTypeEnum userTypeEnum) {
    String userType = getLoginUser().getUserType();
    return userTypeEnum.eq(userType);
  }

  public String getToken() {
    HttpServletRequest request = WebUtil.getRequest();
    if (ObjectUtil.isNull(request)) {}
    return getToken(request);
  }

  /**
   * 获取 Token
   *
   * @param request request
   * @return Token
   */
  public String getToken(HttpServletRequest request) {
    String token = request.getHeader(TokenConstant.AUTHENTICATION);
    if (StrUtil.isBlank(token)) {
      token = request.getParameter(TokenConstant.AUTHENTICATION);
    }
    return replaceTokenPrefix(token);
  }

  /** 去掉 Token 前缀 */
  public String replaceTokenPrefix(String token) {
    // 如果前端设置了令牌前缀，则裁剪掉前缀
    if (StrUtil.isNotEmpty(token) && token.startsWith(TokenConstant.PREFIX)) {
      token = token.replaceFirst(TokenConstant.PREFIX, StrPool.EMPTY);
    }
    return token;
  }
}
