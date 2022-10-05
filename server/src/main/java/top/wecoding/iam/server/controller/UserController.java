package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.sdk.InnerAuth;
import top.wecoding.iam.server.service.UserService;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/user")
public class UserController {

  private final UserService userService;

  @InnerAuth
  @GetMapping("/info/{username}")
  public R<UserInfoResponse> info(@PathVariable("username") String username) {
    return R.ok(userService.getInfoByUsername(username));
  }
}
