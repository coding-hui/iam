package top.wecoding.iam.common.model.response;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import top.wecoding.iam.common.pojo.OAuth2ClientSettings;
import top.wecoding.iam.common.pojo.OAuth2TokenSettings;

import java.time.Instant;
import java.util.Set;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Oauth2ClientInfoResponse {

  private String id;

  private String clientId;

  private Instant clientIdIssuedAt;

  private String clientSecret;

  private Instant clientSecretExpiresAt;

  private String clientName;

  private Set<String> clientAuthenticationMethods;

  private Set<String> authorizationGrantTypes;

  private Set<String> redirectUris;

  private Set<String> scopes;

  private OAuth2ClientSettings clientSettings;

  private OAuth2TokenSettings tokenSettings;
}
