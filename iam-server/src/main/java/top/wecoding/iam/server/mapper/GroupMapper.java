package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import org.apache.ibatis.annotations.Param;
import top.wecoding.iam.common.model.request.GroupInfoListRequest;
import top.wecoding.iam.common.model.request.GroupInfoPageRequest;
import top.wecoding.iam.server.entity.Group;

import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
public interface GroupMapper extends BaseMapper<Group> {

  Group getByGroupId(String groupId);

  Group getByTenantIdAndGroupCode(String tenantId, String groupCode);

  Group getByTenantIdAndGroupId(String tenantId, String groupId);

  List<Group> list(GroupInfoListRequest groupInfoListRequest);

  List<Group> page(@Param("page") IPage<Group> page, @Param("query") GroupInfoPageRequest query);

  int count();
}
