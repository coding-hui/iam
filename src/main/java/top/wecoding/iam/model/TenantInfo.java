package top.wecoding.iam.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Date;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class TenantInfo {

  @JsonProperty("tenant_id")
  private String tenantId;

  @JsonProperty("tenant_name")
  private String tenantName;

  @JsonProperty("owner_id")
  private String ownerId;

  @JsonProperty("username")
  private String username;

  @JsonProperty("description")
  private String annotate;

  @JsonProperty("login_type")
  private Integer loginType;

  @JsonProperty("create_time")
  private Date createTime;

  @JsonProperty("create_ts")
  private Long createTimestamp;
}
