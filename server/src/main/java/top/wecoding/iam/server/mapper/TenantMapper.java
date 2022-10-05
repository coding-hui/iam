package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
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

  int count();
}
