package top.wecoding.iam.api.feign;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestHeader;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.common.model.response.Oauth2ClientInfoResponse;

import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@FeignClient(
    value = "wecoding-iam",
    contextId = "remoteClientDetails",
    url = "${wecoding.feign.iam-server:http://localhost:80}")
public interface RemoteClientDetailsService {

  /**
   * 通过clientId 查询客户端信息
   *
   * @param clientId 用户名
   * @param from 调用标志
   * @return R
   */
  @GetMapping("/api/v1/client/info/{clientId}")
  R<Oauth2ClientInfoResponse> info(
      @PathVariable("clientId") String clientId,
      @RequestHeader(SecurityConstants.FROM) String from);

  /**
   * 查询全部客户端
   *
   * @param from 调用标识
   * @return R
   */
  @GetMapping("/client")
  R<List<Oauth2ClientInfoResponse>> listClientDetails(
      @RequestHeader(SecurityConstants.FROM) String from);
}
