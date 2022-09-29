package top.wecoding.iam.server.service;

import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.server.pojo.Group;
import top.wecoding.iam.server.pojo.User;
import top.wecoding.iam.server.pojo.UserGroup;
import top.wecoding.mybatis.base.BaseService;

import java.util.List;

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