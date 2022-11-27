package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
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
 * @date 2022/10/5
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UpdateOauth2ClientRequest {

  @NotBlank
  @JsonProperty("id")
  private String id;

  @NotBlank
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
