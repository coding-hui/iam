package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import top.wecoding.commons.core.model.PageInfo;
import top.wecoding.commons.core.model.R;
import top.wecoding.iam.common.model.request.*;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.util.AuthUtil;
import top.wecoding.iam.framework.InnerAuth;
import top.wecoding.iam.server.service.UserService;
import top.wecoding.iam.server.web.annotation.RequestParameter;

/**
 * @author liuyuhui
 * @date 2022/10/4
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/users")
public class UserController {

  private final UserService userService;

  @GetMapping("/info")
  public R<UserInfoResponse> info() {
    return R.ok(userService.getInfo(AuthUtil.currentLoginUser()));
  }

  @InnerAuth
  @GetMapping("/{username}/info")
  public R<UserInfoResponse> info(@PathVariable("username") String username) {
    return R.ok(userService.getInfoByUsername(username));
  }

  @GetMapping("")
  public R<PageInfo<UserInfoResponse>> page(
      @RequestParameter UserInfoPageRequest userInfoPageRequest) {
    return R.ok(userService.infoPage(userInfoPageRequest));
  }

  @PostMapping("")
  public R<?> create(@RequestBody @Validated CreateUserRequest createUserRequest) {
    userService.create(createUserRequest);
    return R.ok();
  }

  @PutMapping("/{userId}")
  public R<?> update(
      @PathVariable("userId") String userId,
      @RequestBody @Validated UpdateUserRequest updateUserRequest) {
    userService.update(userId, updateUserRequest);
    return R.ok();
  }

  @DeleteMapping("/{userId}")
  public R<?> delete(@PathVariable("userId") String userId) {
    userService.delete(userId);
    return R.ok();
  }

  @PutMapping("/{id}/disable")
  public R<?> disable(
      @PathVariable("id") String userId,
      @RequestBody @Validated DisableUserRequest disableUserRequest) {
    userService.disable(userId, disableUserRequest);
    return R.ok();
  }

  @PutMapping("/{userId}/password")
  public R<?> password(
      @PathVariable("userId") String userId,
      @RequestBody @Validated UpdatePasswordRequest updatePasswordRequest) {
    userService.password(userId, updatePasswordRequest);
    return R.ok();
  }

  @PutMapping("/password")
  public R<?> password(@RequestBody @Validated UpdatePasswordRequest updatePasswordRequest) {
    userService.password(AuthUtil.currentUserId(), updatePasswordRequest);
    return R.ok();
  }
}