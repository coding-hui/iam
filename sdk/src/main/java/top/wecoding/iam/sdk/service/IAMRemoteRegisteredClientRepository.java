package top.wecoding.iam.sdk.service;

import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.ClientAuthenticationMethod;
import org.springframework.security.oauth2.core.OAuth2TokenFormat;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClient;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClientRepository;
import org.springframework.security.oauth2.server.authorization.config.ClientSettings;
import org.springframework.security.oauth2.server.authorization.config.TokenSettings;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.iam.api.feign.RemoteClientDetailsService;

import java.time.Duration;
import java.util.UUID;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@RequiredArgsConstructor
public class IAMRemoteRegisteredClientRepository implements RegisteredClientRepository {

  /** 刷新令牌有效期默认 30 天 */
  private static final int refreshTokenValiditySeconds = 60 * 60 * 24 * 30;

  /** 请求令牌有效期默认 12 小时 */
  private static final int accessTokenValiditySeconds = 60 * 60 * 12;

  private final RemoteClientDetailsService clientDetailsService;

  @Override
  public void save(RegisteredClient registeredClient) {
    throw new UnsupportedOperationException();
  }

  @Override
  public RegisteredClient findById(String id) {
    throw new UnsupportedOperationException();
  }

  @Override
  public RegisteredClient findByClientId(String clientId) {
    // R<OauthClientResponse> result =
    //     clientDetailsService.getClientDetailsById(clientId, SecurityConstants.INNER);
    // if (result.getData() == null) {
    //   return null;
    // }
    // OauthClientResponse clientResponse = result.getData();
    return RegisteredClient.withId(UUID.randomUUID().toString())
        .clientId(clientId)
        .clientSecret(SecurityConstants.NOOP + "wecoding")
        .redirectUri("http:localhost:80")
        .authorizationGrantType(AuthorizationGrantType.PASSWORD)
        .authorizationGrantType(AuthorizationGrantType.REFRESH_TOKEN)
        .authorizationGrantType(AuthorizationGrantType.AUTHORIZATION_CODE)
        .authorizationGrantType(AuthorizationGrantType.CLIENT_CREDENTIALS)
        .clientAuthenticationMethod(ClientAuthenticationMethod.CLIENT_SECRET_BASIC)
        .scope("server")
        .clientSettings(ClientSettings.builder().requireAuthorizationConsent(false).build())
        .tokenSettings(
            TokenSettings.builder()
                .accessTokenFormat(OAuth2TokenFormat.REFERENCE)
                .accessTokenTimeToLive(Duration.ofSeconds(accessTokenValiditySeconds))
                .refreshTokenTimeToLive(Duration.ofSeconds(refreshTokenValiditySeconds))
                .build())
        .build();
  }
}
