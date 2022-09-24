package top.wecoding.iam.controller;

import io.swagger.annotations.ApiOperation;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import top.wecoding.auth.context.AuthContextHolder;
import top.wecoding.core.constant.TokenConstant;
import top.wecoding.core.result.R;
import top.wecoding.iam.model.request.LoginRequest;
import top.wecoding.iam.model.response.AuthInfoResponse;
import top.wecoding.iam.model.response.CommonLoginResponse;
import top.wecoding.iam.model.response.UserInfoResponse;
import top.wecoding.iam.service.AuthService;
import top.wecoding.iam.service.ValidateService;
import top.wecoding.web.controller.BaseController;

import javax.annotation.Resource;

/**
 * 认证服务
 *
 * @author liuyuhui
 * @qq 1515418211
 */
@RestController
@RequestMapping("/auth")
public class AuthController extends BaseController {

  @Resource private AuthService authService;
  @Resource private ValidateService validateService;

  @GetMapping("/")
  public R<AuthInfoResponse> authInfo() {
    return R.ok(authService.authInfo(AuthContextHolder.getContext()));
  }

  @GetMapping("/user-info")
  public R<UserInfoResponse> userInfo() {
    return R.ok(authService.userInfo(AuthContextHolder.getContext()));
  }

  @PostMapping("token")
  @ApiOperation(value = "获取认证Token", notes = "账号:account,密码:password")
  public R<CommonLoginResponse> token(@Validated @RequestBody LoginRequest loginRequest) {
    return R.ok(authService.login(loginRequest));
  }

  /** 验证码 */
  @GetMapping("/code")
  public R<?> code() {
    return validateService.createCode();
  }

  /** 手机验证码 */
  @GetMapping("/sms-code/{mobile}")
  public R<?> smsCode(@PathVariable("mobile") String mobile) {
    return validateService.createSmsCode(mobile);
  }

  /** 登出 */
  @DeleteMapping("/logout")
  public R<?> logout(@RequestHeader(value = TokenConstant.AUTHENTICATION) String authHeader) {
    authService.logout(AuthContextHolder.getContext());
    return R.ok();
  }
}
