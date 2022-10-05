package top.wecoding.iam.common.model.response;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OauthClientDetailsResponse {

  /** 客户端ID */
  private String clientId;

  /** 客户端密钥 */
  private String clientSecret;

  /** 资源ID */
  private String resourceIds;

  /** 作用域 */
  private String scope;

  /** 授权方式（A,B,C） */
  private String authorizedGrantTypes;

  /** 回调地址 */
  private String webServerRedirectUri;

  /** 权限 */
  private String authorities;

  /** 请求令牌有效时间 */
  private Integer accessTokenValidity;

  /** 刷新令牌有效时间 */
  private Integer refreshTokenValidity;

  /** 扩展信息 */
  private String additionalInformation;

  /** 是否自动放行 */
  private String autoapprove;
}
