package top.wecoding.iam.server.controller;

import cn.hutool.core.util.StrUtil;
import io.swagger.annotations.ApiOperation;
import org.springframework.http.HttpHeaders;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.model.request.LoginRequest;
import top.wecoding.iam.common.model.response.CommonLoginResponse;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.userdetails.LoginUser;
import top.wecoding.iam.sdk.utili.AuthUtil;
import top.wecoding.iam.server.service.AuthService;
import top.wecoding.iam.server.service.ValidateService;
import top.wecoding.web.controller.BaseController;

import javax.annotation.Resource;

/**
 * 认证服务
 *
 * @author liuyuhui
 * @qq 1515418211
 */
@RestController
@RequestMapping("/api/v1/auth")
public class AuthController extends BaseController {

  @Resource private AuthService authService;

  @Resource private ValidateService validateService;

  @GetMapping("")
  public R<LoginUser> authInfo() {
    return R.ok(AuthUtil.currentLoginUser());
  }

  @GetMapping("/user-info")
  public R<UserInfoResponse> userInfo() {
    return R.ok();
  }

  @PostMapping("token")
  @ApiOperation(value = "获取认证Token", notes = "账号:account,密码:password")
  public R<CommonLoginResponse> token(@Validated @RequestBody LoginRequest loginRequest) {
    return R.ok(authService.login(loginRequest));
  }

  @GetMapping("/code")
  public R<?> code() {
    return validateService.createCode();
  }

  @GetMapping("/sms-code/{mobile}")
  public R<?> smsCode(@PathVariable("mobile") String mobile) {
    return validateService.createSmsCode(mobile);
  }

  @DeleteMapping("/logout")
  public R<Boolean> logout(
      @RequestHeader(value = HttpHeaders.AUTHORIZATION, required = false) String authHeader) {
    if (StrUtil.isBlank(authHeader)) {
      return R.ok();
    }
    String tokenValue =
        authHeader.replace(OAuth2AccessToken.TokenType.BEARER.getValue(), StrUtil.EMPTY).trim();
    return R.ok(authService.logout(tokenValue));
  }
}
