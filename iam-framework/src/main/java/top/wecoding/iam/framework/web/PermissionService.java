package top.wecoding.iam.framework.web;

import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.util.PatternMatchUtils;
import org.springframework.util.StringUtils;
import top.wecoding.commons.lang.Objects;

import java.util.Collection;

/**
 * @author liuyuhui
 * @date 2022/10/1
 */
public class PermissionService {

  /**
   * 判断接口是否有任意xxx，xxx权限
   *
   * @param permissions 权限
   * @return {boolean}
   */
  public boolean hasPermission(String... permissions) {
    if (Objects.isEmpty(permissions)) {
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
