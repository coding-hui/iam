package top.wecoding.iam.sdk.service;

import java.util.Optional;
import java.util.Set;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.ClientAuthenticationMethod;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClient;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClientRepository;
import org.springframework.security.oauth2.server.authorization.config.ClientSettings;
import org.springframework.security.oauth2.server.authorization.config.TokenSettings;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.iam.api.feign.RemoteClientDetailsService;
import top.wecoding.iam.common.model.response.Oauth2ClientInfoResponse;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@RequiredArgsConstructor
public class IAMRemoteRegisteredClientRepository implements RegisteredClientRepository {

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
  @SneakyThrows(Exception.class)
  public RegisteredClient findByClientId(String clientId) {
    Oauth2ClientInfoResponse clientResponse =
        Optional.ofNullable(clientDetailsService.info(clientId, SecurityConstants.INNER).getData())
            .orElseThrow(
                () ->
                    new OAuth2AuthenticationException(
                        "Get query client failed, check database chain"));

    Set<String> clientAuthenticationMethods = clientResponse.getClientAuthenticationMethods();
    Set<String> authorizationGrantTypes = clientResponse.getAuthorizationGrantTypes();
    Set<String> redirectUris = clientResponse.getRedirectUris();
    Set<String> clientScopes = clientResponse.getScopes();
    ClientSettings clientSettings = clientResponse.getClientSettings().toClientSettings();
    TokenSettings tokenSettings = clientResponse.getTokenSettings().toTokenSettings();

    RegisteredClient.Builder builder =
        RegisteredClient.withId(clientResponse.getId())
            .clientId(clientResponse.getClientId())
            .clientIdIssuedAt(clientResponse.getClientIdIssuedAt())
            .clientSecret(SecurityConstants.NOOP + clientResponse.getClientSecret())
            .clientSecretExpiresAt(clientResponse.getClientSecretExpiresAt())
            .clientName(clientResponse.getClientName())
            .clientAuthenticationMethods(
                (authenticationMethods) ->
                    clientAuthenticationMethods.forEach(
                        authenticationMethod ->
                            authenticationMethods.add(
                                resolveClientAuthenticationMethod(authenticationMethod))))
            .authorizationGrantTypes(
                (grantTypes) ->
                    authorizationGrantTypes.forEach(
                        grantType -> grantTypes.add(resolveAuthorizationGrantType(grantType))))
            .redirectUris((uris) -> uris.addAll(redirectUris))
            .scopes((scopes) -> scopes.addAll(clientScopes))
            .clientSettings(clientSettings)
            .tokenSettings(tokenSettings);

    return builder.build();
  }

  private AuthorizationGrantType resolveAuthorizationGrantType(String authorizationGrantType) {
    if (AuthorizationGrantType.AUTHORIZATION_CODE.getValue().equals(authorizationGrantType)) {
      return AuthorizationGrantType.AUTHORIZATION_CODE;
    } else if (AuthorizationGrantType.CLIENT_CREDENTIALS
        .getValue()
        .equals(authorizationGrantType)) {
      return AuthorizationGrantType.CLIENT_CREDENTIALS;
    } else if (AuthorizationGrantType.REFRESH_TOKEN.getValue().equals(authorizationGrantType)) {
      return AuthorizationGrantType.REFRESH_TOKEN;
    }
    return new AuthorizationGrantType(authorizationGrantType); // Custom authorization grant type
  }

  private ClientAuthenticationMethod resolveClientAuthenticationMethod(
      String clientAuthenticationMethod) {
    if (ClientAuthenticationMethod.CLIENT_SECRET_BASIC
        .getValue()
        .equals(clientAuthenticationMethod)) {
      return ClientAuthenticationMethod.CLIENT_SECRET_BASIC;
    } else if (ClientAuthenticationMethod.CLIENT_SECRET_POST
        .getValue()
        .equals(clientAuthenticationMethod)) {
      return ClientAuthenticationMethod.CLIENT_SECRET_POST;
    } else if (ClientAuthenticationMethod.NONE.getValue().equals(clientAuthenticationMethod)) {
      return ClientAuthenticationMethod.NONE;
    }
    return new ClientAuthenticationMethod(
        clientAuthenticationMethod); // Custom client authentication method
  }
}
