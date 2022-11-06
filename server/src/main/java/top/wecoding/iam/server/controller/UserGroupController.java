package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.result.R;
import top.wecoding.iam.server.service.UserGroupService;
import top.wecoding.web.controller.BaseController;

/**
 * @author liuyuhui
 * @date 2022/10/8
 * @qq 1515418211
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/groups/{groupId}")
public class UserGroupController extends BaseController {

  private final UserGroupService userGroupService;

  @GetMapping("/users")
  public R<PageInfo<Object>> memberList(@PathVariable("groupId") String groupId) {
    return R.ok();
  }

  @PostMapping("/users/add")
  public R<Object> addMember(@PathVariable("groupId") String groupId) {
    return R.ok();
  }

  @PostMapping("/users/remove")
  public R<Object> removeMember(@PathVariable("groupId") String groupId) {
    return R.ok();
  }
}
