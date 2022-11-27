package top.wecoding.iam.framework.authentication.password;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.OAuth2ErrorCodes;
import org.springframework.security.oauth2.core.OAuth2Token;
import org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClient;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenGenerator;
import top.wecoding.iam.framework.authentication.OAuth2ResourceOwnerBaseAuthenticationProvider;

import java.util.Map;

/**
 * 密码认证
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
public class OAuth2ResourceOwnerPasswordAuthenticationProvider
    extends OAuth2ResourceOwnerBaseAuthenticationProvider<
        OAuth2ResourceOwnerPasswordAuthenticationToken> {

  private static final Logger LOGGER =
      LogManager.getLogger(OAuth2ResourceOwnerPasswordAuthenticationProvider.class);

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
  public UsernamePasswordAuthenticationToken buildToken(Map<String, Object> reqParameters) {
    String username = (String) reqParameters.get(OAuth2ParameterNames.USERNAME);
    String password = (String) reqParameters.get(OAuth2ParameterNames.PASSWORD);
    return new UsernamePasswordAuthenticationToken(username, password);
  }

  @Override
  public boolean supports(Class<?> authentication) {
    boolean supports =
        OAuth2ResourceOwnerPasswordAuthenticationToken.class.isAssignableFrom(authentication);
    LOGGER.debug("supports authentication=" + authentication + " returning " + supports);
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
