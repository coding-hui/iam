package top.wecoding.iam.service.impl;

import org.springframework.stereotype.Service;
import top.wecoding.iam.mapper.UserGroupMapper;
import top.wecoding.iam.model.response.GroupInfoResponse;
import top.wecoding.iam.model.response.UserInfoResponse;
import top.wecoding.iam.pojo.Group;
import top.wecoding.iam.pojo.User;
import top.wecoding.iam.pojo.UserGroup;
import top.wecoding.iam.service.UserGroupService;
import top.wecoding.mybatis.base.BaseServiceImpl;

import javax.annotation.Resource;
import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Service
public class UserGroupServiceImpl extends BaseServiceImpl<UserGroupMapper, UserGroup>
    implements UserGroupService {

  @Resource private UserGroupMapper userGroupMapper;

  @Override
  public UserInfoResponse getUserInfoResponse(String tenantId, User user) {
    return null;
  }

  @Override
  public GroupInfoResponse getGroupInfoResponse(String tenantId, Group group) {
    return null;
  }

  @Override
  public List<UserInfoResponse> getUserInfoResponse(String tenantId, List<User> userList) {
    return null;
  }

  @Override
  public List<GroupInfoResponse> getGroupInfoResponse(String tenantId, List<Group> groupList) {
    return null;
  }
}
