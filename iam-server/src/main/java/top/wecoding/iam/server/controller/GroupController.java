package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import top.wecoding.commons.core.model.PageInfo;
import top.wecoding.commons.core.model.R;
import top.wecoding.iam.common.model.request.CreateGroupRequest;
import top.wecoding.iam.common.model.request.GroupInfoPageRequest;
import top.wecoding.iam.common.model.request.UpdateGroupRequest;
import top.wecoding.iam.common.model.response.CreateGroupResponse;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.server.service.GroupService;
import top.wecoding.iam.server.web.annotation.RequestParameter;
import top.wecoding.web.controller.BaseController;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/groups")
public class GroupController extends BaseController {

  private final GroupService groupService;

  @GetMapping("/{id}")
  public R<GroupInfoResponse> info(@PathVariable("id") String groupId) {
    return R.ok(groupService.getInfo(groupId));
  }

  @GetMapping("")
  public R<PageInfo<GroupInfoResponse>> infoPage(
      @RequestParameter GroupInfoPageRequest groupInfoPageRequest) {
    return R.ok(groupService.infoPage(groupInfoPageRequest));
  }

  @PostMapping("")
  public R<CreateGroupResponse> create(
      @RequestBody @Validated CreateGroupRequest createGroupRequest) {
    return R.ok(groupService.create(createGroupRequest));
  }

  @PutMapping("/{groupId}")
  public R<?> update(
      @PathVariable("groupId") String groupId,
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
