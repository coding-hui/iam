package top.wecoding.iam.server.entity;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.*;
import top.wecoding.mybatis.base.BaseEntity;

import java.util.Date;

/**
 * 系统用户表
 *
 * @author liuyuhui
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_user")
@EqualsAndHashCode(callSuper = true)
public class User extends BaseEntity {

  @TableId private Long id;

  private String userId;

  private String tenantId;

  /** 姓名 */
  private String name;

  /** 用户名（登录名） */
  private String username;

  @TableField("pwd")
  private String password;

  private String nickName;

  private String avatar;

  private Date birthday;

  private String gender;

  private String email;

  private String phone;

  private Integer userType;

  private Integer userState;

  @TableField("def_pwd")
  private Boolean defaultPwd;

  private String country;

  private String company;

  private String address;

  private String province;

  private String city;

  private String streetAddress;

  private String postalCode;

  private String externalId;

  private String lastLoginIp;

  private Date lastLoginTime;

  private Integer loginCount;
}
