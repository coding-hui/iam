package top.wecoding.iam.server.service;

import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.common.model.request.TokenInfoPageRequest;
import top.wecoding.iam.common.model.response.TokenInfoResponse;

/**
 * @author liuyuhui
 * @date 2022/10/6
 * @qq 1515418211
 */
public interface TokenService {

  boolean delete(String tokenValue);

  PageInfo<TokenInfoResponse> infoPage(TokenInfoPageRequest tokenInfoPageRequest);
}
