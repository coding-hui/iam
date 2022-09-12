package top.wecoding.iam.pojo;

import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.*;
import top.wecoding.mybatis.base.LogicDeletedBaseEntity;

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
public class User extends LogicDeletedBaseEntity {

  @TableId private Long id;

  private String userId;

  private String tenantId;

  private String username;

  private String password;

  private String aliasName;

  private String avatar;

  private Date birthday;

  private String gender;

  private String email;

  private String phone;

  private String lastLoginIp;

  private Date lastLoginTime;

  private String userType;

  private String userState;

  private Boolean defaultPwd;

  private String infos;
}
