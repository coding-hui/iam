package top.wecoding.iam.server.entity;

import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.util.Date;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import top.wecoding.mybatis.base.BaseEntity;

/**
 * @author liuyuhui
 * @date 2022/11/6
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_user_profile")
@EqualsAndHashCode(callSuper = true)
public class UserProfile extends BaseEntity {

  /** userid */
  @TableId private String userId;

  private String tenantId;

  /** 姓名 */
  private String name;

  /** 用户账号 */
  private String username;

  /** 名 */
  private String firstName;

  /** 姓 */
  private String lastName;

  /** 中间名 */
  private String middleName;

  /** 昵称 */
  private String nickName;

  /** 用户个人网站 */
  private String website;

  /** 用户头像 */
  private String avatar;

  /** 邮箱 */
  private String email;

  /** 手机号 */
  private String phone;

  /** 生日 */
  private Date birthday;

  /** 性别 */
  private String gender;

  /** 国家 */
  private String country;

  /** 公司 */
  private String company;

  /** 地址 */
  private String address;

  /** 省份 */
  private String province;

  /** 城市 */
  private String city;

  /** 街道 */
  private String streetAddress;

  /** 邮编 */
  private String postalCode;

  /** 第三方系统Id */
  private String externalId;
}
