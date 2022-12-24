package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import java.io.Serializable;
import org.apache.ibatis.annotations.Param;
import top.wecoding.iam.common.entity.UserInfo;
import top.wecoding.iam.common.model.request.UserInfoPageRequest;
import top.wecoding.iam.server.entity.Oauth2Client;
import top.wecoding.iam.server.entity.User;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
public interface UserMapper extends BaseMapper<User> {

  User getById(Serializable id);

  UserInfo getInfoById(Serializable id);

  UserInfo getInfoByUsername(String username);

  Page<UserInfo> page(
      @Param("page") IPage<Oauth2Client> page, @Param("query") UserInfoPageRequest query);

  int flushLastLoginInfo(Serializable id, String lastLoginIp);

  int updateState(Serializable id, int newState, int oldState, String updatedBy);

  int count();
}