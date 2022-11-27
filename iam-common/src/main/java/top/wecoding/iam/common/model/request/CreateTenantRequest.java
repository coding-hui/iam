package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class CreateTenantRequest {

  @NotBlank
  @JsonProperty("tenant_name")
  private String tenantName;

  @NotBlank
  @Size(max = 500)
  @JsonProperty("description")
  private String description;
}
