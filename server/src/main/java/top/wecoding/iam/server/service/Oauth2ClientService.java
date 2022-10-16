package top.wecoding.iam.server.service;

import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.common.model.request.CreateOauth2ClientPageRequest;
import top.wecoding.iam.common.model.request.CreateOauth2ClientRequest;
import top.wecoding.iam.common.model.request.UpdateOauth2ClientRequest;
import top.wecoding.iam.common.model.response.Oauth2ClientInfoResponse;

/**
 * @author liuyuhui
 * @date 2022/10/5
 * @qq 1515418211
 */
public interface Oauth2ClientService {

  Oauth2ClientInfoResponse getInfo(String clientId);

  void create(CreateOauth2ClientRequest createOauth2ClientRequest);

  void update(UpdateOauth2ClientRequest updateOauth2ClientRequest);

  void delete(String clientId);

  PageInfo<Oauth2ClientInfoResponse> infoPage(CreateOauth2ClientPageRequest clientPageRequest);
}
