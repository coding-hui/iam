package top.wecoding.iam.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.*;
import top.wecoding.core.model.request.PageRequest;

import java.util.Set;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class TenantInfoPageRequest extends PageRequest {

  @JsonProperty("tenant_ids")
  private Set<String> tenantIds;

  @JsonProperty("tenant_name_fuzzy")
  private String tenantNameFuzzy;

  @JsonProperty("reverse")
  private Boolean reverse;
}
