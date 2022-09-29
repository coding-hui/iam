package top.wecoding.iam.server.service.impl;

import org.springframework.stereotype.Service;
import top.wecoding.iam.common.model.AuthInfo;
import top.wecoding.iam.common.model.request.LoginRequest;
import top.wecoding.iam.common.model.request.TokenRequest;
import top.wecoding.iam.common.model.response.AuthInfoResponse;
import top.wecoding.iam.common.model.response.CommonLoginResponse;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.server.granter.TokenGranterContext;
import top.wecoding.iam.server.service.AuthService;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Service
public class AuthServiceImpl implements AuthService {

  @Override
  public AuthInfoResponse authInfo(AuthInfo authInfo) {
    return AuthInfoResponse.builder()
        .tenantId(authInfo.getTenantId())
        .userId(authInfo.getUserId())
        .build();
  }

  @Override
  public UserInfoResponse userInfo(AuthInfo authInfo) {
    return null;
  }

  @Override
  public void logout(AuthInfo authInfo) {}

  @Override
  public CommonLoginResponse login(LoginRequest loginRequest) {
    AuthInfo authInfo = TokenGranterContext.grant(TokenRequest.of(loginRequest));
    return CommonLoginResponse.builder()
        .tenantId(authInfo.getTenantId())
        .userId(authInfo.getUserId())
        .username(authInfo.getAccount())
        .build();
  }
}
