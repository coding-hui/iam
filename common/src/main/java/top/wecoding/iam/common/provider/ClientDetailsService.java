package top.wecoding.iam.common.provider;

/**
 * 获取客户端详情接口
 *
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
public interface ClientDetailsService {

  /**
   * 根据 clientId 获取 Client 详情.
   *
   * @param clientId 客户端id
   * @return ClientDetails
   */
  ClientDetails loadClientByClientId(String clientId);
}
