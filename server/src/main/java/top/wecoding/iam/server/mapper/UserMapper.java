package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import java.io.Serializable;
import java.util.Collection;
import java.util.List;
import org.apache.ibatis.annotations.Param;
import top.wecoding.iam.common.model.request.UserInfoListRequest;
import top.wecoding.iam.common.model.request.UserInfoPageRequest;
import top.wecoding.iam.server.entity.Oauth2Client;
import top.wecoding.iam.server.entity.User;

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

  List<User> list(UserInfoListRequest query);

  List<User> listByTenantId(String tenantId);

  List<User> listByTenantIdAndUserIds(String tenantId, Collection<String> userIds);

  Page<User> page(
      @Param("page") Page<Oauth2Client> page, @Param("query") UserInfoPageRequest query);

  int flushLastLoginInfo(Serializable id, String lastLoginIp);

  int updateState(Serializable id, int newState, int oldState, String updatedBy);

  int count();
}
