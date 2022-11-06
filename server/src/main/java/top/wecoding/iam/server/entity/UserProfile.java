package top.wecoding.iam.server.entity;

import com.baomidou.mybatisplus.annotation.TableName;
import lombok.*;
import top.wecoding.mybatis.base.BaseEntity;

import java.util.Date;

/**
 * @author liuyuhui
 * @date 2022/11/6
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_user_profile")
@EqualsAndHashCode(callSuper = true)
public class UserProfile extends BaseEntity {

  private String userId;

  private String nickName;

  private String avatar;

  private Date birthday;

  private String gender;

  private String country;

  private String company;

  private String address;

  private String province;

  private String city;

  private String streetAddress;

  private String postalCode;
}
