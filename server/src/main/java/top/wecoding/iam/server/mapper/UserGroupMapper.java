package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import java.util.List;
import java.util.Set;
import org.apache.ibatis.annotations.Param;
import top.wecoding.iam.server.pojo.UserGroup;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface UserGroupMapper extends BaseMapper<UserGroup> {
  default int insertBatch0(List<UserGroup> userGroupList, String updatedBy) {
    userGroupList.forEach(userGroup -> userGroup.setCreatedBy(updatedBy));
    return insertBatch(userGroupList);
  }

  int insertBatch(List<UserGroup> userGroupList);

  List<UserGroup> listByTenantIdAndGroupId(
      @Param("tenantId") String tenantId, @Param("groupId") String groupId);

  List<UserGroup> listByTenantIdAndGroupIdAndUserIdList(
      @Param("tenantId") String tenantId,
      @Param("groupId") String groupId,
      @Param("userIds") Set<String> userIds);
}
