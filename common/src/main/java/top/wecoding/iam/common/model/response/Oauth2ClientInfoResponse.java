package top.wecoding.iam.common.model.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.Instant;
import java.util.Set;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import top.wecoding.iam.common.entity.OAuth2ClientSettings;
import top.wecoding.iam.common.entity.OAuth2TokenSettings;

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

  @JsonProperty("id")
  private String id;

  @JsonProperty("client_id")
  private String clientId;

  @JsonProperty("client_id_issued_at")
  private Instant clientIdIssuedAt;

  @JsonProperty("client_secret")
  private String clientSecret;

  @JsonProperty("client_secret_expires_at")
  private Instant clientSecretExpiresAt;

  @JsonProperty("client_name")
  private String clientName;

  @JsonProperty("client_authentication_methods")
  private Set<String> clientAuthenticationMethods;

  @JsonProperty("authorization_grant_types")
  private Set<String> authorizationGrantTypes;

  @JsonProperty("redirect_uris")
  private Set<String> redirectUris;

  @JsonProperty("scopes")
  private Set<String> scopes;

  @JsonProperty("client_settings")
  private OAuth2ClientSettings clientSettings;

  @JsonProperty("token_settings")
  private OAuth2TokenSettings tokenSettings;
}
