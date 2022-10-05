package top.wecoding.iam.server.controller;

import cn.hutool.core.collection.CollUtil;
import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.server.authorization.config.ClientSettings;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.model.response.OauthClientResponse;
import top.wecoding.iam.sdk.InnerAuth;
import top.wecoding.web.controller.BaseController;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/client")
public class OauthClientController extends BaseController {

  @InnerAuth
  @GetMapping(value = {"/{clientId}"})
  public R<OauthClientResponse> info(@PathVariable("clientId") String clientId) {
    OauthClientResponse client =
        OauthClientResponse.builder()
            .clientId(clientId)
            .clientSecret(SecurityConstants.NOOP + "wecoding")
            .scopes(CollUtil.newHashSet("server"))
            .clientSettings(ClientSettings.builder().requireAuthorizationConsent(false).build())
            .build();
    return R.ok(client);
  }
}
