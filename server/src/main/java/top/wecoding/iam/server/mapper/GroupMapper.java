package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import top.wecoding.iam.server.pojo.Group;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface GroupMapper extends BaseMapper<Group> {

  Group getByGroupId(String groupId);

  Group getByTenantIdAndGroupName(String tenantId, String groupName);

  Group getByTenantIdAndGroupId(String tenantId, String groupId);

  int count();
}
