package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpHeaders;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import top.wecoding.commons.core.constant.StrPool;
import top.wecoding.commons.core.model.PageInfo;
import top.wecoding.commons.core.model.R;
import top.wecoding.commons.lang.Strings;
import top.wecoding.iam.common.model.request.TokenInfoPageRequest;
import top.wecoding.iam.common.model.response.TokenInfoResponse;
import top.wecoding.iam.framework.InnerAuth;
import top.wecoding.iam.server.service.TokenService;

/**
 * @author liuyuhui
 * @date 2022/10/4
 */
@Slf4j
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/tokens")
public class TokenController {

  private final TokenService tokenService;

  @GetMapping("/page")
  public R<PageInfo<TokenInfoResponse>> page(TokenInfoPageRequest tokenInfoPageRequest) {
    return R.ok(tokenService.infoPage(tokenInfoPageRequest));
  }

  @DeleteMapping("/logout")
  public R<Boolean> logout(
      @RequestHeader(value = HttpHeaders.AUTHORIZATION, required = false) String authHeader) {
    if (Strings.isBlank(authHeader)) {
      return R.ok();
    }
    String tokenValue =
        authHeader.replace(OAuth2AccessToken.TokenType.BEARER.getValue(), StrPool.EMPTY).trim();
    return removeToken(tokenValue);
  }

  @InnerAuth
  @DeleteMapping("/{token}")
  public R<Boolean> removeToken(@PathVariable("token") String token) {
    return R.ok(tokenService.delete(token));
  }
}
