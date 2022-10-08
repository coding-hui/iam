package top.wecoding.iam.api.feign;

import java.util.List;
import java.util.Set;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestParam;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.model.response.UserInfoResponse;

/**
 * @author liuyuhui
 * @date 2022/9/29
 * @qq 1515418211
 */
@FeignClient(name = "iam-user", contextId = "remoteUserService", url = "http://localhost:80")
public interface RemoteUserService {

  /**
   * 通过用户名查询用户、角色信息
   *
   * @param username 用户名
   * @param from 调用标志
   * @return R
   */
  @GetMapping("/api/v1/user/info/{username}")
  R<UserInfoResponse> info(
      @PathVariable("username") String username,
      @RequestHeader(SecurityConstants.FROM) String from);

  /**
   * 通过手机号码查询用户、角色信息
   *
   * @param phone 手机号码
   * @param from 调用标志
   * @return R
   */
  @GetMapping("/app/info/{phone}")
  R<UserInfoResponse> infoByMobile(
      @PathVariable("phone") String phone, @RequestHeader(SecurityConstants.FROM) String from);

  /**
   * 根据部门id，查询对应的用户 id 集合
   *
   * @param deptIds 部门id 集合
   * @param from 调用标志
   * @return 用户 id 集合
   */
  @GetMapping("/user/ids")
  R<List<Long>> listUserIdByDeptIds(
      @RequestParam("deptIds") Set<Long> deptIds,
      @RequestHeader(SecurityConstants.FROM) String from);
}
