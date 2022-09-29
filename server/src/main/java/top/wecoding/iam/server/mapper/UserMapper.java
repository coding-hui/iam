package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import top.wecoding.iam.server.pojo.User;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface UserMapper extends BaseMapper<User> {

  User getInfoByUsernameAntTenantName(String username, String tenantName);
}
