package top.wecoding.iam.sdk.web;

import cn.hutool.core.util.ArrayUtil;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.util.PatternMatchUtils;
import org.springframework.util.StringUtils;

import java.util.Collection;

/**
 * @author liuyuhui
 * @date 2022/10/1
 * @qq 1515418211
 */
public class PermissionService {

  /**
   * 判断接口是否有任意xxx，xxx权限
   *
   * @param permissions 权限
   * @return {boolean}
   */
  public boolean hasPermission(String... permissions) {
    if (ArrayUtil.isEmpty(permissions)) {
      return false;
    }
    Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
    if (authentication == null) {
      return false;
    }
    Collection<? extends GrantedAuthority> authorities = authentication.getAuthorities();
    return authorities.stream()
        .map(GrantedAuthority::getAuthority)
        .filter(StringUtils::hasText)
        .anyMatch(x -> PatternMatchUtils.simpleMatch(permissions, x));
  }
}
