package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import java.io.Serializable;
import org.apache.ibatis.annotations.Param;
import top.wecoding.iam.common.model.request.TenantInfoPageRequest;
import top.wecoding.iam.server.pojo.Tenant;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface TenantMapper extends BaseMapper<Tenant> {

  Tenant getByTenantName(String tenantName);

  Tenant getByUsername(String username);

  Tenant getByTenantId(String tenantId);

  Page<Tenant> page(@Param("page") Page<Tenant> page, @Param("query") TenantInfoPageRequest query);

  int count();

  int updateTenantName(Serializable id, String newTenantName, String oldTenantName);

  int updateTenantAnnotate(Serializable id, String annotate, String updatedBy);
}
