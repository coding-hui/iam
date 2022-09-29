package top.wecoding.iam.common.provider.client;

import lombok.Data;
import top.wecoding.iam.common.provider.ClientDetails;

/**
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
@Data
public class BaseClientDetails implements ClientDetails {

  /** 客户端id */
  private String clientId;

  /** 客户端密钥 */
  private String clientSecret;

  /** 令牌过期秒数 */
  private Long accessTokenValiditySeconds;

  /** 刷新令牌过期秒数 */
  private Long refreshTokenValiditySeconds;
}
