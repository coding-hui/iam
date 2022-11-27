package top.wecoding.iam.server.convert;

import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;
import top.wecoding.iam.common.model.GroupInfo;
import top.wecoding.iam.common.model.request.CreateGroupRequest;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.server.entity.Group;

/**
 * @author liuyuhui
 * @date 2022/10/6
 */
@Mapper
public interface GroupConvert {

  GroupConvert INSTANCE = Mappers.getMapper(GroupConvert.class);

  GroupInfoResponse toGroupInfoResponse(GroupInfo groupInfo);

  GroupInfoResponse toGroupInfoResponse(Group group);

  GroupInfo toGroupInfo(Group group);

  Group toGroup(CreateGroupRequest createGroupRequest);
}
