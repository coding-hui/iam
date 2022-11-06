package top.wecoding.iam.server.service;

import java.util.List;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.server.entity.Group;
import top.wecoding.iam.server.entity.User;
import top.wecoding.iam.server.entity.UserGroup;
import top.wecoding.mybatis.base.BaseService;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface UserGroupService extends BaseService<UserGroup> {

  UserInfoResponse getUserInfoResponse(String tenantId, User user);

  GroupInfoResponse getGroupInfoResponse(String tenantId, Group group);

  List<UserInfoResponse> getUserInfoResponse(String tenantId, List<User> userList);

  List<GroupInfoResponse> getGroupInfoResponse(String tenantId, List<Group> groupList);
}
