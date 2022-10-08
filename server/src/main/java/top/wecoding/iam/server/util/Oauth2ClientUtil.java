package top.wecoding.iam.server.util;

import cn.hutool.core.lang.UUID;
import java.time.Instant;
import java.util.Optional;
import java.util.Set;
import lombok.experimental.UtilityClass;
import org.springframework.util.StringUtils;
import top.wecoding.iam.common.model.request.CreateOauth2ClientRequest;
import top.wecoding.iam.common.model.request.UpdateOauth2ClientRequest;
import top.wecoding.iam.common.model.response.Oauth2ClientInfoResponse;
import top.wecoding.iam.server.pojo.Oauth2Client;

/**
 * @author liuyuhui
 * @date 2022/10/5
 * @qq 1515418211
 */
@UtilityClass
public class Oauth2ClientUtil {

  public Oauth2ClientInfoResponse toOauth2ClientInfoResponse(Oauth2Client client) {
    Set<String> clientAuthenticationMethods =
        StringUtils.commaDelimitedListToSet(client.getClientAuthenticationMethods());
    Set<String> authorizationGrantTypes =
        StringUtils.commaDelimitedListToSet(client.getAuthorizationGrantTypes());
    Set<String> redirectUris = StringUtils.commaDelimitedListToSet(client.getRedirectUris());
    Set<String> clientScopes = StringUtils.commaDelimitedListToSet(client.getScopes());

    return Oauth2ClientInfoResponse.builder()
        .id(client.getId())
        .clientId(client.getClientId())
        .clientIdIssuedAt(client.getClientIdIssuedAt())
        .clientSecret(client.getClientSecret())
        .clientSecretExpiresAt(client.getClientSecretExpiresAt())
        .clientName(client.getClientName())
        .clientAuthenticationMethods(clientAuthenticationMethods)
        .authorizationGrantTypes(authorizationGrantTypes)
        .redirectUris(redirectUris)
        .scopes(clientScopes)
        .clientSettings(client.getClientSettings())
        .tokenSettings(client.getTokenSettings())
        .build();
  }

  public Oauth2Client ofOauth2Client(CreateOauth2ClientRequest createOauth2ClientRequest) {
    String clientId =
        StringUtils.hasText(createOauth2ClientRequest.getClientId())
            ? createOauth2ClientRequest.getClientId()
            : UUID.randomUUID().toString();

    Instant clientIdIssuedAt =
        Optional.ofNullable(createOauth2ClientRequest.getClientIdIssuedAt())
            .orElseGet(Instant::now);

    Instant clientSecretExpiresAt =
        Optional.ofNullable(createOauth2ClientRequest.getClientSecretExpiresAt())
            .orElseGet(Instant::now);

    Set<String> clientAuthenticationMethods =
        createOauth2ClientRequest.getClientAuthenticationMethods();

    Set<String> authorizationGrantTypes = createOauth2ClientRequest.getAuthorizationGrantTypes();

    return Oauth2Client.builder()
        .clientId(clientId)
        .clientIdIssuedAt(clientIdIssuedAt)
        .clientSecret(createOauth2ClientRequest.getClientSecret())
        .clientSecretExpiresAt(clientSecretExpiresAt)
        .clientName(createOauth2ClientRequest.getClientName())
        .clientAuthenticationMethods(
            StringUtils.collectionToCommaDelimitedString(clientAuthenticationMethods))
        .authorizationGrantTypes(
            StringUtils.collectionToCommaDelimitedString(authorizationGrantTypes))
        .redirectUris(
            StringUtils.collectionToCommaDelimitedString(
                createOauth2ClientRequest.getRedirectUris()))
        .scopes(StringUtils.collectionToCommaDelimitedString(createOauth2ClientRequest.getScopes()))
        .clientSettings(createOauth2ClientRequest.getClientSettings())
        .tokenSettings(createOauth2ClientRequest.getTokenSettings())
        .build();
  }

  public Oauth2Client ofOauth2Client(UpdateOauth2ClientRequest updateOauth2ClientRequest) {
    String id = updateOauth2ClientRequest.getId();

    String clientId = updateOauth2ClientRequest.getClientId();

    Instant clientIdIssuedAt =
        Optional.ofNullable(updateOauth2ClientRequest.getClientIdIssuedAt())
            .orElseGet(Instant::now);

    Instant clientSecretExpiresAt =
        Optional.ofNullable(updateOauth2ClientRequest.getClientSecretExpiresAt())
            .orElseGet(Instant::now);

    Set<String> clientAuthenticationMethods =
        updateOauth2ClientRequest.getClientAuthenticationMethods();

    Set<String> authorizationGrantTypes = updateOauth2ClientRequest.getAuthorizationGrantTypes();

    return Oauth2Client.builder()
        .id(id)
        .clientId(clientId)
        .clientIdIssuedAt(clientIdIssuedAt)
        .clientSecret(updateOauth2ClientRequest.getClientSecret())
        .clientSecretExpiresAt(clientSecretExpiresAt)
        .clientName(updateOauth2ClientRequest.getClientName())
        .clientAuthenticationMethods(
            StringUtils.collectionToCommaDelimitedString(clientAuthenticationMethods))
        .authorizationGrantTypes(
            StringUtils.collectionToCommaDelimitedString(authorizationGrantTypes))
        .redirectUris(
            StringUtils.collectionToCommaDelimitedString(
                updateOauth2ClientRequest.getRedirectUris()))
        .scopes(StringUtils.collectionToCommaDelimitedString(updateOauth2ClientRequest.getScopes()))
        .clientSettings(updateOauth2ClientRequest.getClientSettings())
        .tokenSettings(updateOauth2ClientRequest.getTokenSettings())
        .build();
  }
}
