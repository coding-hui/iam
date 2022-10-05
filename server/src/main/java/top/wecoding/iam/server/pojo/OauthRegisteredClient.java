package top.wecoding.iam.server.pojo;

import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.*;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.ClientAuthenticationMethod;
import org.springframework.security.oauth2.server.authorization.config.ClientSettings;
import org.springframework.security.oauth2.server.authorization.config.TokenSettings;
import top.wecoding.mybatis.base.BaseEntity;

import java.time.Instant;
import java.util.Set;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_oauth_registered_client")
@EqualsAndHashCode(callSuper = true)
public class OauthRegisteredClient extends BaseEntity {

  @TableId private Long id;

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
