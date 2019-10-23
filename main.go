package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/workspaces"
	"os"
)

func getListIAMUser(sess *session.Session) ([]string, error) {
	svc := iam.New(sess)
	input := &iam.ListUsersInput{}
	pageNum := 0
	var userList []string
	err := svc.ListUsersPages(input,
		func(page *iam.ListUsersOutput, lastPage bool) bool {
			pageNum++
			for _, key := range page.Users{
				userList = append(userList, *key.UserName)
			}
			return pageNum <= 100
		})
	return userList, err
}

func envRegister(info []string) (*credentials.Value, error) {
	os.Setenv("AWS_ACCESS_KEY_ID", info[1])
	os.Setenv("AWS_SECRET_ACCESS_KEY", info[2])

	creds := credentials.NewEnvCredentials()
	// Retrieve the credentials value
	credValues, err := creds.Get()

	return &credValues, err
}

func createSession(region string) (*session.Session, error) {
	// Create a Session with a custom region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	return sess, err
}

func getListWorkspaceUser(sess *session.Session) ([]string, error) {
	svc := workspaces.New(sess)
	input := &workspaces.DescribeWorkspacesInput{}
	pageNum := 0
	var userList []string
	err := svc.DescribeWorkspacesPages(input,
		func(page *workspaces.DescribeWorkspacesOutput, lastPage bool) bool {
			pageNum++
			for _, user := range page.Workspaces {
				userList = append(userList, *user.UserName)
			}
			return pageNum <= 100
		})
	return userList, err
}

func wrtieCSV(iam, workspace []string) {
	// 파일 생성
	file, err := os.Create("./output.csv")
	if err != nil {
		panic(err)
	}

	// csv writer 생성
	wr := csv.NewWriter(bufio.NewWriter(file))

	// csv 내용 쓰기
	wr.Write([]string{"IAM Users"})
	wr.Write(iam)
	wr.Write([]string{""})
	wr.Write([]string{"Workspace Users"})
	wr.Write(workspace)
	wr.Flush()
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Arguments not enough")
	}

	info := os.Args
	_, err := envRegister(info)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	region := os.Args[3]
	sess, err := createSession(region)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	userIAMList, err := getListIAMUser(sess)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	userWorkspaceList, err := getListWorkspaceUser(sess)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	wrtieCSV(userIAMList, userWorkspaceList)
}