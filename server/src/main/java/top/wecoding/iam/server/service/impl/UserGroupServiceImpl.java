package top.wecoding.iam.server.service.impl;

import java.util.List;
import javax.annotation.Resource;
import org.springframework.stereotype.Service;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.server.mapper.UserGroupMapper;
import top.wecoding.iam.server.pojo.Group;
import top.wecoding.iam.server.pojo.User;
import top.wecoding.iam.server.pojo.UserGroup;
import top.wecoding.iam.server.service.UserGroupService;
import top.wecoding.mybatis.base.BaseServiceImpl;

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
