package top.wecoding.iam.server.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import java.io.Serializable;
import org.apache.ibatis.annotations.Param;
import top.wecoding.iam.common.model.request.CreateOauth2ClientPageRequest;
import top.wecoding.iam.server.entity.Oauth2Client;

/**
 * @author liuyuhui
 * @date 2022/10/5
 * @qq 1515418211
 */
public interface Oauth2ClientMapper extends BaseMapper<Oauth2Client> {

  Oauth2Client getById(Serializable clientId);

  Oauth2Client getByClientId(String clientId);

  Page<Oauth2Client> page(
      @Param("page") Page<Oauth2Client> page, @Param("query") CreateOauth2ClientPageRequest query);
}
