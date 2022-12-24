package top.wecoding.iam.common.util;

import java.util.Optional;
import lombok.experimental.UtilityClass;
import org.springframework.security.core.context.SecurityContext;
import org.springframework.security.core.context.SecurityContextHolder;
import top.wecoding.commons.core.exception.user.UnauthorizedException;
import top.wecoding.iam.common.userdetails.LoginUser;

/**
 * @author liuyuhui
 * @date 2022/10/4
 */
@UtilityClass
public class AuthUtil {

  public String currentUserId() {
    return currentLoginUser().userInfo().getUserId();
  }

  public String currentUsername() {
    return currentLoginUser().userInfo().getUsername();
  }

  public String currentTenantId() {
    return currentLoginUser().userInfo().getTenantId();
  }

  public String currentTenantName() {
    return currentLoginUser().userInfo().getTenantName();
  }

  public LoginUser currentLoginUser() {
    return Optional.ofNullable(SecurityContextHolder.getContext())
        .map(SecurityContext::getAuthentication)
        .map(authentication -> (LoginUser) authentication.getPrincipal())
        .orElseThrow(UnauthorizedException::new);
  }
}
