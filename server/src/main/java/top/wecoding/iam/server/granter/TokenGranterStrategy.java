package top.wecoding.iam.server.granter;

import top.wecoding.iam.common.model.AuthInfo;
import top.wecoding.iam.common.model.request.TokenRequest;

/**
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
public interface TokenGranterStrategy {

  /**
   * 根据不同策略进行登录，返回用户信息
   *
   * @param tokenRequest 授权参数
   * @return UserInfo
   */
  AuthInfo grant(TokenRequest tokenRequest);
}
