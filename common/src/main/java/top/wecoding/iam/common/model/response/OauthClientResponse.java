package top.wecoding.iam.common.model.response;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.ClientAuthenticationMethod;
import org.springframework.security.oauth2.server.authorization.config.ClientSettings;
import org.springframework.security.oauth2.server.authorization.config.TokenSettings;

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
public class OauthClientResponse {

  private String id;

  private String clientId;

  private Instant clientIdIssuedAt;

  private String clientSecret;

  private Instant clientSecretExpiresAt;

  private String clientName;

  private Set<ClientAuthenticationMethod> clientAuthenticationMethods;

  private Set<AuthorizationGrantType> authorizationGrantTypes;

  private Set<String> redirectUris;

  private Set<String> scopes;

  private ClientSettings clientSettings;

  private TokenSettings tokenSettings;
}
