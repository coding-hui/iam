package top.wecoding.iam.server.pojo;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import com.baomidou.mybatisplus.extension.handlers.JacksonTypeHandler;
import lombok.*;
import top.wecoding.iam.common.pojo.OAuth2ClientSettings;
import top.wecoding.iam.common.pojo.OAuth2TokenSettings;
import top.wecoding.mybatis.base.BaseEntity;

import java.time.Instant;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_oauth2_registered_client")
@EqualsAndHashCode(callSuper = true)
public class Oauth2Client extends BaseEntity {

  @TableId private String id;

  private String clientId;

  private Instant clientIdIssuedAt;

  private String clientSecret;

  private Instant clientSecretExpiresAt;

  private String clientName;

  private String clientAuthenticationMethods;

  private String authorizationGrantTypes;

  private String redirectUris;

  private String scopes;

  @TableField(typeHandler = JacksonTypeHandler.class)
  private OAuth2ClientSettings clientSettings;

  @TableField(typeHandler = JacksonTypeHandler.class)
  private OAuth2TokenSettings tokenSettings;
}
