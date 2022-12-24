package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import top.wecoding.commons.core.model.PageInfo;
import top.wecoding.commons.core.model.R;
import top.wecoding.iam.common.model.request.CreateOauth2ClientPageRequest;
import top.wecoding.iam.common.model.request.CreateOauth2ClientRequest;
import top.wecoding.iam.common.model.request.UpdateOauth2ClientRequest;
import top.wecoding.iam.common.model.response.Oauth2ClientInfoResponse;
import top.wecoding.iam.framework.InnerAuth;
import top.wecoding.iam.server.service.Oauth2ClientService;
import top.wecoding.iam.server.web.annotation.RequestParameter;
import top.wecoding.web.controller.BaseController;

/**
 * @author liuyuhui
 * @date 2022/10/4
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/clients")
public class Oauth2ClientController extends BaseController {

  private final Oauth2ClientService clientService;

  @InnerAuth
  @GetMapping("/info/{clientId}")
  public R<Oauth2ClientInfoResponse> info(@PathVariable("clientId") String clientId) {
    Oauth2ClientInfoResponse info = clientService.getInfoByClientId(clientId);
    return R.ok(info);
  }

  @GetMapping("")
  public R<PageInfo<Oauth2ClientInfoResponse>> page(
      @RequestParameter CreateOauth2ClientPageRequest clientPageRequest) {
    return R.ok(clientService.infoPage(clientPageRequest));
  }

  @PostMapping("")
  public R<?> create(@RequestBody @Validated CreateOauth2ClientRequest createOauth2ClientRequest) {
    clientService.create(createOauth2ClientRequest);
    return R.ok();
  }

  @PutMapping("/id/{id}")
  public R<?> update(
      @PathVariable("id") String id,
      @RequestBody @Validated UpdateOauth2ClientRequest updateOauth2ClientRequest) {
    clientService.update(id, updateOauth2ClientRequest);
    return R.ok();
  }

  @DeleteMapping("/{clientId}")
  public R<?> delete(@PathVariable("clientId") String clientId) {
    clientService.delete(clientId);
    return R.ok();
  }
}
