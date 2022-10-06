package top.wecoding.iam.server.convert;

import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;
import top.wecoding.iam.common.model.GroupInfo;
import top.wecoding.iam.common.model.request.CreateGroupRequest;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.server.pojo.Group;
import top.wecoding.iam.server.pojo.UserGroup;

/**
 * @author liuyuhui
 * @date 2022/10/6
 * @qq 1515418211
 */
@Mapper
public interface UserGroupConvert {

  UserGroupConvert INSTANCE = Mappers.getMapper(UserGroupConvert.class);

  GroupInfoResponse toGroupInfoResponse(GroupInfo groupInfo);

  GroupInfo toGroupInfo(Group group);

  UserGroup toUserGroup(CreateGroupRequest createGroupRequest);
}
