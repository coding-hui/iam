package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.Set;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import lombok.experimental.SuperBuilder;
import top.wecoding.commons.core.model.request.PageRequest;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
@Data
@SuperBuilder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class TenantInfoPageRequest extends PageRequest implements CommonRequest {

  @JsonProperty("tenant_ids")
  private Set<String> tenantIds;

  @JsonProperty("tenant_name_fuzzy")
  private String tenantNameFuzzy;
}
