package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import top.wecoding.commons.core.model.PageInfo;
import top.wecoding.commons.core.model.R;
import top.wecoding.iam.common.model.request.CreateTenantRequest;
import top.wecoding.iam.common.model.request.TenantInfoPageRequest;
import top.wecoding.iam.common.model.request.UpdateTenantRequest;
import top.wecoding.iam.common.model.response.CreateTenantResponse;
import top.wecoding.iam.common.model.response.TenantInfoResponse;
import top.wecoding.iam.common.util.AuthUtil;
import top.wecoding.iam.server.service.TenantService;
import top.wecoding.web.controller.BaseController;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/tenants")
public class TenantController extends BaseController {

  private final TenantService tenantService;

  @GetMapping("/info/{id}")
  public R<TenantInfoResponse> info(@PathVariable("id") String tenantId) {
    return R.ok(tenantService.getInfo(tenantId));
  }

  @GetMapping("/info")
  public R<TenantInfoResponse> info() {
    return R.ok(tenantService.getInfo(AuthUtil.currentTenantId()));
  }

  @PostMapping("/page")
  public R<PageInfo<TenantInfoResponse>> page(
      @RequestBody @Validated TenantInfoPageRequest tenantInfoPageRequest) {
    return R.ok(tenantService.infoPage(tenantInfoPageRequest));
  }

  @PostMapping("")
  public R<CreateTenantResponse> create(
      @RequestBody @Validated CreateTenantRequest createTenantRequest) {
    return R.ok(tenantService.create(createTenantRequest));
  }

  @PutMapping("/{id}")
  public R<Object> update(
      @PathVariable("id") String tenantId,
      @RequestBody @Validated UpdateTenantRequest updateTenantRequest) {
    tenantService.update(tenantId, updateTenantRequest);
    return R.ok();
  }

  @DeleteMapping("/{id}")
  public R<Object> delete(@PathVariable("id") String tenantId) {
    tenantService.delete(tenantId);
    return R.ok();
  }
}
