// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package set

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// import mysql driver.
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cmdutil "github.com/coding-hui/iam/internal/iamctl/cmd/util"
	"github.com/coding-hui/iam/internal/iamctl/util/templates"
	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

// DBOptions is an options struct to support 'db' sub command.
type DBOptions struct {
	host     string
	username string
	password string
	Database string

	drop  bool
	admin bool

	genericclioptions.IOStreams
}

var setExample = templates.Examples(`
		# Create new iam platform database and tables
		iamctl set db --mysql.host=127.0.0.1:3306 --mysql.username=iam --mysql.password=iamxxxx --mysql.database=iam

		# Create new iam platform database and tables with a administrator inserted
		iamctl set db --admin --mysql.host=127.0.0.1:3306 --mysql.username=iam --mysql.password=iamxxxx --mysql.database=iam

		# drop and create iam platform database and tables
		iamctl set db -d --mysql.host=127.0.0.1:3306 --mysql.username=iam --mysql.password=iamxxxx --mysql.database=iam`)

// NewDBOptions returns an initialized DBOptions instance.
func NewDBOptions(ioStreams genericclioptions.IOStreams) *DBOptions {
	return &DBOptions{
		host:     "127.0.0.1:3306",
		username: "root",
		password: "root",
		Database: "iam",

		drop:      false,
		admin:     false,
		IOStreams: ioStreams,
	}
}

// NewCmdDB returns new initialized instance of 'db' sub command.
func NewCmdDB(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewDBOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "db",
		DisableFlagsInUseLine: true,
		Short:                 "Initialize the iam database",
		Long:                  "Initialize the iam database.",
		Example:               setExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete())
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
		Aliases:    []string{},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.host, "host", o.host, "MySQL service host address.")
	cmd.Flags().StringVar(&o.username, "username", o.username, "username for access to mysql service.")
	cmd.Flags().StringVar(&o.password, "password", o.password,
		"password for access to mysql, should be used pair with password.")
	cmd.Flags().StringVar(&o.Database, "database", o.Database, "Database name for the server to use.")
	cmd.Flags().BoolVarP(&o.drop, "drop", "d", o.drop, "drop database if exists, pls double check the db name!")
	cmd.Flags().BoolVar(&o.admin, "admin", o.drop, "Insert a administrator user to the database.")

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}

// Complete completes all the required options.
func (o *DBOptions) Complete() error {
	// o.host = viper.GetString("host")
	// o.username = viper.GetString("username")
	// o.password = viper.GetString("password")
	// o.Database = viper.GetString("database")
	// o.drop = viper.GetBool("drop")
	// o.admin = viper.GetBool("admin")

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *DBOptions) Validate() error {
	return nil
}

// Run executes a db sub command using the specified options.
func (o *DBOptions) Run() error {
	if err := o.ensureSchema(); err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		o.username, o.password, o.host, o.Database, true, "Local")

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}

func (o *DBOptions) ensureSchema() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8", o.username, o.password, o.host)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if o.drop {
		dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s", o.Database)
		if err := db.Exec(dropSQL).Error; err != nil {
			return err
		}
		fmt.Fprintf(o.Out, "drop database %s success\n", o.Database)
	}

	createSQL := fmt.Sprintf("CREATE DATABASE if not exists %s CHARSET utf8 COLLATE utf8_general_ci", o.Database)
	if err := db.Exec(createSQL).Error; err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "create database %s success\n", o.Database)

	return nil
}
