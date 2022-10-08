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
@RequestMapping("/api/v1/group/{group_id}")
public class UserGroupController extends BaseController {

  private final UserGroupService userGroupService;

  @GetMapping("/member/list")
  public R<PageInfo<Object>> memberList(@PathVariable("group_id") String groupId) {
    return R.ok();
  }

  @PostMapping("/member/add")
  public R<Object> addMember(@PathVariable("group_id") String groupId) {
    return R.ok();
  }

  @PostMapping("/member/remove")
  public R<Object> removeMember(@PathVariable("group_id") String groupId) {
    return R.ok();
  }
}
