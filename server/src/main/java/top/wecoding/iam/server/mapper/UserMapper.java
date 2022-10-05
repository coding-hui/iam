package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import top.wecoding.iam.server.pojo.User;

import java.io.Serializable;
import java.util.Collection;
import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface UserMapper extends BaseMapper<User> {

  User getById(Serializable id);

  User getByUserId(String userId);

  User getByUsername(String username);

  User getByTenantIdAndUserId(String tenantId, String userId);

  User getByTenantIdAndUsername(String tenantId, String username);

  User getByTenantIdAndUsernameAndState(String tenantId, String username, int state);

  List<User> listByTenantId(String tenantId);

  List<User> listByTenantIdAndUserIds(String tenantId, Collection<String> userIds);

  int flushLastLoginInfo(Serializable id, String lastLoginIp);

  int count();
}
