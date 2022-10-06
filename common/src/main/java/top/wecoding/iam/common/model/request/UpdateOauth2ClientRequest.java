package top.wecoding.iam.common.model.request;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import top.wecoding.iam.common.pojo.OAuth2ClientSettings;
import top.wecoding.iam.common.pojo.OAuth2TokenSettings;

import javax.validation.constraints.NotBlank;
import java.time.Instant;
import java.util.Set;

/**
 * @author liuyuhui
 * @date 2022/10/5
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UpdateOauth2ClientRequest {

  @NotBlank private String id;

  @NotBlank private String clientId;

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
