package top.wecoding.iam.server.security.authorization.authentication.password;

import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.OAuth2ErrorCodes;
import org.springframework.security.oauth2.core.OAuth2Token;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClient;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenGenerator;
import top.wecoding.iam.common.model.request.LoginRequest;
import top.wecoding.iam.server.security.authorization.authentication.OAuth2ResourceOwnerBaseAuthenticationProvider;

/**
 * 密码认证
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
public class OAuth2ResourceOwnerPasswordAuthenticationProvider
    extends OAuth2ResourceOwnerBaseAuthenticationProvider<
        OAuth2ResourceOwnerPasswordAuthenticationToken> {

  /**
   * Constructs an {@code OAuth2AuthorizationCodeAuthenticationProvider} using the provided
   * parameters.
   *
   * @param authorizationService the authorization service
   * @param tokenGenerator the token generator
   * @since 0.2.3
   */
  public OAuth2ResourceOwnerPasswordAuthenticationProvider(
      OAuth2AuthorizationService authorizationService,
      OAuth2TokenGenerator<? extends OAuth2Token> tokenGenerator) {
    super(null, authorizationService, tokenGenerator);
  }

  /**
   * Constructs an {@code OAuth2AuthorizationCodeAuthenticationProvider} using the provided
   * parameters.
   *
   * @param authenticationManager the authentication manager
   * @param authorizationService the authorization service
   * @param tokenGenerator the token generator
   * @since 0.2.3
   */
  public OAuth2ResourceOwnerPasswordAuthenticationProvider(
      AuthenticationManager authenticationManager,
      OAuth2AuthorizationService authorizationService,
      OAuth2TokenGenerator<? extends OAuth2Token> tokenGenerator) {
    super(authenticationManager, authorizationService, tokenGenerator);
  }

  @Override
  public UsernamePasswordAuthenticationToken buildToken(LoginRequest loginRequest) {
    String username = loginRequest.getPasswordPayload().getAccount();
    String password = loginRequest.getPasswordPayload().getPassword();
    return new UsernamePasswordAuthenticationToken(username, password);
  }

  @Override
  public boolean supports(Class<?> authentication) {
    boolean supports =
        OAuth2ResourceOwnerPasswordAuthenticationToken.class.isAssignableFrom(authentication);
    if (log.isDebugEnabled()) {
      log.debug("supports authentication=" + authentication + " returning " + supports);
    }
    return supports;
  }

  @Override
  public void checkClient(RegisteredClient registeredClient) {
    assert registeredClient != null;
    if (!registeredClient.getAuthorizationGrantTypes().contains(AuthorizationGrantType.PASSWORD)) {
      throw new OAuth2AuthenticationException(OAuth2ErrorCodes.UNAUTHORIZED_CLIENT);
    }
  }
}
