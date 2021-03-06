package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSIAMGroup_normal(t *testing.T) {
	var conf iam.GetGroupOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSGroupExists("aws_iam_group.group", &conf),
					testAccCheckAWSGroupAttributes(&conf),
				),
			},
		},
	})
}

func testAccCheckAWSGroupDestroy(s *terraform.State) error {
	iamconn := testAccProvider.Meta().(*AWSClient).iamconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_iam_group" {
			continue
		}

		// Try to get group
		_, err := iamconn.GetGroup(&iam.GetGroupInput{
			GroupName: aws.String(rs.Primary.ID),
		})
		if err == nil {
			return fmt.Errorf("still exist.")
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ec2err.Code() != "NoSuchEntity" {
			return err
		}
	}

	return nil
}

func testAccCheckAWSGroupExists(n string, res *iam.GetGroupOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Group name is set")
		}

		iamconn := testAccProvider.Meta().(*AWSClient).iamconn

		resp, err := iamconn.GetGroup(&iam.GetGroupInput{
			GroupName: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		*res = *resp

		return nil
	}
}

func testAccCheckAWSGroupAttributes(group *iam.GetGroupOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *group.Group.GroupName != "test-group" {
			return fmt.Errorf("Bad name: %s", *group.Group.GroupName)
		}

		if *group.Group.Path != "/" {
			return fmt.Errorf("Bad path: %s", *group.Group.Path)
		}

		return nil
	}
}

const testAccAWSGroupConfig = `
resource "aws_iam_group" "group" {
	name = "test-group"
	path = "/"
}
`
