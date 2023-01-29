package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import top.wecoding.commons.core.model.R;
import top.wecoding.iam.common.model.request.UpdateUserRequest;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.server.service.UserProfileService;
import top.wecoding.iam.server.service.UserService;

/**
 * @author liuyuhui
 * @date 2022/10/4
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/users/{userId:\\d+}/profiles")
public class UserProfileController {

  private final UserService userService;

  private final UserProfileService userProfileService;

  @GetMapping("")
  public R<UserInfoResponse> details(@PathVariable("userId") String userId) {
    return R.ok(userService.getInfoById(userId));
  }

  @PutMapping("")
  public R<?> update(
      @PathVariable("userId") String userId,
      @RequestBody @Validated UpdateUserRequest updateUserRequest) {
    userProfileService.update(userId, updateUserRequest);
    return R.ok();
  }
}
