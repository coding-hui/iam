package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.model.request.CreateGroupRequest;
import top.wecoding.iam.common.model.request.GroupInfoListRequest;
import top.wecoding.iam.common.model.request.GroupInfoPageRequest;
import top.wecoding.iam.common.model.request.UpdateGroupRequest;
import top.wecoding.iam.common.model.response.CreateGroupResponse;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.server.service.GroupService;
import top.wecoding.web.controller.BaseController;

import java.util.List;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/group")
public class GroupController extends BaseController {

  private final GroupService groupService;

  @GetMapping("/{id}")
  public R<GroupInfoResponse> info(@PathVariable("id") String groupId) {
    return R.ok(groupService.getInfo(groupId));
  }

  @PostMapping("/list")
  public R<List<GroupInfoResponse>> infoList(
      @RequestBody @Validated GroupInfoListRequest groupInfoListRequest) {
    return R.ok(groupService.infoList(groupInfoListRequest));
  }

  @PostMapping("/page")
  public R<PageInfo<GroupInfoResponse>> infoPage(
      @RequestBody @Validated GroupInfoPageRequest groupInfoPageRequest) {
    return R.ok(groupService.infoPage(groupInfoPageRequest));
  }

  @PostMapping("")
  public R<CreateGroupResponse> create(
      @RequestBody @Validated CreateGroupRequest createGroupRequest) {
    return R.ok(groupService.create(createGroupRequest));
  }

  @PutMapping("/{id}")
  public R<?> update(
      @PathVariable("id") String groupId,
      @RequestBody @Validated UpdateGroupRequest updateGroupRequest) {
    groupService.update(groupId, updateGroupRequest);
    return R.ok();
  }

  @DeleteMapping("/{id}")
  public R<?> delete(@PathVariable("id") String groupId) {
    groupService.delete(groupId);
    return R.ok();
  }
}
