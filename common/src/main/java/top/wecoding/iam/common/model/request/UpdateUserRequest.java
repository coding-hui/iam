package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.fasterxml.jackson.databind.annotation.JsonNaming;
import java.util.Date;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@JsonNaming(PropertyNamingStrategies.SnakeCaseStrategy.class)
public class UpdateUserRequest {

  private String userId;

  private String name;

  private String username;

  private String nickName;

  private Date birthday;

  private String gender;

  private String email;

  private String phone;

  private String country;

  private String company;

  private String address;

  private String province;

  private String city;

  private String streetAddress;

  private String postalCode;

  private String externalId;
}
