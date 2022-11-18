package top.wecoding.iam.server.controller;

import cn.hutool.core.util.StrUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpHeaders;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.web.bind.annotation.*;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.model.request.TokenInfoPageRequest;
import top.wecoding.iam.common.model.response.TokenInfoResponse;
import top.wecoding.iam.framework.InnerAuth;
import top.wecoding.iam.server.service.TokenService;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
@Slf4j
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/token")
public class TokenController {

  private final TokenService tokenService;

  @GetMapping("/page")
  public R<PageInfo<TokenInfoResponse>> page(TokenInfoPageRequest tokenInfoPageRequest) {
    return R.ok(tokenService.infoPage(tokenInfoPageRequest));
  }

  @DeleteMapping("/logout")
  public R<Boolean> logout(
      @RequestHeader(value = HttpHeaders.AUTHORIZATION, required = false) String authHeader) {
    if (StrUtil.isBlank(authHeader)) {
      return R.ok();
    }
    String tokenValue =
        authHeader.replace(OAuth2AccessToken.TokenType.BEARER.getValue(), StrUtil.EMPTY).trim();
    return removeToken(tokenValue);
  }

  @InnerAuth
  @DeleteMapping("/{token}")
  public R<Boolean> removeToken(@PathVariable("token") String token) {
    return R.ok(tokenService.delete(token));
  }
}
