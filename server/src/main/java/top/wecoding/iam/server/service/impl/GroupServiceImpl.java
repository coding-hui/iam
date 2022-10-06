package top.wecoding.iam.server.service.impl;

import com.baomidou.mybatisplus.core.toolkit.IdWorker;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.iam.common.enums.IamErrorCode;
import top.wecoding.iam.common.model.GroupInfo;
import top.wecoding.iam.common.model.request.CreateGroupRequest;
import top.wecoding.iam.common.model.request.GroupInfoListRequest;
import top.wecoding.iam.common.model.request.GroupInfoPageRequest;
import top.wecoding.iam.common.model.request.UpdateGroupRequest;
import top.wecoding.iam.common.model.response.CreateGroupResponse;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.sdk.utili.AuthUtil;
import top.wecoding.iam.server.convert.GroupConvert;
import top.wecoding.iam.server.mapper.GroupMapper;
import top.wecoding.iam.server.mapper.UserGroupMapper;
import top.wecoding.iam.server.pojo.Group;
import top.wecoding.iam.server.pojo.UserGroup;
import top.wecoding.iam.server.service.GroupService;
import top.wecoding.mybatis.base.BaseServiceImpl;

import java.io.Serializable;
import java.util.*;
import java.util.stream.Collectors;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Service
@RequiredArgsConstructor
public class GroupServiceImpl extends BaseServiceImpl<GroupMapper, Group> implements GroupService {

  private final GroupMapper groupMapper;

  private final UserGroupMapper userGroupMapper;

  @Override
  public GroupInfoResponse getInfo(String groupId) {
    Group group = groupMapper.getByGroupId(groupId);

    AssertUtils.isNotNull(group, IamErrorCode.GROUP_DOES_NOT_EXIST);

    GroupInfo groupInfo = GroupConvert.INSTANCE.toGroupInfo(group);

    return GroupInfoResponse.builder().groupInfo(groupInfo).build();
  }

  @Override
  @Transactional(rollbackFor = Exception.class)
  public CreateGroupResponse create(CreateGroupRequest createGroupRequest) {
    String tenantId = AuthUtil.currentTenantId();
    String groupName = createGroupRequest.getGroupName();
    String groupId = IdWorker.getIdStr();

    Set<String> userIds = createGroupRequest.getUserIds();

    AssertUtils.isNull(
        groupMapper.getByTenantIdAndGroupName(tenantId, groupName),
        IamErrorCode.GROUP_ALREADY_EXIST);

    Group group = GroupConvert.INSTANCE.toGroup(createGroupRequest);

    List<UserGroup> userGroupList =
        Optional.ofNullable(userIds).orElse(Collections.emptySet()).stream()
            .map(
                userId ->
                    UserGroup.builder().tenantId(tenantId).groupId(groupId).userId(userId).build())
            .collect(Collectors.toList());

    if (1 != groupMapper.insert(group)) {
      AssertUtils.error(IamErrorCode.GROUP_ADD_FAILED);
    }
    if (0 != userGroupList.size()) {
      userGroupMapper.insertBatch0(userGroupList, AuthUtil.currentUsername());
    }

    return CreateGroupResponse.builder().groupId(groupId).groupName(groupName).build();
  }

  @Override
  public void update(String groupId, UpdateGroupRequest updateGroupRequest) {
    String tenantId = AuthUtil.currentTenantId();
    Set<String> inputUserIdSet = updateGroupRequest.getInputIdSet();
    Set<String> outputUserIdSet = updateGroupRequest.getOutputIdSet();
    Set<String> ignoreSet =
        inputUserIdSet.size() < outputUserIdSet.size()
            ? inputUserIdSet.stream().filter(outputUserIdSet::contains).collect(Collectors.toSet())
            : outputUserIdSet.stream().filter(inputUserIdSet::contains).collect(Collectors.toSet());

    inputUserIdSet.removeAll(ignoreSet);
    outputUserIdSet.removeAll(ignoreSet);

    if (0 == inputUserIdSet.size() && 0 == outputUserIdSet.size()) {
      return;
    }

    Set<String> mergeSet = new HashSet<>();

    mergeSet.addAll(inputUserIdSet);
    mergeSet.addAll(outputUserIdSet);

    List<UserGroup> userGroupList =
        userGroupMapper.listByTenantIdAndGroupIdAndUserIdList(tenantId, groupId, mergeSet);

    Map<String, Serializable> userGroupMap =
        userGroupList.stream().collect(Collectors.toMap(UserGroup::getUserId, UserGroup::getId));

    List<UserGroup> inputUserGroupList =
        inputUserIdSet.stream()
            .filter(userId -> !userGroupMap.containsKey(userId))
            .map(
                userId ->
                    UserGroup.builder().tenantId(tenantId).groupId(groupId).userId(userId).build())
            .collect(Collectors.toList());

    Set<Serializable> outputIdSet =
        outputUserIdSet.stream()
            .filter(userGroupMap::containsKey)
            .map(userGroupMap::get)
            .collect(Collectors.toSet());

    if (0 == inputUserGroupList.size() && 0 == outputIdSet.size()) {
      return;
    }

    if (0 != outputIdSet.size()) {
      userGroupMapper.deleteBatchIds(outputIdSet);
    }
    if (0 != inputUserGroupList.size()) {
      userGroupMapper.insertBatch0(inputUserGroupList, AuthUtil.currentUsername());
    }
  }

  @Override
  public void delete(String groupId) {
    String tenantId = AuthUtil.currentTenantId();

    Group group = groupMapper.getByTenantIdAndGroupId(tenantId, groupId);

    AssertUtils.isNotNull(group, IamErrorCode.GROUP_DOES_NOT_EXIST);

    List<UserGroup> userGroupList = userGroupMapper.listByTenantIdAndGroupId(tenantId, groupId);
    List<Serializable> userGroupIds =
        Optional.ofNullable(userGroupList).orElse(Collections.emptyList()).stream()
            .map(UserGroup::getId)
            .collect(Collectors.toList());

    AssertUtils.isTrue(
        1 == groupMapper.deleteById(group.getId()), IamErrorCode.GROUP_DELETE_FAILED);
    userGroupMapper.deleteBatchIds(userGroupIds);
  }

  @Override
  public PageInfo<GroupInfoResponse> infoPage(GroupInfoPageRequest groupInfoPageRequest) {
    return null;
  }

  @Override
  public List<GroupInfoResponse> infoList(GroupInfoListRequest groupInfoListRequest) {
    return null;
  }
}
