package top.wecoding.iam.server.convert;

import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;
import top.wecoding.iam.common.model.GroupInfo;
import top.wecoding.iam.common.model.request.CreateGroupRequest;
import top.wecoding.iam.common.model.response.GroupInfoResponse;
import top.wecoding.iam.server.pojo.Group;

/**
 * @author liuyuhui
 * @date 2022/10/6
 * @qq 1515418211
 */
@Mapper
public interface GroupConvert {

  GroupConvert INSTANCE = Mappers.getMapper(GroupConvert.class);

  GroupInfoResponse toGroupInfoResponse(GroupInfo groupInfo);

  GroupInfo toGroupInfo(Group group);

  Group toGroup(CreateGroupRequest createGroupRequest);
}
