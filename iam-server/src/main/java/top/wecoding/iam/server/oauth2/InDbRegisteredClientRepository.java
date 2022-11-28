package top.wecoding.iam.server.oauth2;

import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.ClientAuthenticationMethod;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClient;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClientRepository;
import org.springframework.security.oauth2.server.authorization.settings.ClientSettings;
import org.springframework.security.oauth2.server.authorization.settings.TokenSettings;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;
import top.wecoding.iam.common.model.response.Oauth2ClientInfoResponse;
import top.wecoding.iam.server.service.Oauth2ClientService;

import java.util.Set;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@Service
@RequiredArgsConstructor
public class InDbRegisteredClientRepository implements RegisteredClientRepository {

  private final Oauth2ClientService clientDetailsService;

  @Override
  public void save(RegisteredClient registeredClient) {
    throw new UnsupportedOperationException();
  }

  @Override
  public RegisteredClient findById(String id) {
    Assert.hasText(id, "client id cannot be empty");
    Oauth2ClientInfoResponse clientResponse = clientDetailsService.getInfoById(id);
    return buildRegisteredClient(clientResponse);
  }

  @Override
  public RegisteredClient findByClientId(String clientId) {
    Assert.hasText(clientId, "clientId cannot be empty");
    Oauth2ClientInfoResponse clientResponse = clientDetailsService.getInfoByClientId(clientId);
    return buildRegisteredClient(clientResponse);
  }

  private RegisteredClient buildRegisteredClient(Oauth2ClientInfoResponse clientResponse) {
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
            .clientSecret(clientResponse.getClientSecret())
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
    // Custom authorization grant type
    return new AuthorizationGrantType(authorizationGrantType);
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
    // Custom client authentication method
    return new ClientAuthenticationMethod(clientAuthenticationMethod);
  }
}
