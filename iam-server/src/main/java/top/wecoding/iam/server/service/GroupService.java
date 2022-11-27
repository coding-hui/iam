package top.wecoding.iam.server.service;

import top.wecoding.commons.core.model.PageInfo;
import top.wecoding.iam.common.model.request.CreateGroupRequest;
import top.wecoding.iam.common.model.request.GroupInfoListRequest;
import top.wecoding.iam.common.model.request.GroupInfoPageRequest;
import top.wecoding.iam.common.model.request.UpdateGroupRequest;
import top.wecoding.iam.common.model.response.CreateGroupResponse;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.server.entity.Group;
import top.wecoding.mybatis.base.BaseService;

import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
public interface GroupService extends BaseService<Group> {

  GroupInfoResponse getInfo(String groupId);

  CreateGroupResponse create(CreateGroupRequest createGroupRequest);

  void update(String groupId, UpdateGroupRequest updateGroupRequest);

  void delete(String groupId);

  PageInfo<GroupInfoResponse> infoPage(GroupInfoPageRequest groupInfoPageRequest);

  List<GroupInfoResponse> infoList(GroupInfoListRequest groupInfoListRequest);
}
