package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.model.request.CreateOauth2ClientPageRequest;
import top.wecoding.iam.common.model.request.CreateOauth2ClientRequest;
import top.wecoding.iam.common.model.request.UpdateOauth2ClientRequest;
import top.wecoding.iam.common.model.response.Oauth2ClientInfoResponse;
import top.wecoding.iam.sdk.InnerAuth;
import top.wecoding.iam.server.service.Oauth2ClientService;
import top.wecoding.web.controller.BaseController;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/client")
public class Oauth2ClientController extends BaseController {

  private final Oauth2ClientService clientService;

  @InnerAuth
  @GetMapping("/info/{clientId}")
  public R<Oauth2ClientInfoResponse> info(@PathVariable("clientId") String clientId) {
    Oauth2ClientInfoResponse info = clientService.getInfo(clientId);
    return R.ok(info);
  }

  @PostMapping("/page")
  public R<PageInfo<Oauth2ClientInfoResponse>> page(
      CreateOauth2ClientPageRequest clientPageRequest) {
    return R.ok(clientService.infoPage(clientPageRequest));
  }

  @PostMapping("")
  public R<?> create(@RequestBody @Validated CreateOauth2ClientRequest createOauth2ClientRequest) {
    clientService.create(createOauth2ClientRequest);
    return R.ok();
  }

  @PutMapping("")
  public R<?> update(@RequestBody @Validated UpdateOauth2ClientRequest updateOauth2ClientRequest) {
    clientService.update(updateOauth2ClientRequest);
    return R.ok();
  }

  @DeleteMapping("/{clientId}")
  public R<?> delete(@PathVariable("clientId") String clientId) {
    clientService.delete(clientId);
    return R.ok();
  }
}
