package top.wecoding.iam.server.security.authorization.authentication;

import java.util.Collections;
import java.util.HashSet;
import java.util.Set;
import lombok.Getter;
import org.springframework.lang.Nullable;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.util.Assert;
import top.wecoding.commons.core.constant.StrPool;
import top.wecoding.iam.common.enums.AuthType;
import top.wecoding.iam.common.model.request.LoginRequest;

/**
 * 自定义授权模式抽象
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
public class OAuth2ResourceOwnerBaseAuthenticationToken extends AbstractAuthenticationToken {

  @Getter private final AuthType authType;

  @Getter private final Authentication clientPrincipal;

  @Getter private final Set<String> scopes;

  @Getter private final LoginRequest loginRequest;

  public OAuth2ResourceOwnerBaseAuthenticationToken(
      AuthType authType,
      Authentication clientPrincipal,
      @Nullable Set<String> scopes,
      LoginRequest loginRequest) {
    super(Collections.emptyList());
    Assert.notNull(authType, "authorizationType cannot be null");
    Assert.notNull(clientPrincipal, "clientPrincipal cannot be null");
    this.authType = authType;
    this.clientPrincipal = clientPrincipal;
    this.scopes =
        Collections.unmodifiableSet(
            scopes != null ? new HashSet<>(scopes) : Collections.emptySet());
    this.loginRequest = loginRequest;
  }

  @Override
  public Object getCredentials() {
    return StrPool.EMPTY;
  }

  @Override
  public Object getPrincipal() {
    return this.clientPrincipal;
  }
}
