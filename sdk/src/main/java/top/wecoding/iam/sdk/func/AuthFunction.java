package top.wecoding.iam.sdk.func;

import cn.hutool.core.util.ObjectUtil;
import cn.hutool.core.util.StrUtil;
import org.springframework.util.PatternMatchUtils;
import org.springframework.util.StringUtils;
import top.wecoding.core.exception.user.UnauthorizedException;
import top.wecoding.iam.common.constant.RoleConstant;
import top.wecoding.iam.common.model.LoginUser;
import top.wecoding.iam.sdk.util.AuthUtil;

import java.util.Collection;
import java.util.HashSet;
import java.util.Set;

/**
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
public class AuthFunction {

  /**
   * 需要登录才可访问
   *
   * @return {boolean}
   */
  public boolean requiresLogin() {
    String token = AuthUtil.getToken();
    if (StrUtil.isBlank(token)) {
      throw new UnauthorizedException();
    }
    LoginUser loginUser = AuthUtil.getLoginUser();
    if (ObjectUtil.isNull(loginUser)) {
      throw new UnauthorizedException();
    }

    return true;
  }

  /**
   * 放行所有请求
   *
   * @return {boolean}
   */
  public boolean permitAll() {
    return true;
  }

  /**
   * 只有超管角色才可访问
   *
   * @return {boolean}
   */
  public boolean denyAll() {
    return hasRole(RoleConstant.ADMIN);
  }

  /**
   * 验证用户是否具备某权限
   *
   * @param permit 权限字符串
   * @return {boolean}
   */
  public boolean hasPermission(String permit) {
    return hasPermission(getAllPermissions(), permit);
  }

  /**
   * 验证用户是否具备某权限
   *
   * @param permits 权限字符串
   * @return {boolean}
   */
  public boolean hasAnyPermission(String... permits) {
    Set<String> permissionList = getAllPermissions();
    for (String permission : permits) {
      if (hasPermission(permissionList, permission)) {
        return true;
      }
    }
    throw new UnauthorizedException();
  }

  /**
   * 判断是否有该角色权限
   *
   * @param role 单角色
   * @return {boolean}
   */
  public boolean hasRole(String role) {
    return hasAnyRole(role);
  }

  /**
   * 判断是否有该角色权限
   *
   * @param roles 角色集合
   * @return {boolean}
   */
  public boolean hasAnyRole(String... roles) {
    Set<String> roleKeyList = getAllRoleKeys();
    for (String role : roles) {
      if (hasRole(roleKeyList, role)) {
        return true;
      }
    }
    throw new UnauthorizedException();
  }

  /**
   * 获取当前账号的角色列表
   *
   * @return 角色列表
   */
  public Set<String> getAllRoleKeys() {
    try {
      LoginUser loginUser = AuthUtil.getLoginUser();
      return loginUser.getRoleKeys();
    } catch (Exception e) {
      return new HashSet<>();
    }
  }

  /**
   * 获取当前账号的权限列表
   *
   * @return 权限列表
   */
  public Set<String> getAllPermissions() {
    try {
      LoginUser loginUser = AuthUtil.getLoginUser();
      return loginUser.getPermissions();
    } catch (Exception e) {
      return new HashSet<>();
    }
  }

  /**
   * 判断是否包含权限
   *
   * @param authorities 权限列表
   * @param permission 权限字符串
   * @return 用户是否具备某权限
   */
  public boolean hasPermission(Collection<String> authorities, String permission) {
    return authorities.stream()
        .filter(StringUtils::hasText)
        .anyMatch(
            x ->
                RoleConstant.ALL_PERMISSION.contains(x)
                    || PatternMatchUtils.simpleMatch(x, permission));
  }

  /**
   * 判断是否包含角色
   *
   * @param roles 角色列表
   * @param role 角色
   * @return 用户是否具备某角色权限
   */
  public boolean hasRole(Collection<String> roles, String role) {
    return roles.stream()
        .filter(StringUtils::hasText)
        .anyMatch(x -> RoleConstant.ADMIN.contains(x) || PatternMatchUtils.simpleMatch(x, role));
  }
}
