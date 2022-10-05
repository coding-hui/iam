package top.wecoding.iam.server.service.impl;

import cn.hutool.core.util.StrUtil;
import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2TokenType;
import org.springframework.security.oauth2.server.authorization.OAuth2Authorization;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.stereotype.Service;
import top.wecoding.iam.common.model.request.LoginRequest;
import top.wecoding.iam.common.model.response.AuthInfoResponse;
import top.wecoding.iam.common.model.response.CommonLoginResponse;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.pojo.AuthInfo;
import top.wecoding.iam.server.service.AuthService;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Service
@RequiredArgsConstructor
public class AuthServiceImpl implements AuthService {

  private final OAuth2AuthorizationService authorizationService;

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
  public boolean logout(String tokenValue) {
    OAuth2Authorization authorization =
        authorizationService.findByToken(tokenValue, OAuth2TokenType.ACCESS_TOKEN);
    if (authorization == null) {
      return true;
    }

    OAuth2Authorization.Token<OAuth2AccessToken> accessToken = authorization.getAccessToken();
    if (accessToken == null || StrUtil.isBlank(accessToken.getToken().getTokenValue())) {
      return true;
    }

    authorizationService.remove(authorization);

    return true;
  }

  @Override
  public CommonLoginResponse login(LoginRequest loginRequest) {
    return CommonLoginResponse.builder().build();
  }
}
