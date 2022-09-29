package top.wecoding.iam.common.provider;

import java.io.Serializable;

/**
 * Oauth2 客户端详情
 *
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
public interface ClientDetails extends Serializable {

  /**
   * 客户端id.
   *
   * @return The client id.
   */
  String getClientId();

  /**
   * 客户端密钥.
   *
   * @return The client secret.
   */
  String getClientSecret();

  /**
   * 客户端token过期时间.
   *
   * @return the access token validity period
   */
  Long getAccessTokenValiditySeconds();

  /**
   * 客户端刷新token过期时间.
   *
   * @return the refresh token validity period
   */
  Long getRefreshTokenValiditySeconds();
}
