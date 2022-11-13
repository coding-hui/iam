package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import top.wecoding.iam.server.entity.UserProfile;

/**
 * @author liuyuhui
 * @date 2022/11/9
 * @qq 1515418211
 */
public interface UserProfileMapper extends BaseMapper<UserProfile> {

  UserProfile getByUserId(String userId);

  UserProfile getByUsername(String username);

  UserProfile getByTenantIdAndUserId(String tenantId, String userId);

  UserProfile getByTenantIdAndUsername(String tenantId, String username);

  UserProfile getByTenantIdAndUsernameAndState(String tenantId, String username, int state);

  int count();
}
