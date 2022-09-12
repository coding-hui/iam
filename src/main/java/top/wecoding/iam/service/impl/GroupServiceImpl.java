package top.wecoding.iam.service.impl;

import org.springframework.stereotype.Service;
import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.mapper.GroupMapper;
import top.wecoding.iam.model.request.CreateGroupRequest;
import top.wecoding.iam.model.request.GroupInfoListRequest;
import top.wecoding.iam.model.request.GroupInfoPageRequest;
import top.wecoding.iam.model.request.UpdateGroupRequest;
import top.wecoding.iam.model.response.CreateGroupResponse;
import top.wecoding.iam.model.response.GroupInfoResponse;
import top.wecoding.iam.pojo.Group;
import top.wecoding.iam.service.GroupService;
import top.wecoding.mybatis.base.BaseServiceImpl;

import javax.annotation.Resource;
import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Service
public class GroupServiceImpl extends BaseServiceImpl<GroupMapper, Group> implements GroupService {

  @Resource private GroupMapper groupMapper;

  @Override
  public GroupInfoResponse getInfo(String groupId) {
    return null;
  }

  @Override
  public CreateGroupResponse create(CreateGroupRequest createGroupRequest) {
    return null;
  }

  @Override
  public void update(String groupId, UpdateGroupRequest updateGroupRequest) {}

  @Override
  public void delete(String groupId) {}

  @Override
  public PageInfo<GroupInfoResponse> infoPage(GroupInfoPageRequest groupInfoPageRequest) {
    return null;
  }

  @Override
  public List<GroupInfoResponse> infoList(GroupInfoListRequest groupInfoListRequest) {
    return null;
  }
}
