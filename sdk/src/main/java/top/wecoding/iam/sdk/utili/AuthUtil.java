package top.wecoding.iam.sdk.utili;

import lombok.experimental.UtilityClass;
import org.springframework.security.core.context.SecurityContext;
import org.springframework.security.core.context.SecurityContextHolder;
import top.wecoding.core.exception.user.UnauthorizedException;
import top.wecoding.iam.common.userdetails.LoginUser;

import java.util.Optional;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
@UtilityClass
public class AuthUtil {

  public String currentUserId() {
    return currentLoginUser().getUserId();
  }

  public String currentTenantId() {
    return currentLoginUser().getTenantId();
  }

  public LoginUser currentLoginUser() {
    return Optional.ofNullable(SecurityContextHolder.getContext())
        .map(SecurityContext::getAuthentication)
        .map(authentication -> (LoginUser) authentication.getPrincipal())
        .orElseThrow(UnauthorizedException::new);
  }
}
