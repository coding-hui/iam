package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import org.apache.ibatis.annotations.Param;
import top.wecoding.iam.common.model.request.GroupInfoListRequest;
import top.wecoding.iam.common.model.request.GroupInfoPageRequest;
import top.wecoding.iam.server.pojo.Group;

import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface GroupMapper extends BaseMapper<Group> {

  Group getByGroupId(String groupId);

  Group getByTenantIdAndGroupName(String tenantId, String groupName);

  Group getByTenantIdAndGroupId(String tenantId, String groupId);

  List<Group> list(GroupInfoListRequest groupInfoListRequest);

  Page<Group> page(@Param("page") Page<Group> page, @Param("query") GroupInfoPageRequest query);

  int count();
}
