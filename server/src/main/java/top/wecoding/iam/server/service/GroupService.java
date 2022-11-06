package top.wecoding.iam.server.service;

import java.util.List;
import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.common.model.request.CreateGroupRequest;
import top.wecoding.iam.common.model.request.GroupInfoListRequest;
import top.wecoding.iam.common.model.request.GroupInfoPageRequest;
import top.wecoding.iam.common.model.request.UpdateGroupRequest;
import top.wecoding.iam.common.model.response.CreateGroupResponse;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.server.entity.Group;
import top.wecoding.mybatis.base.BaseService;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface GroupService extends BaseService<Group> {

  GroupInfoResponse getInfo(String groupId);

  CreateGroupResponse create(CreateGroupRequest createGroupRequest);

  void update(String groupId, UpdateGroupRequest updateGroupRequest);

  void delete(String groupId);

  PageInfo<GroupInfoResponse> infoPage(GroupInfoPageRequest groupInfoPageRequest);

  List<GroupInfoResponse> infoList(GroupInfoListRequest groupInfoListRequest);
}
